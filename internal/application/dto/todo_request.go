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
	// `binding:"required"` : Ginフレームワークのバリデーションタグ（必須チェック）
	// `validate:"required,min=1,max=100"` : validator パッケージのバリデーション
	Title string `json:"title" binding:"required" validate:"required,min=1,max=100"`

	// Description はTodoの詳細説明（任意項目）
	// バリデーションは設定していませんが、長さ制限を設ける場合は
	// `validate:"max=500"` などを追加できます
	Description string `json:"description" validate:"max=500"`
}

// UpdateTodoRequest はTodo更新時のHTTPリクエストボディを表すDTOです
// 作成時とは異なり、全てのフィールドが任意更新可能な設計にしています
// （部分更新：PATCHメソッドの考え方）
type UpdateTodoRequest struct {
	// Title の更新（任意）
	// ポインタ型 (*string) を使用することで、フィールドが送信されたかどうかを判別可能
	// nil の場合は更新しない、値がある場合は更新する
	Title *string `json:"title,omitempty" validate:"omitempty,min=1,max=100"`

	// Description の更新（任意）
	Description *string `json:"description,omitempty" validate:"omitempty,max=500"`

	// IsCompleted の更新（任意）
	// bool のポインタ型で、完了状態の変更を任意にします
	IsCompleted *bool `json:"is_completed,omitempty"`
}

// CompleteTodoRequest はTodo完了/未完了切り替え専用のリクエストです
// シンプルなアクション用のDTOとして定義
type CompleteTodoRequest struct {
	// IsCompleted で完了状態を指定
	// true: 完了, false: 未完了
	IsCompleted bool `json:"is_completed" binding:"required"`
}

// TodoListRequest はTodo一覧取得時のクエリパラメータを表すDTOです
// 将来的な拡張（ページング、フィルタリング、ソート）を想定した構造
type TodoListRequest struct {
	// ページング関連（将来的な拡張用）
	// Page は取得するページ番号（1から開始）
	Page int `form:"page" validate:"min=1"`

	// Limit は1ページあたりの取得件数
	Limit int `form:"limit" validate:"min=1,max=100"`

	// フィルタリング関連（将来的な拡張用）
	// IsCompleted で完了状態によるフィルタ（任意）
	// nil の場合は全て、true/false で絞り込み
	IsCompleted *bool `form:"is_completed"`

	// ソート関連（将来的な拡張用）
	// SortBy はソートするフィールド名
	SortBy string `form:"sort_by" validate:"omitempty,oneof=id title created_at updated_at"`

	// SortOrder はソート順序（asc/desc）
	SortOrder string `form:"sort_order" validate:"omitempty,oneof=asc desc"`
}

// DTOのバリデーションルール解説：
//
// binding:"required" - Ginのバリデーション（必須）
// validate:"required" - validator パッケージのバリデーション（必須）
// validate:"min=1,max=100" - 最小1文字、最大100文字
// validate:"omitempty" - 空の場合はバリデーションをスキップ
// validate:"oneof=asc desc" - 指定した値のいずれかのみ許可
// json:"field_name,omitempty" - 空の場合はJSONに含めない
// form:"field_name" - URLクエリパラメータやフォームデータのキー名