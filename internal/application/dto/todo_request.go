package dto

// CreateTodoRequest はTodo作成時のHTTPリクエストボディを表すDTO（Data Transfer Object）です
// DTOの役割：
// 1. HTTPリクエスト/レスポンスの構造を定義
// 2. 外部システム（クライアント）とのデータ交換フォーマット
// 3. ドメインエンティティとの変換（マッピング）
// 4. 入力値の基本的なバリデーション
type CreateTodoRequest struct {
	// Title はTodoのタイトル（必須項目）
	// `json:"title"` : JSONキー名を指定（Goの命名規則と異なる場合に使用）
	// 標準パッケージではバリデーションタグは使用せず、手動バリデーションを行います
	Title string `json:"title"`

	// Description はTodoの詳細説明（任意項目）
	// 長さ制限などのバリデーションは実装層で手動実装します
	Description string `json:"description"`
}

// UpdateTodoRequest はTodo更新時のHTTPリクエストボディを表すDTOです
// 作成時とは異なり、全てのフィールドが任意更新可能な設計にしています
// （部分更新：PATCHメソッドの考え方）
type UpdateTodoRequest struct {
	// Title の更新（任意）
	// ポインタ型 (*string) を使用することで、フィールドが送信されたかどうかを判別可能
	// nil の場合は更新しない、値がある場合は更新する
	Title *string `json:"title,omitempty"`

	// Description の更新（任意）
	// 長さ制限などのバリデーションは実装層で行います
	Description *string `json:"description,omitempty"`

	// IsCompleted の更新（任意）
	// bool のポインタ型で、完了状態の変更を任意にします
	IsCompleted *bool `json:"is_completed,omitempty"`
}

// CompleteTodoRequest はTodo完了/未完了切り替え専用のリクエストです
// シンプルなアクション用のDTOとして定義
type CompleteTodoRequest struct {
	// IsCompleted で完了状態を指定
	// true: 完了, false: 未完了
	// バリデーションは実装層で手動実装します
	IsCompleted bool `json:"is_completed"`
}

// TodoListRequest はTodo一覧取得時のクエリパラメータを表すDTOです
// 将来的な拡張（ページング、フィルタリング、ソート）を想定した構造
type TodoListRequest struct {
	// ページング関連（将来的な拡張用）
	// Page は取得するページ番号（1から開始）
	Page int `json:"page"`

	// Limit は1ページあたりの取得件数
	Limit int `json:"limit"`

	// フィルタリング関連（将来的な拡張用）
	// IsCompleted で完了状態によるフィルタ（任意）
	// nil の場合は全て、true/false で絞り込み
	IsCompleted *bool `json:"is_completed"`

	// ソート関連（将来的な拡張用）
	// SortBy はソートするフィールド名
	SortBy string `json:"sort_by"`

	// SortOrder はソート順序（asc/desc）
	SortOrder string `json:"sort_order"`
}

// 標準パッケージでのDTO設計の学習ポイント：
//
// 1. 構造体タグ：
//    - json:"field_name" - JSONキー名の指定
//    - json:"field_name,omitempty" - 空値の場合はJSONに含めない
//
// 2. バリデーション：
//    - タグベースのバリデーションは使用しない
//    - 実装層（ハンドラー）で手動バリデーションを実装
//    - ビジネスルールに応じた柔軟な検証が可能
//
// 3. ポインタ型の活用：
//    - *string, *bool でフィールドの送信有無を判別
//    - nil = 未送信、値あり = 送信済み
//    - 部分更新（PATCH）の実装に有効
//
// 4. 将来拡張性：
//    - フィルタリング、ソート、ページングを考慮した設計
//    - 標準パッケージでも十分な機能を提供可能
