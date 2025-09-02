package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"todoapp-api-golang/internal/domain/entity"
)

// MockTodoService はテスト用のTodoServiceのモック実装です
// HTTPハンドラーテストでサービス層の依存関係を分離するために使用
type MockTodoService struct {
	todos       map[int]*entity.Todo
	nextID      int
	shouldError bool
	errorMsg    string
	callCounts  map[string]int
}

// NewMockTodoService はモックサービスのコンストラクタです
func NewMockTodoService() *MockTodoService {
	return &MockTodoService{
		todos:      make(map[int]*entity.Todo),
		nextID:     1,
		callCounts: make(map[string]int),
	}
}

// SetError はモックがエラーを返すように設定します
func (m *MockTodoService) SetError(shouldError bool, errorMsg string) {
	m.shouldError = shouldError
	m.errorMsg = errorMsg
}

// CreateTodo のモック実装
func (m *MockTodoService) CreateTodo(ctx context.Context, todo *entity.Todo) (*entity.Todo, error) {
	m.callCounts["CreateTodo"]++
	
	if m.shouldError {
		return nil, errors.New(m.errorMsg)
	}

	todo.ID = m.nextID
	todo.CreatedAt = time.Now()
	todo.UpdatedAt = time.Now()
	m.nextID++

	savedTodo := *todo
	m.todos[todo.ID] = &savedTodo

	return &savedTodo, nil
}

// GetTodoByID のモック実装
func (m *MockTodoService) GetTodoByID(ctx context.Context, id int) (*entity.Todo, error) {
	m.callCounts["GetTodoByID"]++
	
	if m.shouldError {
		return nil, errors.New(m.errorMsg)
	}

	todo, exists := m.todos[id]
	if !exists {
		return nil, errors.New("todo not found")
	}

	result := *todo
	return &result, nil
}

// GetAllTodos のモック実装
func (m *MockTodoService) GetAllTodos(ctx context.Context) ([]*entity.Todo, error) {
	m.callCounts["GetAllTodos"]++
	
	if m.shouldError {
		return nil, errors.New(m.errorMsg)
	}

	result := make([]*entity.Todo, 0, len(m.todos))
	for _, todo := range m.todos {
		todoCopy := *todo
		result = append(result, &todoCopy)
	}

	return result, nil
}

// UpdateTodo のモック実装
func (m *MockTodoService) UpdateTodo(ctx context.Context, todo *entity.Todo) (*entity.Todo, error) {
	m.callCounts["UpdateTodo"]++
	
	if m.shouldError {
		return nil, errors.New(m.errorMsg)
	}

	_, exists := m.todos[todo.ID]
	if !exists {
		return nil, errors.New("todo not found")
	}

	todo.UpdatedAt = time.Now()
	savedTodo := *todo
	m.todos[todo.ID] = &savedTodo

	return &savedTodo, nil
}

// DeleteTodo のモック実装
func (m *MockTodoService) DeleteTodo(ctx context.Context, id int) error {
	m.callCounts["DeleteTodo"]++
	
	if m.shouldError {
		return errors.New(m.errorMsg)
	}

	_, exists := m.todos[id]
	if !exists {
		return errors.New("todo not found")
	}

	delete(m.todos, id)
	return nil
}

// CompleteTodo のモック実装
func (m *MockTodoService) CompleteTodo(ctx context.Context, id int) (*entity.Todo, error) {
	m.callCounts["CompleteTodo"]++
	
	if m.shouldError {
		return nil, errors.New(m.errorMsg)
	}

	todo, exists := m.todos[id]
	if !exists {
		return nil, errors.New("todo not found")
	}

	todo.MarkAsCompleted()
	todo.UpdatedAt = time.Now()
	
	result := *todo
	return &result, nil
}

// IncompleteTodo のモック実装
func (m *MockTodoService) IncompleteTodo(ctx context.Context, id int) (*entity.Todo, error) {
	m.callCounts["IncompleteTodo"]++
	
	if m.shouldError {
		return nil, errors.New(m.errorMsg)
	}

	todo, exists := m.todos[id]
	if !exists {
		return nil, errors.New("todo not found")
	}

	todo.MarkAsIncomplete()
	todo.UpdatedAt = time.Now()
	
	result := *todo
	return &result, nil
}

// TestNewTodoHandler はTodoHandlerのコンストラクタをテストします
func TestNewTodoHandler(t *testing.T) {
	mockService := NewMockTodoService()
	handler := NewTodoHandler(mockService)

	if handler == nil {
		t.Error("NewTodoHandler() は nil を返すべきではありません")
	}
}

// TestTodoHandler_CreateTodo はTodo作成ハンドラーをテストします
func TestTodoHandler_CreateTodo(t *testing.T) {
	mockService := NewMockTodoService()
	handler := NewTodoHandler(mockService)

	tests := []struct {
		name           string
		method         string
		body           string
		setupMock      func(*MockTodoService)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:           "正常なTodo作成",
			method:         http.MethodPost,
			body:           `{"title":"テストタスク","description":"テスト説明"}`,
			setupMock:      func(m *MockTodoService) {},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("レスポンスのJSONパースに失敗: %v", err)
				}
				if response["title"] != "テストタスク" {
					t.Errorf("レスポンスのタイトルが正しくありません: %v", response["title"])
				}
			},
		},
		{
			name:           "不正なHTTPメソッド",
			method:         http.MethodGet,
			body:           "",
			setupMock:      func(m *MockTodoService) {},
			expectedStatus: http.StatusMethodNotAllowed,
			checkResponse:  func(t *testing.T, rec *httptest.ResponseRecorder) {},
		},
		{
			name:           "不正なJSONフォーマット",
			method:         http.MethodPost,
			body:           `{"title": invalid json}`,
			setupMock:      func(m *MockTodoService) {},
			expectedStatus: http.StatusBadRequest,
			checkResponse:  func(t *testing.T, rec *httptest.ResponseRecorder) {},
		},
		{
			name:           "空のタイトル",
			method:         http.MethodPost,
			body:           `{"title":"","description":"説明"}`,
			setupMock:      func(m *MockTodoService) {},
			expectedStatus: http.StatusBadRequest,
			checkResponse:  func(t *testing.T, rec *httptest.ResponseRecorder) {},
		},
		{
			name:   "サービス層エラー",
			method: http.MethodPost,
			body:   `{"title":"テスト","description":"説明"}`,
			setupMock: func(m *MockTodoService) {
				m.SetError(true, "database error")
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse:  func(t *testing.T, rec *httptest.ResponseRecorder) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックのセットアップ
			tt.setupMock(mockService)

			// リクエストの作成
			req := httptest.NewRequest(tt.method, "/api/v1/todos", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")

			// レスポンスレコーダーの作成
			rec := httptest.NewRecorder()

			// ハンドラーの実行
			handler.CreateTodo(rec, req)

			// ステータスコードの確認
			if rec.Code != tt.expectedStatus {
				t.Errorf("ステータスコード = %v, 期待値 = %v", rec.Code, tt.expectedStatus)
			}

			// レスポンス内容の確認
			tt.checkResponse(t, rec)

			// モックのリセット
			mockService.SetError(false, "")
		})
	}
}

// TestTodoHandler_GetAllTodos は全Todo取得ハンドラーをテストします
func TestTodoHandler_GetAllTodos(t *testing.T) {
	mockService := NewMockTodoService()
	handler := NewTodoHandler(mockService)

	tests := []struct {
		name           string
		method         string
		setupData      func(*MockTodoService)
		setupMock      func(*MockTodoService)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:      "空のTodoリスト取得",
			method:    http.MethodGet,
			setupData: func(m *MockTodoService) {},
			setupMock: func(m *MockTodoService) {},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("レスポンスのJSONパースに失敗: %v", err)
				}
				todos, ok := response["todos"].([]interface{})
				if !ok {
					t.Error("todos フィールドが配列ではありません")
					return
				}
				if len(todos) != 0 {
					t.Errorf("空のリストが期待されましたが、%d個の要素がありました", len(todos))
				}
			},
		},
		{
			name:   "複数のTodo取得",
			method: http.MethodGet,
			setupData: func(m *MockTodoService) {
				m.todos[1] = &entity.Todo{ID: 1, Title: "タスク1", Description: "説明1"}
				m.todos[2] = &entity.Todo{ID: 2, Title: "タスク2", Description: "説明2"}
			},
			setupMock: func(m *MockTodoService) {},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("レスポンスのJSONパースに失敗: %v", err)
				}
				todos, ok := response["todos"].([]interface{})
				if !ok {
					t.Error("todos フィールドが配列ではありません")
					return
				}
				if len(todos) != 2 {
					t.Errorf("2個の要素が期待されましたが、%d個の要素がありました", len(todos))
				}
			},
		},
		{
			name:      "不正なHTTPメソッド",
			method:    http.MethodPost,
			setupData: func(m *MockTodoService) {},
			setupMock: func(m *MockTodoService) {},
			expectedStatus: http.StatusMethodNotAllowed,
			checkResponse:  func(t *testing.T, rec *httptest.ResponseRecorder) {},
		},
		{
			name:      "サービス層エラー",
			method:    http.MethodGet,
			setupData: func(m *MockTodoService) {},
			setupMock: func(m *MockTodoService) {
				m.SetError(true, "database connection error")
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse:  func(t *testing.T, rec *httptest.ResponseRecorder) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// データとモックのセットアップ
			tt.setupData(mockService)
			tt.setupMock(mockService)

			// リクエストの作成
			req := httptest.NewRequest(tt.method, "/api/v1/todos", nil)

			// レスポンスレコーダーの作成
			rec := httptest.NewRecorder()

			// ハンドラーの実行
			handler.GetAllTodos(rec, req)

			// ステータスコードの確認
			if rec.Code != tt.expectedStatus {
				t.Errorf("ステータスコード = %v, 期待値 = %v", rec.Code, tt.expectedStatus)
			}

			// レスポンス内容の確認
			tt.checkResponse(t, rec)

			// クリーンアップ
			mockService.SetError(false, "")
			mockService.todos = make(map[int]*entity.Todo)
		})
	}
}

// TestTodoHandler_GetTodoByID はID指定Todo取得ハンドラーをテストします
func TestTodoHandler_GetTodoByID(t *testing.T) {
	mockService := NewMockTodoService()
	handler := NewTodoHandler(mockService)

	// テスト用データの準備
	testTodo := &entity.Todo{
		ID:          1,
		Title:       "テストタスク",
		Description: "説明",
		IsCompleted: false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	mockService.todos[1] = testTodo

	tests := []struct {
		name           string
		method         string
		setupMock      func(*MockTodoService)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:           "存在するTodo取得",
			method:         http.MethodGet,
			setupMock:      func(m *MockTodoService) {},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("レスポンスのJSONパースに失敗: %v", err)
				}
				if response["title"] != "テストタスク" {
					t.Errorf("レスポンスのタイトルが正しくありません: %v", response["title"])
				}
			},
		},
		{
			name:           "不正なHTTPメソッド",
			method:         http.MethodPost,
			setupMock:      func(m *MockTodoService) {},
			expectedStatus: http.StatusMethodNotAllowed,
			checkResponse:  func(t *testing.T, rec *httptest.ResponseRecorder) {},
		},
		{
			name:   "サービス層エラー",
			method: http.MethodGet,
			setupMock: func(m *MockTodoService) {
				m.SetError(true, "todo not found")
			},
			expectedStatus: http.StatusNotFound,
			checkResponse:  func(t *testing.T, rec *httptest.ResponseRecorder) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(mockService)

			// URLにIDパラメータを設定したリクエストを作成
			// 実際の実装ではルーターからIDが抽出されますが、
			// テストでは直接設定するかコンテキスト経由で渡す必要があります
			req := httptest.NewRequest(tt.method, "/api/v1/todos/1", nil)

			rec := httptest.NewRecorder()
			handler.GetTodoByID(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("ステータスコード = %v, 期待値 = %v", rec.Code, tt.expectedStatus)
			}

			tt.checkResponse(t, rec)
			mockService.SetError(false, "")
		})
	}
}

// TestTodoHandler_UpdateTodo はTodo更新ハンドラーをテストします
func TestTodoHandler_UpdateTodo(t *testing.T) {
	mockService := NewMockTodoService()
	handler := NewTodoHandler(mockService)

	// テスト用データの準備
	testTodo := &entity.Todo{
		ID:          1,
		Title:       "元のタイトル",
		Description: "元の説明",
		IsCompleted: false,
	}
	mockService.todos[1] = testTodo

	tests := []struct {
		name           string
		method         string
		body           string
		setupMock      func(*MockTodoService)
		expectedStatus int
	}{
		{
			name:           "正常なTodo更新",
			method:         http.MethodPut,
			body:           `{"title":"更新されたタイトル","description":"更新された説明"}`,
			setupMock:      func(m *MockTodoService) {},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "不正なHTTPメソッド",
			method:         http.MethodGet,
			body:           "",
			setupMock:      func(m *MockTodoService) {},
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "不正なJSONフォーマット",
			method:         http.MethodPut,
			body:           `invalid json`,
			setupMock:      func(m *MockTodoService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "サービス層エラー",
			method: http.MethodPut,
			body:   `{"title":"更新タイトル","description":"説明"}`,
			setupMock: func(m *MockTodoService) {
				m.SetError(true, "update failed")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(mockService)

			req := httptest.NewRequest(tt.method, "/api/v1/todos/1", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			handler.UpdateTodo(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("ステータスコード = %v, 期待値 = %v", rec.Code, tt.expectedStatus)
			}

			mockService.SetError(false, "")
		})
	}
}

// TestTodoHandler_DeleteTodo はTodo削除ハンドラーをテストします
func TestTodoHandler_DeleteTodo(t *testing.T) {
	mockService := NewMockTodoService()
	handler := NewTodoHandler(mockService)

	// テスト用データの準備
	mockService.todos[1] = &entity.Todo{ID: 1, Title: "削除対象"}

	tests := []struct {
		name           string
		method         string
		setupMock      func(*MockTodoService)
		expectedStatus int
	}{
		{
			name:           "正常なTodo削除",
			method:         http.MethodDelete,
			setupMock:      func(m *MockTodoService) {},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "不正なHTTPメソッド",
			method:         http.MethodGet,
			setupMock:      func(m *MockTodoService) {},
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:   "サービス層エラー",
			method: http.MethodDelete,
			setupMock: func(m *MockTodoService) {
				m.SetError(true, "delete failed")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// テスト用データを再設定
			mockService.todos[1] = &entity.Todo{ID: 1, Title: "削除対象"}
			tt.setupMock(mockService)

			req := httptest.NewRequest(tt.method, "/api/v1/todos/1", nil)
			rec := httptest.NewRecorder()
			handler.DeleteTodo(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("ステータスコード = %v, 期待値 = %v", rec.Code, tt.expectedStatus)
			}

			mockService.SetError(false, "")
		})
	}
}

// 標準パッケージでのHTTPハンドラーテストの学習ポイント：
//
// 1. net/http/httptest パッケージの活用：
//    - httptest.NewRequest() でテスト用リクエスト作成
//    - httptest.ResponseRecorder でレスポンス記録
//    - HTTPステータスコードとボディの検証
//
// 2. JSONハンドリングのテスト：
//    - encoding/json パッケージでのシリアライゼーション
//    - リクエスト・レスポンスのJSON検証
//    - 不正なJSONフォーマットのテスト
//
// 3. HTTPメソッドのテスト：
//    - 各エンドポイントでサポートされるメソッドの確認
//    - 不正なメソッドでのエラーハンドリング
//
// 4. エラーケースの網羅：
//    - サービス層エラーの伝播
//    - バリデーションエラー
//    - HTTPステータスコードの適切な設定
//
// 5. モックを使った単体テスト：
//    - 外部依存（サービス層）の分離
//    - テスト専用のモック実装
//    - エラーケースのシミュレーション