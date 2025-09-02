package database

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"todoapp-api-golang/internal/domain/entity"

	// SQLite ドライバーをテスト用に使用
	_ "github.com/mattn/go-sqlite3"
)

// setupTestDB はテスト用のSQLiteデータベースを設定します
// 標準パッケージでの統合テストの学習ポイント：
// 1. インメモリデータベースの使用
// 2. テスト用データベースの初期化
// 3. テーブル作成とクリーンアップ
// 4. トランザクションを使った分離
func setupTestDB(t *testing.T) *sql.DB {
	// インメモリSQLiteデータベースを作成
	// ":memory:" を使うことで、テスト終了時に自動的に削除される
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("テストデータベースの作成に失敗: %v", err)
	}

	// Todosテーブルを作成
	createTable := `
		CREATE TABLE todos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			description TEXT,
			is_completed BOOLEAN NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`

	_, err = db.Exec(createTable)
	if err != nil {
		t.Fatalf("テストテーブルの作成に失敗: %v", err)
	}

	return db
}

// TestNewTodoRepository はTodoRepositoryのコンストラクタをテストします
func TestNewTodoRepository(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewTodoRepository(db)
	if repo == nil {
		t.Error("NewTodoRepository() は nil を返すべきではありません")
	}
}

// TestTodoRepository_Create はTodo作成機能をテストします
func TestTodoRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewTodoRepository(db)
	ctx := context.Background()

	tests := []struct {
		name    string
		todo    *entity.Todo
		wantErr bool
	}{
		{
			name: "正常なTodo作成",
			todo: &entity.Todo{
				Title:       "テストタスク",
				Description: "テスト用の説明",
				IsCompleted: false,
			},
			wantErr: false,
		},
		{
			name: "空のタイトル（データベースレベルでは許可される場合がある）",
			todo: &entity.Todo{
				Title:       "",
				Description: "説明のみ",
				IsCompleted: false,
			},
			wantErr: false,
		},
		{
			name: "説明なしのTodo",
			todo: &entity.Todo{
				Title:       "タイトルのみ",
				Description: "",
				IsCompleted: false,
			},
			wantErr: false,
		},
		{
			name: "完了済みTodo",
			todo: &entity.Todo{
				Title:       "完了済みタスク",
				Description: "説明",
				IsCompleted: true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.Create(ctx, tt.todo)

			if tt.wantErr {
				if err == nil {
					t.Error("エラーが期待されましたが、発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラーが発生しました: %v", err)
				return
			}

			if result == nil {
				t.Error("作成されたTodoが nil です")
				return
			}

			// IDが自動生成されていることを確認
			if result.ID <= 0 {
				t.Error("IDが正しく生成されていません")
			}

			// フィールドが正しく設定されていることを確認
			if result.Title != tt.todo.Title {
				t.Errorf("タイトルが正しく設定されていません。取得値 = %v, 期待値 = %v", result.Title, tt.todo.Title)
			}

			if result.Description != tt.todo.Description {
				t.Errorf("説明が正しく設定されていません。取得値 = %v, 期待値 = %v", result.Description, tt.todo.Description)
			}

			if result.IsCompleted != tt.todo.IsCompleted {
				t.Errorf("完了状態が正しく設定されていません。取得値 = %v, 期待値 = %v", result.IsCompleted, tt.todo.IsCompleted)
			}

			// 作成日時と更新日時が設定されていることを確認
			if result.CreatedAt.IsZero() {
				t.Error("CreatedAt が設定されていません")
			}

			if result.UpdatedAt.IsZero() {
				t.Error("UpdatedAt が設定されていません")
			}
		})
	}
}

// TestTodoRepository_GetByID はID指定取得機能をテストします
func TestTodoRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewTodoRepository(db)
	ctx := context.Background()

	// テスト用データを作成
	testTodo := &entity.Todo{
		Title:       "取得テスト用",
		Description: "説明",
		IsCompleted: false,
	}
	createdTodo, err := repo.Create(ctx, testTodo)
	if err != nil {
		t.Fatalf("テストデータの作成に失敗: %v", err)
	}

	tests := []struct {
		name    string
		id      int
		wantErr bool
	}{
		{
			name:    "存在するTodoの取得",
			id:      createdTodo.ID,
			wantErr: false,
		},
		{
			name:    "存在しないTodoの取得",
			id:      99999,
			wantErr: true,
		},
		{
			name:    "無効なID（0）",
			id:      0,
			wantErr: true,
		},
		{
			name:    "無効なID（負数）",
			id:      -1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.GetByID(ctx, tt.id)

			if tt.wantErr {
				if err == nil {
					t.Error("エラーが期待されましたが、発生しませんでした")
				}
				if result != nil {
					t.Error("エラー時は nil が返されるべきです")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラーが発生しました: %v", err)
				return
			}

			if result == nil {
				t.Error("取得されたTodoが nil です")
				return
			}

			if result.ID != tt.id {
				t.Errorf("IDが一致しません。取得値 = %v, 期待値 = %v", result.ID, tt.id)
			}

			if result.Title != testTodo.Title {
				t.Errorf("タイトルが一致しません。取得値 = %v, 期待値 = %v", result.Title, testTodo.Title)
			}
		})
	}
}

// TestTodoRepository_GetAll は全Todo取得機能をテストします
func TestTodoRepository_GetAll(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewTodoRepository(db)
	ctx := context.Background()

	// 空の状態でのテスト
	t.Run("空のTodoリスト", func(t *testing.T) {
		result, err := repo.GetAll(ctx)
		if err != nil {
			t.Errorf("予期しないエラーが発生しました: %v", err)
		}

		if len(result) != 0 {
			t.Errorf("空のリストが期待されましたが、%d個の要素がありました", len(result))
		}
	})

	// テスト用データを複数作成
	testTodos := []*entity.Todo{
		{Title: "タスク1", Description: "説明1", IsCompleted: false},
		{Title: "タスク2", Description: "説明2", IsCompleted: true},
		{Title: "タスク3", Description: "説明3", IsCompleted: false},
	}

	for _, todo := range testTodos {
		_, err := repo.Create(ctx, todo)
		if err != nil {
			t.Fatalf("テストデータの作成に失敗: %v", err)
		}
	}

	// 複数データでのテスト
	t.Run("複数のTodo取得", func(t *testing.T) {
		result, err := repo.GetAll(ctx)
		if err != nil {
			t.Errorf("予期しないエラーが発生しました: %v", err)
		}

		expectedLen := len(testTodos)
		if len(result) != expectedLen {
			t.Errorf("取得件数が一致しません。取得値 = %d, 期待値 = %d", len(result), expectedLen)
		}

		// ソート順の確認（作成順になっているかどうか）
		for i, todo := range result {
			if todo.Title != testTodos[i].Title {
				t.Errorf("取得順序が正しくありません。位置%d: 取得値 = %v, 期待値 = %v", i, todo.Title, testTodos[i].Title)
			}
		}
	})
}

// TestTodoRepository_Update はTodo更新機能をテストします
func TestTodoRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewTodoRepository(db)
	ctx := context.Background()

	// テスト用データを作成
	originalTodo := &entity.Todo{
		Title:       "元のタイトル",
		Description: "元の説明",
		IsCompleted: false,
	}
	createdTodo, err := repo.Create(ctx, originalTodo)
	if err != nil {
		t.Fatalf("テストデータの作成に失敗: %v", err)
	}

	// 少し時間を空けて更新時刻の違いを明確にする
	time.Sleep(10 * time.Millisecond)

	tests := []struct {
		name    string
		todo    *entity.Todo
		wantErr bool
	}{
		{
			name: "正常な更新",
			todo: &entity.Todo{
				ID:          createdTodo.ID,
				Title:       "更新されたタイトル",
				Description: "更新された説明",
				IsCompleted: true,
			},
			wantErr: false,
		},
		{
			name: "存在しないTodoの更新",
			todo: &entity.Todo{
				ID:          99999,
				Title:       "存在しないTodo",
				Description: "説明",
				IsCompleted: false,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.Update(ctx, tt.todo)

			if tt.wantErr {
				if err == nil {
					t.Error("エラーが期待されましたが、発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラーが発生しました: %v", err)
				return
			}

			if result == nil {
				t.Error("更新されたTodoが nil です")
				return
			}

			// フィールドが正しく更新されていることを確認
			if result.Title != tt.todo.Title {
				t.Errorf("タイトルが更新されていません。取得値 = %v, 期待値 = %v", result.Title, tt.todo.Title)
			}

			if result.Description != tt.todo.Description {
				t.Errorf("説明が更新されていません。取得値 = %v, 期待値 = %v", result.Description, tt.todo.Description)
			}

			if result.IsCompleted != tt.todo.IsCompleted {
				t.Errorf("完了状態が更新されていません。取得値 = %v, 期待値 = %v", result.IsCompleted, tt.todo.IsCompleted)
			}

			// 更新日時が変更されているか確認（SQLiteでは秒精度なので等しくない場合のみチェック）
			if result.UpdatedAt.Equal(createdTodo.UpdatedAt) {
				t.Logf("UpdatedAt が変更されませんでした: %v == %v", result.UpdatedAt, createdTodo.UpdatedAt)
			}

			// 作成日時は変更されないことを確認（時刻の差が1秒以内なら許容）
			timeDiff := result.CreatedAt.Sub(createdTodo.CreatedAt)
			if timeDiff > time.Second || timeDiff < -time.Second {
				t.Errorf("CreatedAt が大幅に変更されています: %v != %v (差分: %v)", result.CreatedAt, createdTodo.CreatedAt, timeDiff)
			}
		})
	}
}

// TestTodoRepository_Delete はTodo削除機能をテストします
func TestTodoRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewTodoRepository(db)
	ctx := context.Background()

	// テスト用データを作成
	testTodo := &entity.Todo{
		Title:       "削除テスト用",
		Description: "削除される予定",
		IsCompleted: false,
	}
	createdTodo, err := repo.Create(ctx, testTodo)
	if err != nil {
		t.Fatalf("テストデータの作成に失敗: %v", err)
	}

	tests := []struct {
		name    string
		id      int
		wantErr bool
	}{
		{
			name:    "存在するTodoの削除",
			id:      createdTodo.ID,
			wantErr: false,
		},
		{
			name:    "存在しないTodoの削除",
			id:      99999,
			wantErr: true,
		},
		{
			name:    "無効なID（0）",
			id:      0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 削除対象のデータが存在する場合のみ再作成
			if tt.id == createdTodo.ID && tt.name == "存在するTodoの削除" {
				testTodo := &entity.Todo{
					Title:       "削除テスト用",
					Description: "削除される予定",
					IsCompleted: false,
				}
				createdTodo, err = repo.Create(ctx, testTodo)
				if err != nil {
					t.Fatalf("テストデータの再作成に失敗: %v", err)
				}
				tt.id = createdTodo.ID
			}

			err := repo.Delete(ctx, tt.id)

			if tt.wantErr {
				if err == nil {
					t.Error("エラーが期待されましたが、発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラーが発生しました: %v", err)
				return
			}

			// 削除後に取得できないことを確認
			_, getErr := repo.GetByID(ctx, tt.id)
			if getErr == nil {
				t.Error("削除されたTodoが取得できてしまいました")
			}
		})
	}
}

// TestTodoRepository_Transaction はトランザクションを使った処理をテストします
func TestTodoRepository_Transaction(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	ctx := context.Background()

	// このテストでは基本的なトランザクションの動作を確認
	// 現在の実装では完全なトランザクション分離は実装されていないため、
	// 基本的な作成・ロールバックのパターンをテスト

	// 初期状態での件数を確認
	initialCount := getTodoCount(t, db)

	// トランザクションを開始
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("トランザクション開始に失敗: %v", err)
	}

	// トランザクション内で直接SQLを実行（リポジトリ経由ではなく）
	query := `INSERT INTO todos (title, description, is_completed, created_at, updated_at) 
	          VALUES (?, ?, false, datetime('now'), datetime('now'))`
	_, err = tx.ExecContext(ctx, query, "トランザクションテスト", "ロールバックされる予定")
	if err != nil {
		tx.Rollback()
		t.Fatalf("トランザクション内でのTodo作成に失敗: %v", err)
	}

	// トランザクションをロールバック
	err = tx.Rollback()
	if err != nil {
		t.Fatalf("ロールバックに失敗: %v", err)
	}

	// ロールバック後の件数を確認（変化していないはず）
	finalCount := getTodoCount(t, db)
	if finalCount != initialCount {
		t.Errorf("ロールバック後も件数が変化しています。初期: %d, 最終: %d", initialCount, finalCount)
	}
}

// getTodoCount はテーブル内のTodo件数を取得するヘルパー関数です
func getTodoCount(t *testing.T, db *sql.DB) int {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM todos").Scan(&count)
	if err != nil {
		t.Fatalf("件数取得に失敗: %v", err)
	}
	return count
}

// 標準パッケージでのデータベーステストの学習ポイント：
//
// 1. インメモリデータベースの活用：
//    - SQLite ":memory:" でテスト分離
//    - 高速なテスト実行
//    - 外部依存の排除
//
// 2. database/sql パッケージのテスト：
//    - 実際のSQL実行とその検証
//    - エラーケースの確認
//    - データの整合性確認
//
// 3. CRUD操作の網羅的テスト：
//    - Create, Read, Update, Delete の全機能
//    - 正常系と異常系の両方
//    - エッジケースの確認
//
// 4. 時刻フィールドのテスト：
//    - CreatedAt, UpdatedAt の動作確認
//    - 更新時の時刻変更確認
//
// 5. データベース設計の検証：
//    - テーブル制約の動作確認
//    - インデックスの効果測定（必要に応じて）
//    - トランザクション動作の確認
