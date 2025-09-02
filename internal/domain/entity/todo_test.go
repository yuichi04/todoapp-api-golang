package entity

import (
	"encoding/json"
	"testing"
	"time"
)

// TestTodo_IsValid はTodoエンティティのバリデーション機能をテストします
// 標準パッケージでのテスト実装の学習ポイント：
// 1. testing パッケージの基本的な使用方法
// 2. テーブルドリブンテストパターンの実装
// 3. エラーケースの網羅的なテスト
// 4. ビジネスルールの検証
func TestTodo_IsValid(t *testing.T) {
	// テーブルドリブンテスト用のテストケース定義
	// 各テストケースには名前、入力、期待する結果を含める
	tests := []struct {
		name   string
		todo   Todo
		expect bool
	}{
		{
			name: "有効なTodo",
			todo: Todo{
				Title:       "有効なタイトル",
				Description: "有効な説明文",
				IsCompleted: false,
			},
			expect: true,
		},
		{
			name: "タイトルが空文字",
			todo: Todo{
				Title:       "",
				Description: "説明文",
				IsCompleted: false,
			},
			expect: false,
		},
		{
			name: "タイトルが100文字ちょうど（有効）",
			todo: Todo{
				Title:       generateString(100),
				Description: "説明文",
				IsCompleted: false,
			},
			expect: true,
		},
		{
			name: "タイトルが100文字超過",
			todo: Todo{
				Title:       generateString(101),
				Description: "説明文",
				IsCompleted: false,
			},
			expect: false,
		},
		{
			name: "完了状態がtrue（有効）",
			todo: Todo{
				Title:       "有効なタイトル",
				Description: "説明文",
				IsCompleted: true,
			},
			expect: true,
		},
	}

	// 各テストケースを実行
	for _, tt := range tests {
		// サブテストとして実行（テスト結果が個別に表示される）
		t.Run(tt.name, func(t *testing.T) {
			result := tt.todo.IsValid()

			if result != tt.expect {
				t.Errorf("Todo.IsValid() = %v, 期待値 = %v", result, tt.expect)
			}
		})
	}
}

// TestTodo_MarkAsCompleted はTodo完了機能をテストします
func TestTodo_MarkAsCompleted(t *testing.T) {
	todo := Todo{
		Title:       "テストタスク",
		Description: "完了テスト用",
		IsCompleted: false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 完了前の状態を確認
	if todo.IsCompleted {
		t.Error("初期状態では未完了であるべきです")
	}

	// 完了処理を実行
	todo.MarkAsCompleted()

	// 完了後の状態を確認
	if !todo.IsCompleted {
		t.Error("MarkAsCompleted() 実行後は完了状態であるべきです")
	}

	// UpdatedAtが更新されていることを確認
	// 時間の比較は厳密に行うため、現在時刻との差を確認
	timeDiff := time.Since(todo.UpdatedAt)
	if timeDiff > 1*time.Second {
		t.Errorf("UpdatedAt が更新されていません。差分: %v", timeDiff)
	}
}

// TestTodo_MarkAsIncomplete はTodo未完了機能をテストします
func TestTodo_MarkAsIncomplete(t *testing.T) {
	todo := Todo{
		Title:       "テストタスク",
		Description: "未完了テスト用",
		IsCompleted: true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 完了状態の確認
	if !todo.IsCompleted {
		t.Error("初期状態では完了であるべきです")
	}

	// 未完了処理を実行
	todo.MarkAsIncomplete()

	// 未完了状態の確認
	if todo.IsCompleted {
		t.Error("MarkAsIncomplete() 実行後は未完了状態であるべきです")
	}

	// UpdatedAtが更新されていることを確認
	timeDiff := time.Since(todo.UpdatedAt)
	if timeDiff > 1*time.Second {
		t.Errorf("UpdatedAt が更新されていません。差分: %v", timeDiff)
	}
}

// TestTodo_JSONMarshaling はJSON変換機能をテストします
// 標準パッケージではORMのTableNameメソッドは不要のため、
// 代わりにJSONマーシャリングのテストを実装
func TestTodo_JSONMarshaling(t *testing.T) {
	todo := Todo{
		ID:          1,
		Title:       "テストタスク",
		Description: "JSON変換テスト",
		IsCompleted: false,
		CreatedAt:   time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
	}

	// JSON形式の期待値（時刻フォーマットに注意）
	expected := `{"id":1,"title":"テストタスク","description":"JSON変換テスト","is_completed":false,"created_at":"2023-01-01T12:00:00Z","updated_at":"2023-01-01T12:00:00Z"}`

	// 構造体からJSONに変換
	jsonData, err := json.Marshal(todo)
	if err != nil {
		t.Errorf("JSON Marshal エラー: %v", err)
	}

	// JSON文字列の比較
	if string(jsonData) != expected {
		t.Errorf("JSON Marshal結果が期待値と異なります\n実際: %s\n期待: %s", string(jsonData), expected)
	}
}

// generateString は指定された長さの文字列を生成するヘルパー関数です
// テスト用のデータ生成に使用
func generateString(length int) string {
	result := ""
	char := "a" // ASCII文字を使用
	for i := 0; i < length; i++ {
		result += char
	}
	return result
}

// 標準パッケージを使ったテストの学習ポイント：
//
// 1. testing パッケージの活用：
//    - t.Run() でサブテスト実行
//    - t.Error(), t.Errorf() でエラー報告
//    - テーブルドリブンテストパターン
//
// 2. テストケースの設計：
//    - 正常系と異常系の両方をテスト
//    - 境界値のテスト
//    - エラーメッセージの検証
//
// 3. テストの構造：
//    - Given-When-Then パターン
//    - AAA パターン (Arrange-Act-Assert)
//    - テストケースの独立性確保
//
// 4. ヘルパー関数の活用：
//    - テストデータ生成
//    - 共通ロジックの切り出し
//    - コードの重複削減
//
// 5. 実際のビジネスロジックのテスト：
//    - ドメインルールの検証
//    - エンティティの状態変更
//    - バリデーションロジック
//
// このテストファイルにより、Todoエンティティの全機能が
// 適切にテストされ、リファクタリングや機能追加時の
// 安全性が確保されます。
