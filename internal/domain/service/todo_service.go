package service

import (
	"context"
	"errors"
	"fmt"

	"todoapp-api-golang/internal/domain/entity"
	"todoapp-api-golang/internal/domain/repository"
)

// TodoService はTodoに関するビジネスロジックを管理するドメインサービスです
// ドメインサービスの役割：
// 1. 複数のエンティティやリポジトリを組み合わせた業務処理
// 2. ビジネスルールの実装と検証
// 3. トランザクション境界の定義
// 4. ドメイン固有のエラーハンドリング
type TodoService struct {
	// todoRepo はTodoRepositoryインターフェースを通じてデータアクセスを行います
	// インターフェース経由で実装することで、依存関係を逆転させています
	// （ドメイン層がインフラ層に依存しない設計）
	todoRepo repository.TodoRepository
}

// NewTodoService はTodoServiceのコンストラクタ関数です
// 依存性注入（Dependency Injection）のパターンを使用しています
// 引数:
//   - todoRepo: TodoRepositoryインターフェースの実装
//
// 戻り値:
//   - *TodoService: 初期化されたTodoService
func NewTodoService(todoRepo repository.TodoRepository) *TodoService {
	return &TodoService{
		todoRepo: todoRepo,
	}
}

// CreateTodo は新しいTodoを作成するビジネスロジックです
// ここではドメインルールの検証を行った後、リポジトリに処理を委譲します
func (s *TodoService) CreateTodo(ctx context.Context, todo *entity.Todo) (*entity.Todo, error) {
	// 1. 入力値のドメインレベルバリデーション
	// エンティティのIsValid()メソッドでビジネスルールをチェック
	if !todo.IsValid() {
		return nil, errors.New("todo validation failed: title is required and must be 100 characters or less")
	}

	// 2. 追加のビジネスルールチェック（例：タイトルの重複チェックなど）
	// 実際のアプリケーションでは、「同じタイトルのTodoは作成できない」
	// などのルールがある場合があります

	// 3. リポジトリを通じてデータ永続化
	createdTodo, err := s.todoRepo.Create(ctx, todo)
	if err != nil {
		// エラーラッピング：下位層のエラーに追加情報を付与
		return nil, fmt.Errorf("failed to create todo: %w", err)
	}

	return createdTodo, nil
}

// GetTodoByID は指定されたIDのTodoを取得します
func (s *TodoService) GetTodoByID(ctx context.Context, id int) (*entity.Todo, error) {
	// 1. 入力値の基本バリデーション
	if id <= 0 {
		return nil, errors.New("invalid todo ID: must be greater than 0")
	}

	// 2. リポジトリから取得
	todo, err := s.todoRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get todo with ID %d: %w", id, err)
	}

	return todo, nil
}

// GetAllTodos は全てのTodoを取得します
func (s *TodoService) GetAllTodos(ctx context.Context) ([]*entity.Todo, error) {
	todos, err := s.todoRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all todos: %w", err)
	}

	// ビジネスロジック：取得したTodoに追加の処理を行う場合
	// 例：完了済みのTodoを先頭に移動、期限切れのチェックなど
	// この例では単純に取得した結果をそのまま返します

	return todos, nil
}

// UpdateTodo は既存のTodoを更新します
func (s *TodoService) UpdateTodo(ctx context.Context, todo *entity.Todo) (*entity.Todo, error) {
	// 1. 入力値バリデーション
	if todo.ID <= 0 {
		return nil, errors.New("invalid todo ID: must be greater than 0")
	}

	if !todo.IsValid() {
		return nil, errors.New("todo validation failed: title is required and must be 100 characters or less")
	}

	// 2. 存在チェック（更新前にレコードが存在するか確認）
	existingTodo, err := s.todoRepo.GetByID(ctx, todo.ID)
	if err != nil {
		return nil, fmt.Errorf("todo with ID %d not found: %w", todo.ID, err)
	}

	// 3. ビジネスルールに基づく更新制御
	// 例：「完了済みのTodoは編集できない」などのルールがある場合
	// この例では特に制約を設けていません
	_ = existingTodo // 存在チェックのみで使用

	// 4. リポジトリを通じて更新実行
	updatedTodo, err := s.todoRepo.Update(ctx, todo)
	if err != nil {
		return nil, fmt.Errorf("failed to update todo: %w", err)
	}

	return updatedTodo, nil
}

// DeleteTodo は指定されたIDのTodoを削除します
func (s *TodoService) DeleteTodo(ctx context.Context, id int) error {
	// 1. 入力値バリデーション
	if id <= 0 {
		return errors.New("invalid todo ID: must be greater than 0")
	}

	// 2. 存在チェック（削除前にレコードが存在するか確認）
	_, err := s.todoRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("todo with ID %d not found: %w", id, err)
	}

	// 3. ビジネスルールチェック
	// 例：「作成から24時間以内のTodoは削除できない」などのルール
	// この例では特に制約を設けていません

	// 4. リポジトリを通じて削除実行
	err = s.todoRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}

	return nil
}

// CompleteTodo はTodoを完了状態にする専用メソッドです
// エンティティのビジネスロジック（MarkAsCompleted）を使用した例
func (s *TodoService) CompleteTodo(ctx context.Context, id int) (*entity.Todo, error) {
	// 1. 対象のTodoを取得
	todo, err := s.todoRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("todo with ID %d not found: %w", id, err)
	}

	// 2. エンティティのビジネスロジックを使用して状態変更
	todo.MarkAsCompleted()

	// 3. 変更をデータベースに保存
	updatedTodo, err := s.todoRepo.Update(ctx, todo)
	if err != nil {
		return nil, fmt.Errorf("failed to complete todo: %w", err)
	}

	return updatedTodo, nil
}

// IncompleteTodo はTodoを未完了状態に戻す専用メソッドです
func (s *TodoService) IncompleteTodo(ctx context.Context, id int) (*entity.Todo, error) {
	// 1. 対象のTodoを取得
	todo, err := s.todoRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("todo with ID %d not found: %w", id, err)
	}

	// 2. エンティティのビジネスロジックを使用して状態変更
	todo.MarkAsIncomplete()

	// 3. 変更をデータベースに保存
	updatedTodo, err := s.todoRepo.Update(ctx, todo)
	if err != nil {
		return nil, fmt.Errorf("failed to mark todo as incomplete: %w", err)
	}

	return updatedTodo, nil
}
