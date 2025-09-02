package dto

import (
	"time"

	"todoapp-api-golang/internal/domain/entity"
)

// TodoResponse はTodo情報をクライアントに返すためのレスポンスDTOです
// レスポンスDTOの役割：
// 1. 外部に公開する情報の制御（セキュリティ）
// 2. クライアントに最適化されたデータ構造の提供
// 3. APIのバージョニング対応
// 4. 内部実装の隠蔽（エンティティの変更がAPIに影響しないようにする）
type TodoResponse struct {
	// ID はTodoの一意識別子
	ID int `json:"id"`

	// Title はTodoのタイトル
	Title string `json:"title"`

	// Description はTodoの詳細説明
	Description string `json:"description"`

	// IsCompleted はTodoの完了状態
	IsCompleted bool `json:"is_completed"`

	// CreatedAt は作成日時（RFC3339形式でJSONシリアライズ）
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt は最終更新日時
	UpdatedAt time.Time `json:"updated_at"`
}

// TodoListResponse はTodo一覧取得時のレスポンスDTOです
// 将来的なページング情報なども含められる構造にしています
type TodoListResponse struct {
	// Todos はTodoのリスト
	Todos []TodoResponse `json:"todos"`

	// Meta はメタ情報（ページング等）
	Meta ListMetaResponse `json:"meta"`
}

// ListMetaResponse は一覧取得時のメタ情報を表すDTOです
// ページング情報や総件数など、一覧表示に必要な付加情報を含みます
type ListMetaResponse struct {
	// Total は総件数
	Total int `json:"total"`

	// Page は現在のページ番号
	Page int `json:"page"`

	// Limit は1ページあたりの表示件数
	Limit int `json:"limit"`

	// TotalPages は総ページ数
	TotalPages int `json:"total_pages"`
}

// ErrorResponse はエラー発生時のレスポンスDTOです
// 統一的なエラーレスポンス形式を提供します
type ErrorResponse struct {
	// Error はエラーメッセージ
	Error string `json:"error"`

	// Code はエラーコード（任意、アプリケーション固有のコード）
	Code string `json:"code,omitempty"`

	// Details は詳細情報（バリデーションエラー等）
	Details interface{} `json:"details,omitempty"`
}

// ValidationErrorResponse はバリデーションエラー専用のレスポンスDTOです
type ValidationErrorResponse struct {
	// Error は基本エラーメッセージ
	Error string `json:"error"`

	// ValidationErrors はフィールド別のバリデーションエラー
	ValidationErrors []FieldError `json:"validation_errors"`
}

// FieldError はフィールド単位のバリデーションエラー情報です
type FieldError struct {
	// Field はエラーが発生したフィールド名
	Field string `json:"field"`

	// Message はエラーメッセージ
	Message string `json:"message"`

	// Value は入力された値（セキュリティ上問題ない場合のみ）
	Value interface{} `json:"value,omitempty"`
}

// --- 変換関数（Mapper functions） ---

// ToTodoResponse はEntityをResponseDTOに変換します
// エンティティ → レスポンスDTO の変換ロジック
func ToTodoResponse(todo *entity.Todo) TodoResponse {
	return TodoResponse{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
		IsCompleted: todo.IsCompleted,
		CreatedAt:   todo.CreatedAt,
		UpdatedAt:   todo.UpdatedAt,
	}
}

// ToTodoListResponse はEntity配列をResponseDTOに変換します
func ToTodoListResponse(todos []*entity.Todo, page, limit, total int) TodoListResponse {
	// Entity配列を Response配列に変換
	todoResponses := make([]TodoResponse, len(todos))
	for i, todo := range todos {
		todoResponses[i] = ToTodoResponse(todo)
	}

	// ページ数の計算
	totalPages := total / limit
	if total%limit != 0 {
		totalPages++
	}

	return TodoListResponse{
		Todos: todoResponses,
		Meta: ListMetaResponse{
			Total:      total,
			Page:       page,
			Limit:      limit,
			TotalPages: totalPages,
		},
	}
}

// ToEntity はリクエストDTOをEntityに変換します（Create用）
func (req CreateTodoRequest) ToEntity() *entity.Todo {
	return &entity.Todo{
		Title:       req.Title,
		Description: req.Description,
		// IsCompleted は新規作成時は常にfalse（デフォルト値）
		IsCompleted: false,
	}
}

// ApplyToEntity は更新リクエストDTOを既存Entityに適用します（Update用）
// nil チェックを行い、送信されたフィールドのみを更新します
func (req UpdateTodoRequest) ApplyToEntity(todo *entity.Todo) {
	// タイトルが送信された場合のみ更新
	if req.Title != nil {
		todo.Title = *req.Title
	}

	// 説明が送信された場合のみ更新
	if req.Description != nil {
		todo.Description = *req.Description
	}

	// 完了状態が送信された場合のみ更新
	if req.IsCompleted != nil {
		todo.IsCompleted = *req.IsCompleted
	}
}

// DTOパターンの利点：
// 1. セキュリティ: 内部IDやパスワードなど、外部に公開したくない情報を隠蔽
// 2. 進化性: APIの変更を内部実装の変更から分離
// 3. 最適化: クライアント要件に合わせたデータ構造の提供
// 4. バリデーション: 入力値の検証と制御
// 5. ドキュメント化: APIドキュメント生成のための明確な構造定義
