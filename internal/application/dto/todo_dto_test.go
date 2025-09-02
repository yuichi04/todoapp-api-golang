package dto

import (
	"encoding/json"
	"testing"
	"time"

	"todoapp-api-golang/internal/domain/entity"
)

// TestCreateTodoRequest_ToEntity はCreateRequestのエンティティ変換をテストします
// 標準パッケージでのDTO変換テストの学習ポイント：
// 1. データ変換ロジックの検証
// 2. バリデーション処理の確認
// 3. マッピング処理の正確性確認
func TestCreateTodoRequest_ToEntity(t *testing.T) {
	tests := []struct {
		name    string
		request CreateTodoRequest
		want    *entity.Todo
	}{
		{
			name: "正常なリクエスト変換",
			request: CreateTodoRequest{
				Title:       "テストタスク",
				Description: "説明文",
			},
			want: &entity.Todo{
				Title:       "テストタスク",
				Description: "説明文",
				IsCompleted: false,
			},
		},
		{
			name: "説明なしのリクエスト",
			request: CreateTodoRequest{
				Title:       "タイトルのみ",
				Description: "",
			},
			want: &entity.Todo{
				Title:       "タイトルのみ",
				Description: "",
				IsCompleted: false,
			},
		},
		{
			name: "空のタイトル",
			request: CreateTodoRequest{
				Title:       "",
				Description: "説明文",
			},
			want: &entity.Todo{
				Title:       "",
				Description: "説明文",
				IsCompleted: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.request.ToEntity()

			if got == nil {
				t.Error("変換されたエンティティが nil です")
				return
			}

			if got.Title != tt.want.Title {
				t.Errorf("タイトルが一致しません。取得値 = %v, 期待値 = %v", got.Title, tt.want.Title)
			}

			if got.Description != tt.want.Description {
				t.Errorf("説明が一致しません。取得値 = %v, 期待値 = %v", got.Description, tt.want.Description)
			}

			if got.IsCompleted != tt.want.IsCompleted {
				t.Errorf("完了状態が一致しません。取得値 = %v, 期待値 = %v", got.IsCompleted, tt.want.IsCompleted)
			}
		})
	}
}

// TestUpdateTodoRequest_ApplyToEntity は更新リクエストの適用をテストします
func TestUpdateTodoRequest_ApplyToEntity(t *testing.T) {
	tests := []struct {
		name     string
		request  UpdateTodoRequest
		original *entity.Todo
		want     *entity.Todo
	}{
		{
			name: "タイトルのみ更新",
			request: UpdateTodoRequest{
				Title: stringPtr("更新されたタイトル"),
			},
			original: &entity.Todo{
				ID:          1,
				Title:       "元のタイトル",
				Description: "元の説明",
				IsCompleted: false,
			},
			want: &entity.Todo{
				ID:          1,
				Title:       "更新されたタイトル",
				Description: "元の説明",
				IsCompleted: false,
			},
		},
		{
			name: "全フィールド更新",
			request: UpdateTodoRequest{
				Title:       stringPtr("新しいタイトル"),
				Description: stringPtr("新しい説明"),
				IsCompleted: boolPtr(true),
			},
			original: &entity.Todo{
				ID:          1,
				Title:       "元のタイトル",
				Description: "元の説明",
				IsCompleted: false,
			},
			want: &entity.Todo{
				ID:          1,
				Title:       "新しいタイトル",
				Description: "新しい説明",
				IsCompleted: true,
			},
		},
		{
			name:     "何も更新しない",
			request:  UpdateTodoRequest{},
			original: &entity.Todo{
				ID:          1,
				Title:       "元のタイトル",
				Description: "元の説明",
				IsCompleted: false,
			},
			want: &entity.Todo{
				ID:          1,
				Title:       "元のタイトル",
				Description: "元の説明",
				IsCompleted: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// オリジナルをコピー
			got := *tt.original
			
			// リクエストを適用
			tt.request.ApplyToEntity(&got)

			if got.ID != tt.want.ID {
				t.Errorf("IDが一致しません。取得値 = %v, 期待値 = %v", got.ID, tt.want.ID)
			}

			if got.Title != tt.want.Title {
				t.Errorf("タイトルが一致しません。取得値 = %v, 期待値 = %v", got.Title, tt.want.Title)
			}

			if got.Description != tt.want.Description {
				t.Errorf("説明が一致しません。取得値 = %v, 期待値 = %v", got.Description, tt.want.Description)
			}

			if got.IsCompleted != tt.want.IsCompleted {
				t.Errorf("完了状態が一致しません。取得値 = %v, 期待値 = %v", got.IsCompleted, tt.want.IsCompleted)
			}
		})
	}
}

// TestToTodoResponse はエンティティからレスポンスへの変換をテストします
func TestToTodoResponse(t *testing.T) {
	// テスト用の時刻を固定
	fixedTime := time.Date(2023, 12, 25, 15, 30, 45, 0, time.UTC)

	tests := []struct {
		name   string
		entity *entity.Todo
		want   TodoResponse
	}{
		{
			name: "完全なエンティティ変換",
			entity: &entity.Todo{
				ID:          1,
				Title:       "テストタスク",
				Description: "説明文",
				IsCompleted: true,
				CreatedAt:   fixedTime,
				UpdatedAt:   fixedTime.Add(1 * time.Hour),
			},
			want: TodoResponse{
				ID:          1,
				Title:       "テストタスク",
				Description: "説明文",
				IsCompleted: true,
				CreatedAt:   fixedTime,
				UpdatedAt:   fixedTime.Add(1 * time.Hour),
			},
		},
		{
			name: "未完了タスクの変換",
			entity: &entity.Todo{
				ID:          2,
				Title:       "未完了タスク",
				Description: "",
				IsCompleted: false,
				CreatedAt:   fixedTime,
				UpdatedAt:   fixedTime,
			},
			want: TodoResponse{
				ID:          2,
				Title:       "未完了タスク",
				Description: "",
				IsCompleted: false,
				CreatedAt:   fixedTime,
				UpdatedAt:   fixedTime,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToTodoResponse(tt.entity)

			if got.ID != tt.want.ID {
				t.Errorf("ID = %v, 期待値 = %v", got.ID, tt.want.ID)
			}

			if got.Title != tt.want.Title {
				t.Errorf("タイトル = %v, 期待値 = %v", got.Title, tt.want.Title)
			}

			if got.Description != tt.want.Description {
				t.Errorf("説明 = %v, 期待値 = %v", got.Description, tt.want.Description)
			}

			if got.IsCompleted != tt.want.IsCompleted {
				t.Errorf("完了状態 = %v, 期待値 = %v", got.IsCompleted, tt.want.IsCompleted)
			}

			if !got.CreatedAt.Equal(tt.want.CreatedAt) {
				t.Errorf("作成日時 = %v, 期待値 = %v", got.CreatedAt, tt.want.CreatedAt)
			}

			if !got.UpdatedAt.Equal(tt.want.UpdatedAt) {
				t.Errorf("更新日時 = %v, 期待値 = %v", got.UpdatedAt, tt.want.UpdatedAt)
			}
		})
	}
}

// TestTodoResponse_JSONSerialization はJSONシリアライゼーションをテストします
func TestTodoResponse_JSONSerialization(t *testing.T) {
	fixedTime := time.Date(2023, 12, 25, 15, 30, 45, 0, time.UTC)
	
	response := TodoResponse{
		ID:          1,
		Title:       "テストタスク",
		Description: "説明文",
		IsCompleted: true,
		CreatedAt:   fixedTime,
		UpdatedAt:   fixedTime,
	}

	// JSONにシリアライズ
	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("JSONシリアライゼーションに失敗: %v", err)
	}

	// 期待されるJSONフィールドが含まれていることを確認
	jsonStr := string(jsonData)
	expectedFields := []string{
		`"id":1`,
		`"title":"テストタスク"`,
		`"description":"説明文"`,
		`"is_completed":true`,
		`"created_at":"2023-12-25T15:30:45Z"`,
		`"updated_at":"2023-12-25T15:30:45Z"`,
	}

	for _, field := range expectedFields {
		if !contains(jsonStr, field) {
			t.Errorf("JSONに期待されるフィールドが含まれていません: %s", field)
		}
	}

	// JSONからデシリアライズして元に戻る確認
	var deserialized TodoResponse
	err = json.Unmarshal(jsonData, &deserialized)
	if err != nil {
		t.Fatalf("JSONデシリアライゼーションに失敗: %v", err)
	}

	if deserialized.ID != response.ID {
		t.Errorf("デシリアライズ後のID = %v, 期待値 = %v", deserialized.ID, response.ID)
	}

	if deserialized.Title != response.Title {
		t.Errorf("デシリアライズ後のタイトル = %v, 期待値 = %v", deserialized.Title, response.Title)
	}
}

// TestCreateTodoRequest_JSONDeserialization はリクエストのJSONデシリアライゼーションをテストします
func TestCreateTodoRequest_JSONDeserialization(t *testing.T) {
	tests := []struct {
		name     string
		jsonStr  string
		want     CreateTodoRequest
		wantErr  bool
	}{
		{
			name:    "正常なJSON",
			jsonStr: `{"title":"テストタスク","description":"説明文"}`,
			want: CreateTodoRequest{
				Title:       "テストタスク",
				Description: "説明文",
			},
			wantErr: false,
		},
		{
			name:    "説明なしJSON",
			jsonStr: `{"title":"タイトルのみ"}`,
			want: CreateTodoRequest{
				Title:       "タイトルのみ",
				Description: "",
			},
			wantErr: false,
		},
		{
			name:     "不正なJSON",
			jsonStr:  `{"title": invalid json}`,
			want:     CreateTodoRequest{},
			wantErr:  true,
		},
		{
			name:     "空のJSON",
			jsonStr:  `{}`,
			want:     CreateTodoRequest{},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got CreateTodoRequest
			err := json.Unmarshal([]byte(tt.jsonStr), &got)

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

			if got.Title != tt.want.Title {
				t.Errorf("タイトル = %v, 期待値 = %v", got.Title, tt.want.Title)
			}

			if got.Description != tt.want.Description {
				t.Errorf("説明 = %v, 期待値 = %v", got.Description, tt.want.Description)
			}
		})
	}
}

// TestErrorResponse はエラーレスポンスの構造をテストします
func TestErrorResponse(t *testing.T) {
	tests := []struct {
		name    string
		error   string
		details interface{}
	}{
		{
			name:    "基本的なエラーレスポンス",
			error:   "エラーが発生しました",
			details: nil,
		},
		{
			name:    "詳細付きエラーレスポンス",
			error:   "バリデーションエラー",
			details: "タイトルが空です",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errorResp := ErrorResponse{
				Error:   tt.error,
				Details: tt.details,
			}

			// JSONシリアライゼーションの確認
			jsonData, err := json.Marshal(errorResp)
			if err != nil {
				t.Errorf("JSONシリアライゼーションに失敗: %v", err)
			}

			// 必要なフィールドが含まれていることを確認
			jsonStr := string(jsonData)
			if !contains(jsonStr, `"error":"`+tt.error+`"`) {
				t.Error("JSONにエラーフィールドが含まれていません")
			}

			if tt.details != nil {
				if details, ok := tt.details.(string); ok && !contains(jsonStr, `"details":"`+details+`"`) {
					t.Error("JSONに詳細フィールドが含まれていません")
				}
			}
		})
	}
}

// stringPtr はstring値のポインタを返すヘルパー関数です
func stringPtr(s string) *string {
	return &s
}

// boolPtr はbool値のポインタを返すヘルパー関数です
func boolPtr(b bool) *bool {
	return &b
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

// contains は文字列に指定のサブ文字列が含まれているかチェックするヘルパー関数です
func contains(s, substr string) bool {
	return len(s) >= len(substr) && indexOf(s, substr) >= 0
}

// indexOf は文字列内での指定サブ文字列のインデックスを返すヘルパー関数です
func indexOf(s, substr string) int {
	if len(substr) == 0 {
		return 0
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// 標準パッケージでのDTOテストの学習ポイント：
//
// 1. データ変換ロジックのテスト：
//    - エンティティ ↔ DTO の変換確認
//    - バリデーションロジックの検証
//    - エラーケースの適切な処理
//
// 2. JSONハンドリングのテスト：
//    - encoding/json パッケージの使用
//    - シリアライゼーション/デシリアライゼーション
//    - JSONタグの動作確認
//
// 3. 時刻フォーマットのテスト：
//    - RFC3339 フォーマットでの時刻変換
//    - 異なる時刻でのフォーマット確認
//
// 4. バリデーションテスト：
//    - 入力値の妥当性確認
//    - エラーメッセージの適切性
//    - 境界値の処理
//
// 5. API契約のテスト：
//    - リクエスト/レスポンス構造の確認
//    - フィールド名とタイプの検証
//    - エラーレスポンスの一貫性