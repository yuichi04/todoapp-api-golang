package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// TestChainMiddleware はミドルウェアチェーン機能をテストします
// 標準パッケージでのミドルウェアテストの学習ポイント：
// 1. httptest パッケージを使ったHTTPテスト
// 2. ミドルウェア実行順序の検証
// 3. リクエスト・レスポンスの変更確認
// 4. エラーハンドリングのテスト
func TestChainMiddleware(t *testing.T) {
	// 実行順序を記録するためのスライス
	var executionOrder []string

	// テスト用ミドルウェア1
	middleware1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			executionOrder = append(executionOrder, "middleware1-before")
			next.ServeHTTP(w, r)
			executionOrder = append(executionOrder, "middleware1-after")
		})
	}

	// テスト用ミドルウェア2
	middleware2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			executionOrder = append(executionOrder, "middleware2-before")
			next.ServeHTTP(w, r)
			executionOrder = append(executionOrder, "middleware2-after")
		})
	}

	// 最終ハンドラー
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		executionOrder = append(executionOrder, "final-handler")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// ミドルウェアチェーンを構築
	chainedHandler := ChainMiddleware(middleware1, middleware2)(finalHandler)

	// テスト実行
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	chainedHandler.ServeHTTP(rec, req)

	// 実行順序の確認
	expectedOrder := []string{
		"middleware1-before",
		"middleware2-before",
		"final-handler",
		"middleware2-after",
		"middleware1-after",
	}

	if len(executionOrder) != len(expectedOrder) {
		t.Fatalf("実行順序の長さが一致しません。取得値 = %d, 期待値 = %d", len(executionOrder), len(expectedOrder))
	}

	for i, expected := range expectedOrder {
		if executionOrder[i] != expected {
			t.Errorf("実行順序が正しくありません。位置%d: 取得値 = %s, 期待値 = %s", i, executionOrder[i], expected)
		}
	}

	// レスポンスの確認
	if rec.Code != http.StatusOK {
		t.Errorf("ステータスコード = %d, 期待値 = %d", rec.Code, http.StatusOK)
	}

	if rec.Body.String() != "OK" {
		t.Errorf("レスポンスボディ = %s, 期待値 = OK", rec.Body.String())
	}
}

// TestLoggingMiddleware はログ出力ミドルウェアをテストします
func TestLoggingMiddleware(t *testing.T) {
	// テスト用ハンドラー
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})

	// ログ出力ミドルウェアを適用
	handler := LoggingMiddleware(testHandler)

	tests := []struct {
		name   string
		method string
		path   string
		body   string
	}{
		{
			name:   "GETリクエスト",
			method: http.MethodGet,
			path:   "/api/v1/todos",
			body:   "",
		},
		{
			name:   "POSTリクエスト",
			method: http.MethodPost,
			path:   "/api/v1/todos",
			body:   `{"title":"test"}`,
		},
		{
			name:   "ルートパス",
			method: http.MethodGet,
			path:   "/",
			body:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body *bytes.Buffer
			if tt.body != "" {
				body = bytes.NewBufferString(tt.body)
			} else {
				body = bytes.NewBuffer(nil)
			}

			req := httptest.NewRequest(tt.method, tt.path, body)
			if tt.body != "" {
				req.Header.Set("Content-Type", "application/json")
			}

			rec := httptest.NewRecorder()

			// ログミドルウェアを実行
			handler.ServeHTTP(rec, req)

			// レスポンスの基本確認
			if rec.Code != http.StatusOK {
				t.Errorf("ステータスコード = %d, 期待値 = %d", rec.Code, http.StatusOK)
			}

			// NOTE: 実際のログ出力のテストは標準出力をキャプチャする必要があります
			// この例では、ミドルウェアが正常に動作することを確認しています
		})
	}
}

// TestDetailedLoggingMiddleware は詳細ログミドルウェアをテストします
func TestDetailedLoggingMiddleware(t *testing.T) {
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Test-Header", "test-value")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("created"))
	})

	handler := DetailedLoggingMiddleware(testHandler)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/todos", bytes.NewBufferString(`{"title":"test"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "test-agent")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	// レスポンスの確認
	if rec.Code != http.StatusCreated {
		t.Errorf("ステータスコード = %d, 期待値 = %d", rec.Code, http.StatusCreated)
	}

	if rec.Body.String() != "created" {
		t.Errorf("レスポンスボディ = %s, 期待値 = created", rec.Body.String())
	}

	// レスポンスヘッダーが設定されていることを確認
	if rec.Header().Get("X-Test-Header") != "test-value" {
		t.Errorf("レスポンスヘッダーが正しく設定されていません")
	}
}

// TestRequestIDMiddleware はリクエストIDミドルウェアをテストします
func TestRequestIDMiddleware(t *testing.T) {
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	handler := RequestIDMiddleware(testHandler)

	tests := []struct {
		name            string
		existingReqID   string
		expectGenerated bool
	}{
		{
			name:            "既存のリクエストIDなし",
			existingReqID:   "",
			expectGenerated: true,
		},
		{
			name:            "既存のリクエストIDあり",
			existingReqID:   "existing-request-id-123",
			expectGenerated: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/test", nil)

			if tt.existingReqID != "" {
				req.Header.Set("X-Request-ID", tt.existingReqID)
			}

			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			// レスポンスの確認
			if rec.Code != http.StatusOK {
				t.Errorf("ステータスコード = %d, 期待値 = %d", rec.Code, http.StatusOK)
			}

			// レスポンスヘッダーにリクエストIDが含まれていることを確認
			responseReqID := rec.Header().Get("X-Request-ID")
			if responseReqID == "" {
				t.Error("レスポンスにX-Request-IDヘッダーが設定されていません")
			}

			if tt.expectGenerated {
				// 生成されたIDの形式確認（prefix "req_" で始まることを確認）
				if !strings.HasPrefix(responseReqID, "req_") {
					t.Errorf("生成されたリクエストIDの形式が正しくありません: %s", responseReqID)
				}
			} else {
				// 既存のIDがそのまま使用されることを確認
				if responseReqID != tt.existingReqID {
					t.Errorf("既存のリクエストIDが保持されていません。取得値 = %s, 期待値 = %s", responseReqID, tt.existingReqID)
				}
			}
		})
	}
}

// TestRecoveryMiddleware はパニック回復ミドルウェアをテストします
func TestRecoveryMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		shouldPanic    bool
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "正常処理",
			shouldPanic:    false,
			expectedStatus: http.StatusOK,
			expectedBody:   "OK",
		},
		{
			name:           "パニック発生時の回復",
			shouldPanic:    true,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Internal Server Error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.shouldPanic {
					panic("テスト用パニック")
				}
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			})

			handler := RecoveryMiddleware(testHandler)

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			rec := httptest.NewRecorder()

			// パニックが発生しても回復することを確認
			handler.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("ステータスコード = %d, 期待値 = %d", rec.Code, tt.expectedStatus)
			}

			if rec.Body.String() != tt.expectedBody {
				t.Errorf("レスポンスボディ = %s, 期待値 = %s", rec.Body.String(), tt.expectedBody)
			}
		})
	}
}

// TestResponseRecorder はResponseRecorderの動作をテストします
func TestResponseRecorder(t *testing.T) {
	// 元のResponseWriterを作成
	rec := httptest.NewRecorder()

	// ResponseRecorderでラップ
	recorder := NewResponseRecorder(rec)

	// デフォルト値の確認
	if recorder.statusCode != http.StatusOK {
		t.Errorf("デフォルトステータスコード = %d, 期待値 = %d", recorder.statusCode, http.StatusOK)
	}

	if recorder.responseSize != 0 {
		t.Errorf("デフォルトレスポンスサイズ = %d, 期待値 = 0", recorder.responseSize)
	}

	// WriteHeaderのテスト
	recorder.WriteHeader(http.StatusCreated)
	if recorder.statusCode != http.StatusCreated {
		t.Errorf("ステータスコード = %d, 期待値 = %d", recorder.statusCode, http.StatusCreated)
	}

	// Writeのテスト
	testData := []byte("test response data")
	size, err := recorder.Write(testData)
	if err != nil {
		t.Errorf("Write エラー: %v", err)
	}

	expectedSize := len(testData)
	if size != expectedSize {
		t.Errorf("書き込みサイズ = %d, 期待値 = %d", size, expectedSize)
	}

	if recorder.responseSize != expectedSize {
		t.Errorf("記録されたレスポンスサイズ = %d, 期待値 = %d", recorder.responseSize, expectedSize)
	}

	// 複数回のWriteテスト
	additionalData := []byte(" additional")
	_, err = recorder.Write(additionalData)
	if err != nil {
		t.Errorf("2回目のWrite エラー: %v", err)
	}

	totalExpectedSize := len(testData) + len(additionalData)
	if recorder.responseSize != totalExpectedSize {
		t.Errorf("累積レスポンスサイズ = %d, 期待値 = %d", recorder.responseSize, totalExpectedSize)
	}
}

// TestGenerateRequestID はリクエストID生成機能をテストします
func TestGenerateRequestID(t *testing.T) {
	// 複数のIDを生成して一意性を確認
	ids := make(map[string]bool)
	numTests := 100

	for i := 0; i < numTests; i++ {
		id := generateRequestID()

		// 空でないことを確認
		if id == "" {
			t.Error("生成されたリクエストIDが空です")
		}

		// プレフィックスの確認
		if !strings.HasPrefix(id, "req_") {
			t.Errorf("リクエストIDのプレフィックスが正しくありません: %s", id)
		}

		// 一意性の確認
		if ids[id] {
			t.Errorf("重複するリクエストIDが生成されました: %s", id)
		}
		ids[id] = true

		// 少し時間を空けて異なる時刻でのID生成を確認
		if i%10 == 0 {
			time.Sleep(time.Nanosecond)
		}
	}

	// 全て異なるIDが生成されたことを確認
	if len(ids) != numTests {
		t.Errorf("一意なIDの数 = %d, 期待値 = %d", len(ids), numTests)
	}
}

// 標準パッケージでのミドルウェアテストの学習ポイント：
//
// 1. httptest パッケージの活用：
//    - NewRequest() でテスト用リクエスト作成
//    - ResponseRecorder でレスポンス記録
//    - HTTPハンドラーのテスト手法
//
// 2. ミドルウェアパターンのテスト：
//    - チェーン実行順序の確認
//    - リクエスト/レスポンス変更の検証
//    - エラーハンドリングの確認
//
// 3. 時刻依存処理のテスト：
//    - リクエストID生成の一意性確認
//    - 処理時間測定の検証
//
// 4. パニック処理のテスト：
//    - recover() 動作の確認
//    - エラーレスポンスの検証
//
// 5. HTTP機能の詳細テスト：
//    - ヘッダー操作の確認
//    - ステータスコード設定
//    - レスポンスボディの内容確認
