package service

import (
	"context"
	"todoapp-api-golang/internal/domain/entity"
)

// TodoServiceInterface は Todo サービスのインターフェースです
// テスタビリティ向上のため、ハンドラー層のテストでモック実装を使用できます
// 標準パッケージでのインターフェース設計の学習ポイント：
// 1. 小さく、焦点を絞ったインターフェース
// 2. 具象型ではなくインターフェースに依存する設計
// 3. テスト時のモック実装による依存関係の分離
type TodoServiceInterface interface {
	// CreateTodo は新しいTodoを作成します
	CreateTodo(ctx context.Context, todo *entity.Todo) (*entity.Todo, error)

	// GetTodoByID は指定されたIDのTodoを取得します
	GetTodoByID(ctx context.Context, id int) (*entity.Todo, error)

	// GetAllTodos は全てのTodoを取得します
	GetAllTodos(ctx context.Context) ([]*entity.Todo, error)

	// UpdateTodo は既存のTodoを更新します
	UpdateTodo(ctx context.Context, todo *entity.Todo) (*entity.Todo, error)

	// DeleteTodo は指定されたIDのTodoを削除します
	DeleteTodo(ctx context.Context, id int) error

	// CompleteTodo はTodoを完了状態にします
	CompleteTodo(ctx context.Context, id int) (*entity.Todo, error)

	// IncompleteTodo はTodoを未完了状態にします
	IncompleteTodo(ctx context.Context, id int) (*entity.Todo, error)
}

// コンパイル時インターフェース実装確認
// この行により、TodoService が TodoServiceInterface を実装していることを
// コンパイル時に確認できます
var _ TodoServiceInterface = (*TodoService)(nil)