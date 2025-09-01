package middleware

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// ResponseRecorder は標準のhttp.ResponseWriterをラップして
// ステータスコードとレスポンスサイズを記録するための構造体です
//
// 標準パッケージでのレスポンス記録の学習ポイント：
// 1. http.ResponseWriter インターフェースの実装
// 2. 埋め込み（Embedding）によるラップ
// 3. WriteHeader() のオーバーライド
// 4. レスポンス情報の記録
type ResponseRecorder struct {
	http.ResponseWriter                // 埋め込みで元のResponseWriterの機能を継承
	statusCode          int            // HTTPステータスコードを記録
	responseSize        int            // レスポンスサイズ（バイト数）を記録
}

// NewResponseRecorder はResponseRecorderのコンストラクタです
func NewResponseRecorder(w http.ResponseWriter) *ResponseRecorder {
	return &ResponseRecorder{
		ResponseWriter: w,
		statusCode:     http.StatusOK, // デフォルトは200 OK
		responseSize:   0,
	}
}

// WriteHeader はHTTPステータスコードを設定し、それを記録します
// http.ResponseWriter インターフェースのメソッドをオーバーライド
func (r *ResponseRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

// Write はレスポンスボディを書き込み、そのサイズを記録します
// http.ResponseWriter インターフェースのメソッドをオーバーライド
func (r *ResponseRecorder) Write(data []byte) (int, error) {
	size, err := r.ResponseWriter.Write(data)
	r.responseSize += size
	return size, err
}

// LoggingMiddleware はHTTPリクエストとレスポンスをログ出力するミドルウェアです
//
// 標準パッケージでのログ機能の学習ポイント：
// 1. log パッケージを使った構造化ログ
// 2. リクエスト処理時間の計測
// 3. レスポンス情報の記録
// 4. 標準的なアクセスログフォーマット
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. 処理開始時刻を記録
		start := time.Now()

		// 2. ResponseWriterをラップしてレスポンス情報を記録可能にする
		recorder := NewResponseRecorder(w)

		// 3. 次のハンドラーを呼び出し
		// ここで実際のAPI処理が実行される
		next.ServeHTTP(recorder, r)

		// 4. 処理完了後にログを出力
		duration := time.Since(start)

		// Apache Combined Log Format に近い形式でログ出力
		// [timestamp] method path status size duration
		log.Printf("%s %s %s %d %d %v",
			r.RemoteAddr,           // クライアントのIPアドレス
			r.Method,               // HTTPメソッド（GET, POST, etc）
			r.URL.Path,             // リクエストパス
			recorder.statusCode,    // HTTPステータスコード
			recorder.responseSize,  // レスポンスサイズ（バイト）
			duration,               // 処理時間
		)
	})
}

// DetailedLoggingMiddleware はより詳細な情報をログ出力するミドルウェアです
// 開発環境やデバッグ用途で使用
func DetailedLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 処理開始時刻を記録
		start := time.Now()

		// リクエスト情報をログ出力
		log.Printf("→ %s %s %s", r.Method, r.URL.Path, r.Proto)
		log.Printf("  Host: %s", r.Host)
		log.Printf("  User-Agent: %s", r.Header.Get("User-Agent"))
		log.Printf("  Content-Type: %s", r.Header.Get("Content-Type"))
		log.Printf("  Content-Length: %s", r.Header.Get("Content-Length"))

		// ResponseWriterをラップ
		recorder := NewResponseRecorder(w)

		// 次のハンドラーを呼び出し
		next.ServeHTTP(recorder, r)

		// 処理完了後の詳細ログ出力
		duration := time.Since(start)
		
		log.Printf("← %s %s %d %d %v",
			r.Method,
			r.URL.Path,
			recorder.statusCode,
			recorder.responseSize,
			duration,
		)

		// レスポンスヘッダー情報も出力（開発時のデバッグ用）
		for key, values := range recorder.Header() {
			for _, value := range values {
				log.Printf("  %s: %s", key, value)
			}
		}
	})
}

// RequestIDMiddleware は各リクエストに一意のIDを付与するミドルウェアです
// 分散システムでのリクエスト追跡に使用
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. 既存のリクエストIDをチェック（ロードバランサー等から）
		requestID := r.Header.Get("X-Request-ID")
		
		// 2. リクエストIDがない場合は生成
		if requestID == "" {
			requestID = generateRequestID()
		}

		// 3. レスポンスヘッダーにリクエストIDを設定
		w.Header().Set("X-Request-ID", requestID)

		// 4. ログにリクエストIDを出力
		log.Printf("Request ID: %s - %s %s", requestID, r.Method, r.URL.Path)

		// 5. 次のハンドラーを呼び出し
		next.ServeHTTP(w, r)
	})
}

// RecoveryMiddleware はパニックを捕捉して適切にエラーレスポンスを返すミドルウェアです
// アプリケーションのクラッシュを防ぐ重要な安全装置
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// defer と recover() でパニックを捕捉
		defer func() {
			if err := recover(); err != nil {
				// パニックをログに記録
				log.Printf("PANIC: %v", err)
				
				// スタックトレースも出力（開発環境）
				// 本番環境では機密情報を含む可能性があるため注意
				log.Printf("Request: %s %s", r.Method, r.URL.Path)
				
				// クライアントには500エラーを返す
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		// 次のハンドラーを呼び出し
		next.ServeHTTP(w, r)
	})
}

// --- ヘルパー関数 ---

// generateRequestID は簡単なリクエストID生成関数です
// 実際のアプリケーションではUUID等を使用することを推奨
func generateRequestID() string {
	// 現在時刻をベースにした簡単なID生成
	// 本格的な実装ではcrypto/randやgoogle/uuidパッケージを使用
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("req_%d", timestamp)
}


// ChainMiddleware は複数のミドルウェアを連鎖させるためのヘルパー関数です
// 標準パッケージでのミドルウェアチェーンの学習
func ChainMiddleware(middlewares ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(final http.Handler) http.Handler {
		// 右から左にミドルウェアを適用（逆順）
		// 例：Chain(A, B, C)(handler) → A(B(C(handler)))
		for i := len(middlewares) - 1; i >= 0; i-- {
			final = middlewares[i](final)
		}
		return final
	}
}

// 標準パッケージでのログミドルウェアの学習ポイント：
//
// 1. ResponseWriter のラッピング：
//    - 埋め込み（embedding）での機能拡張
//    - インターフェース実装の継承
//    - メソッドオーバーライドによる拡張
//
// 2. 時間計測：
//    - time.Now() と time.Since() の使用
//    - パフォーマンス測定の基本パターン
//
// 3. ログフォーマット：
//    - 標準的なアクセスログ形式
//    - 構造化ログの実装
//    - デバッグ情報の適切な出力
//
// 4. エラー処理：
//    - パニックの捕捉と復旧
//    - ログ記録と適切なエラーレスポンス
//
// 5. ミドルウェアチェーン：
//    - 複数ミドルウェアの組み合わせ
//    - 処理順序の制御
//    - 再利用可能な設計
//
// 使用例：
// ```go
// handler := ChainMiddleware(
//     RecoveryMiddleware,
//     LoggingMiddleware,
//     CORSMiddleware(DefaultCORSConfig()),
// )(todoHandler)
// ```