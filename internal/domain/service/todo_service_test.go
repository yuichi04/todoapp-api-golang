package service

import (
	"context"
	"errors"
	"testing"

	"todoapp-api-golang/internal/domain/entity"
)

// MockTodoRepository はテスト用のTodoRepositoryのモック実装です
// 標準パッケージでのモック実装の学習ポイント：
// 1. インターフェースベースのモックパターン
// 2. テスト用データの管理
// 3. エラーケースのシミュレーション
// 4. 呼び出し回数や引数の検証
type MockTodoRepository struct {
	todos       map[int]*entity.Todo
	nextID      int
	shouldError bool
	errorMsg    string
	callCounts  map[string]int
	lastCalls   map[string][]interface{}
}

// NewMockTodoRepository はモックリポジトリのコンストラクタです
func NewMockTodoRepository() *MockTodoRepository {
	return &MockTodoRepository{
		todos:      make(map[int]*entity.Todo),
		nextID:     1,
		callCounts: make(map[string]int),
		lastCalls:  make(map[string][]interface{}),
	}
}

// SetError はモックがエラーを返すように設定します
func (m *MockTodoRepository) SetError(shouldError bool, errorMsg string) {
	m.shouldError = shouldError
	m.errorMsg = errorMsg
}

// GetCallCount は指定されたメソッドの呼び出し回数を返します
func (m *MockTodoRepository) GetCallCount(method string) int {
	return m.callCounts[method]
}

// GetLastCall は指定されたメソッドの最後の呼び出し引数を返します
func (m *MockTodoRepository) GetLastCall(method string) []interface{} {
	return m.lastCalls[method]
}

// Create はTodoを作成します（モック実装）
func (m *MockTodoRepository) Create(ctx context.Context, todo *entity.Todo) (*entity.Todo, error) {
	m.callCounts["Create"]++
	m.lastCalls["Create"] = []interface{}{ctx, todo}

	if m.shouldError {
		return nil, errors.New(m.errorMsg)
	}

	// IDを設定して保存
	todo.ID = m.nextID
	m.nextID++
	
	// コピーを作成して保存（参照の問題を避ける）
	savedTodo := *todo
	m.todos[todo.ID] = &savedTodo

	return &savedTodo, nil
}

// GetByID はIDによってTodoを取得します（モック実装）
func (m *MockTodoRepository) GetByID(ctx context.Context, id int) (*entity.Todo, error) {
	m.callCounts["GetByID"]++
	m.lastCalls["GetByID"] = []interface{}{ctx, id}

	if m.shouldError {
		return nil, errors.New(m.errorMsg)
	}

	todo, exists := m.todos[id]
	if !exists {
		return nil, errors.New("todo not found")
	}

	// コピーを返す（参照の問題を避ける）
	result := *todo
	return &result, nil
}

// GetAll は全てのTodoを取得します（モック実装）
func (m *MockTodoRepository) GetAll(ctx context.Context) ([]*entity.Todo, error) {
	m.callCounts["GetAll"]++
	m.lastCalls["GetAll"] = []interface{}{ctx}

	if m.shouldError {
		return nil, errors.New(m.errorMsg)
	}

	result := make([]*entity.Todo, 0, len(m.todos))
	for _, todo := range m.todos {
		// コピーを作成
		todoCopy := *todo
		result = append(result, &todoCopy)
	}

	return result, nil
}

// Update はTodoを更新します（モック実装）
func (m *MockTodoRepository) Update(ctx context.Context, todo *entity.Todo) (*entity.Todo, error) {
	m.callCounts["Update"]++
	m.lastCalls["Update"] = []interface{}{ctx, todo}

	if m.shouldError {
		return nil, errors.New(m.errorMsg)
	}

	_, exists := m.todos[todo.ID]
	if !exists {
		return nil, errors.New("todo not found")
	}

	// コピーを作成して保存
	savedTodo := *todo
	m.todos[todo.ID] = &savedTodo

	return &savedTodo, nil
}

// Delete はTodoを削除します（モック実装）
func (m *MockTodoRepository) Delete(ctx context.Context, id int) error {
	m.callCounts["Delete"]++
	m.lastCalls["Delete"] = []interface{}{ctx, id}

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

// TestNewTodoService はTodoServiceのコンストラクタをテストします
func TestNewTodoService(t *testing.T) {
	mockRepo := NewMockTodoRepository()
	service := NewTodoService(mockRepo)

	if service == nil {
		t.Error("NewTodoService() は nil を返すべきではありません")
	}

	// サービスが正しく初期化されているかを確認
	// プライベートフィールドのテストは通常行わないが、学習用として実装
}

// TestTodoService_CreateTodo はTodo作成機能をテストします
func TestTodoService_CreateTodo(t *testing.T) {
	mockRepo := NewMockTodoRepository()
	service := NewTodoService(mockRepo)
	ctx := context.Background()

	tests := []struct {
		name      string
		todo      *entity.Todo
		wantErr   bool
		setupMock func(*MockTodoRepository)
	}{
		{
			name: "正常なTodo作成",
			todo: &entity.Todo{
				Title:       "テストタスク",
				Description: "テスト用の説明",
			},
			wantErr:   false,
			setupMock: func(m *MockTodoRepository) {},
		},
		{
			name: "空のタイトルでエラー",
			todo: &entity.Todo{
				Title:       "",
				Description: "説明",
			},
			wantErr:   true,
			setupMock: func(m *MockTodoRepository) {},
		},
		{
			name: "長すぎるタイトルでエラー",
			todo: &entity.Todo{
				Title:       generateLongString(101),
				Description: "説明",
			},
			wantErr:   true,
			setupMock: func(m *MockTodoRepository) {},
		},
		{
			name: "リポジトリエラー",
			todo: &entity.Todo{
				Title:       "有効なタイトル",
				Description: "説明",
			},
			wantErr: true,
			setupMock: func(m *MockTodoRepository) {
				m.SetError(true, "database error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックのセットアップ
			tt.setupMock(mockRepo)

			// テスト実行
			result, err := service.CreateTodo(ctx, tt.todo)

			// エラーケースの検証
			if tt.wantErr {
				if err == nil {
					t.Error("エラーが期待されましたが、発生しませんでした")
				}
				if result != nil {
					t.Error("エラー時は nil が返されるべきです")
				}
			} else {
				if err != nil {
					t.Errorf("予期しないエラーが発生しました: %v", err)
				}
				if result == nil {
					t.Error("成功時は Todo が返されるべきです")
				}
				if result.Title != tt.todo.Title {
					t.Errorf("タイトルが正しく設定されていません。取得値 = %v, 期待値 = %v", result.Title, tt.todo.Title)
				}
			}

			// モックのリセット
			mockRepo.SetError(false, "")
		})
	}
}

// TestTodoService_GetTodoByID はID指定のTodo取得機能をテストします
func TestTodoService_GetTodoByID(t *testing.T) {
	mockRepo := NewMockTodoRepository()
	service := NewTodoService(mockRepo)
	ctx := context.Background()

	// テスト用のTodoを事前に作成
	testTodo := &entity.Todo{
		ID:          1,
		Title:       "テストタスク",
		Description: "説明",
		IsCompleted: false,
	}
	mockRepo.todos[1] = testTodo

	tests := []struct {
		name      string
		id        int
		wantErr   bool
		setupMock func(*MockTodoRepository)
	}{
		{
			name:      "存在するTodoの取得",
			id:        1,
			wantErr:   false,
			setupMock: func(m *MockTodoRepository) {},
		},
		{
			name:      "存在しないTodoの取得",
			id:        999,
			wantErr:   true,
			setupMock: func(m *MockTodoRepository) {},
		},
		{
			name:    "無効なID（0）",
			id:      0,
			wantErr: true,
			setupMock: func(m *MockTodoRepository) {},
		},
		{
			name:    "無効なID（負数）",
			id:      -1,
			wantErr: true,
			setupMock: func(m *MockTodoRepository) {},
		},
		{
			name:    "リポジトリエラー",
			id:      1,
			wantErr: true,
			setupMock: func(m *MockTodoRepository) {
				m.SetError(true, "database connection error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(mockRepo)

			result, err := service.GetTodoByID(ctx, tt.id)

			if tt.wantErr {
				if err == nil {
					t.Error("エラーが期待されましたが、発生しませんでした")
				}
			} else {
				if err != nil {
					t.Errorf("予期しないエラーが発生しました: %v", err)
				}
				if result == nil {
					t.Error("成功時は Todo が返されるべきです")
				}
				if result.ID != tt.id {
					t.Errorf("IDが一致しません。取得値 = %v, 期待値 = %v", result.ID, tt.id)
				}
			}

			mockRepo.SetError(false, "")
		})
	}
}

// TestTodoService_GetAllTodos は全Todo取得機能をテストします
func TestTodoService_GetAllTodos(t *testing.T) {
	mockRepo := NewMockTodoRepository()
	service := NewTodoService(mockRepo)
	ctx := context.Background()

	tests := []struct {
		name         string
		setupData    func(*MockTodoRepository)
		setupMock    func(*MockTodoRepository)
		wantErr      bool
		expectedLen  int
	}{
		{
			name: "空のTodoリスト",
			setupData: func(m *MockTodoRepository) {},
			setupMock: func(m *MockTodoRepository) {},
			wantErr:   false,
			expectedLen: 0,
		},
		{
			name: "複数のTodo取得",
			setupData: func(m *MockTodoRepository) {
				m.todos[1] = &entity.Todo{ID: 1, Title: "タスク1"}
				m.todos[2] = &entity.Todo{ID: 2, Title: "タスク2"}
				m.todos[3] = &entity.Todo{ID: 3, Title: "タスク3"}
			},
			setupMock: func(m *MockTodoRepository) {},
			wantErr:   false,
			expectedLen: 3,
		},
		{
			name: "リポジトリエラー",
			setupData: func(m *MockTodoRepository) {},
			setupMock: func(m *MockTodoRepository) {
				m.SetError(true, "database error")
			},
			wantErr:     true,
			expectedLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// データとモックのセットアップ
			tt.setupData(mockRepo)
			tt.setupMock(mockRepo)

			result, err := service.GetAllTodos(ctx)

			if tt.wantErr {
				if err == nil {
					t.Error("エラーが期待されましたが、発生しませんでした")
				}
			} else {
				if err != nil {
					t.Errorf("予期しないエラーが発生しました: %v", err)
				}
				if len(result) != tt.expectedLen {
					t.Errorf("結果の長さが一致しません。取得値 = %v, 期待値 = %v", len(result), tt.expectedLen)
				}
			}

			// クリーンアップ
			mockRepo.SetError(false, "")
			mockRepo.todos = make(map[int]*entity.Todo)
		})
	}
}

// TestTodoService_UpdateTodo はTodo更新機能をテストします
func TestTodoService_UpdateTodo(t *testing.T) {
	mockRepo := NewMockTodoRepository()
	service := NewTodoService(mockRepo)
	ctx := context.Background()

	// テスト用のTodoを事前に作成
	existingTodo := &entity.Todo{
		ID:          1,
		Title:       "元のタイトル",
		Description: "元の説明",
		IsCompleted: false,
	}
	mockRepo.todos[1] = existingTodo

	tests := []struct {
		name      string
		todo      *entity.Todo
		wantErr   bool
		setupMock func(*MockTodoRepository)
	}{
		{
			name: "正常な更新",
			todo: &entity.Todo{
				ID:          1,
				Title:       "更新されたタイトル",
				Description: "更新された説明",
			},
			wantErr:   false,
			setupMock: func(m *MockTodoRepository) {},
		},
		{
			name: "存在しないTodoの更新",
			todo: &entity.Todo{
				ID:          999,
				Title:       "タイトル",
				Description: "説明",
			},
			wantErr:   true,
			setupMock: func(m *MockTodoRepository) {},
		},
		{
			name: "無効なタイトルで更新",
			todo: &entity.Todo{
				ID:          1,
				Title:       "",
				Description: "説明",
			},
			wantErr:   true,
			setupMock: func(m *MockTodoRepository) {},
		},
		{
			name: "リポジトリエラー",
			todo: &entity.Todo{
				ID:          1,
				Title:       "有効なタイトル",
				Description: "説明",
			},
			wantErr: true,
			setupMock: func(m *MockTodoRepository) {
				m.SetError(true, "update failed")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(mockRepo)

			result, err := service.UpdateTodo(ctx, tt.todo)

			if tt.wantErr {
				if err == nil {
					t.Error("エラーが期待されましたが、発生しませんでした")
				}
			} else {
				if err != nil {
					t.Errorf("予期しないエラーが発生しました: %v", err)
				}
				if result == nil {
					t.Error("成功時は Todo が返されるべきです")
				}
				if result.Title != tt.todo.Title {
					t.Errorf("タイトルが更新されていません。取得値 = %v, 期待値 = %v", result.Title, tt.todo.Title)
				}
			}

			mockRepo.SetError(false, "")
		})
	}
}

// TestTodoService_DeleteTodo はTodo削除機能をテストします
func TestTodoService_DeleteTodo(t *testing.T) {
	mockRepo := NewMockTodoRepository()
	service := NewTodoService(mockRepo)
	ctx := context.Background()

	// テスト用のTodoを事前に作成
	mockRepo.todos[1] = &entity.Todo{ID: 1, Title: "削除対象"}

	tests := []struct {
		name      string
		id        int
		wantErr   bool
		setupMock func(*MockTodoRepository)
	}{
		{
			name:      "正常な削除",
			id:        1,
			wantErr:   false,
			setupMock: func(m *MockTodoRepository) {},
		},
		{
			name:      "存在しないTodoの削除",
			id:        999,
			wantErr:   true,
			setupMock: func(m *MockTodoRepository) {},
		},
		{
			name:    "無効なID",
			id:      0,
			wantErr: true,
			setupMock: func(m *MockTodoRepository) {},
		},
		{
			name:    "リポジトリエラー",
			id:      1,
			wantErr: true,
			setupMock: func(m *MockTodoRepository) {
				m.SetError(true, "delete failed")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// テスト用データを再設定（前のテストで削除される可能性があるため）
			if tt.id == 1 && !tt.wantErr {
				mockRepo.todos[1] = &entity.Todo{ID: 1, Title: "削除対象"}
			}

			tt.setupMock(mockRepo)

			err := service.DeleteTodo(ctx, tt.id)

			if tt.wantErr {
				if err == nil {
					t.Error("エラーが期待されましたが、発生しませんでした")
				}
			} else {
				if err != nil {
					t.Errorf("予期しないエラーが発生しました: %v", err)
				}
			}

			mockRepo.SetError(false, "")
		})
	}
}

// TestTodoService_CompleteTodo はTodo完了機能をテストします
func TestTodoService_CompleteTodo(t *testing.T) {
	mockRepo := NewMockTodoRepository()
	service := NewTodoService(mockRepo)
	ctx := context.Background()

	// テスト用の未完了Todoを事前に作成
	incompleteTodo := &entity.Todo{
		ID:          1,
		Title:       "未完了タスク",
		Description: "説明",
		IsCompleted: false,
	}
	mockRepo.todos[1] = incompleteTodo

	tests := []struct {
		name      string
		id        int
		wantErr   bool
		setupMock func(*MockTodoRepository)
	}{
		{
			name:      "正常な完了処理",
			id:        1,
			wantErr:   false,
			setupMock: func(m *MockTodoRepository) {},
		},
		{
			name:      "存在しないTodoの完了",
			id:        999,
			wantErr:   true,
			setupMock: func(m *MockTodoRepository) {},
		},
		{
			name:    "リポジトリエラー",
			id:      1,
			wantErr: true,
			setupMock: func(m *MockTodoRepository) {
				m.SetError(true, "update failed")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// テスト用データを再設定
			mockRepo.todos[1] = &entity.Todo{
				ID:          1,
				Title:       "未完了タスク",
				IsCompleted: false,
			}

			tt.setupMock(mockRepo)

			result, err := service.CompleteTodo(ctx, tt.id)

			if tt.wantErr {
				if err == nil {
					t.Error("エラーが期待されましたが、発生しませんでした")
				}
			} else {
				if err != nil {
					t.Errorf("予期しないエラーが発生しました: %v", err)
				}
				if result == nil {
					t.Error("成功時は Todo が返されるべきです")
				}
				if !result.IsCompleted {
					t.Error("Todo が完了状態になっていません")
				}
			}

			mockRepo.SetError(false, "")
		})
	}
}

// generateLongString は指定された長さの文字列を生成するヘルパー関数です
func generateLongString(length int) string {
	result := ""
	char := "a"
	for i := 0; i < length; i++ {
		result += char
	}
	return result
}

// 標準パッケージでのサービス層テストの学習ポイント：
//
// 1. モックパターンの実装：
//    - インターフェースベースのモック作成
//    - テストデータの管理
//    - メソッド呼び出しの追跡
//
// 2. ビジネスロジックのテスト：
//    - ドメインルールの検証
//    - エラーケースの網羅
//    - エッジケースの確認
//
// 3. 依存関係の分離：
//    - リポジトリ層の抽象化
//    - テスト用の実装切り替え
//    - 外部依存の排除
//
// 4. テストの構造化：
//    - テーブルドリブンテスト
//    - セットアップ・クリーンアップ
//    - テストケースの独立性
//
// 5. カバレッジの確保：
//    - 正常系・異常系の両方
//    - 境界値テスト
//    - エラーハンドリングの検証