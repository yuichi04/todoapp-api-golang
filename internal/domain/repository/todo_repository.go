package repository

import (
	"context"

	"todoapp-api-golang/internal/domain/entity"
)

// TodoRepository はTodoエンティティのデータアクセスを抽象化するインターフェースです
// Clean Architectureでは、ドメイン層でインターフェースを定義し、
// インフラストラクチャ層で具体的な実装を行います（依存関係逆転の原則）
//
// インターフェースを使う利点：
// 1. テスタビリティ: モックを作成してユニットテストが可能
// 2. 疎結合: データベースの変更（PostgreSQL → MySQL等）が容易
// 3. 保守性: ビジネスロジックとデータアクセスロジックの分離
type TodoRepository interface {
	// Create は新しいTodoを作成します
	// 引数:
	//   - ctx: リクエストのキャンセルやタイムアウト制御に使用
	//   - todo: 作成するTodoエンティティ（IDは自動生成される）
	// 戻り値:
	//   - *entity.Todo: 作成されたTodo（IDが設定済み）
	//   - error: エラーが発生した場合のエラー情報
	Create(ctx context.Context, todo *entity.Todo) (*entity.Todo, error)

	// GetByID は指定されたIDのTodoを1件取得します
	// 引数:
	//   - ctx: コンテキスト（リクエストライフサイクル管理）
	//   - id: 取得したいTodoのID
	// 戻り値:
	//   - *entity.Todo: 見つかったTodoエンティティ
	//   - error: Todo が見つからない場合やDBエラーの場合
	GetByID(ctx context.Context, id int) (*entity.Todo, error)

	// GetAll は全てのTodoを取得します
	// 実際のアプリケーションでは、ページング（limit/offset）や
	// フィルタリング、ソート機能を追加することが多いです
	// 引数:
	//   - ctx: コンテキスト
	// 戻り値:
	//   - []*entity.Todo: Todoのスライス（配列）
	//   - error: DBエラーの場合
	GetAll(ctx context.Context) ([]*entity.Todo, error)

	// Update は既存のTodoを更新します
	// 引数:
	//   - ctx: コンテキスト
	//   - todo: 更新するTodoエンティティ（IDは必須）
	// 戻り値:
	//   - *entity.Todo: 更新されたTodo
	//   - error: Todo が見つからない場合やDBエラーの場合
	Update(ctx context.Context, todo *entity.Todo) (*entity.Todo, error)

	// Delete は指定されたIDのTodoを削除します
	// 引数:
	//   - ctx: コンテキスト
	//   - id: 削除するTodoのID
	// 戻り値:
	//   - error: Todo が見つからない場合やDBエラーの場合
	// Note: 戻り値はerrorのみです（削除されたレコードの情報は不要なため）
	Delete(ctx context.Context, id int) error
}

// メモ：なぜcontextパッケージを使うのか？
// 1. タイムアウト制御: 長時間のDB処理をキャンセルできる
// 2. キャンセル処理: ユーザーがリクエストをキャンセルした時の伝播
// 3. 値の伝播: リクエストIDやユーザー情報の受け渡し
// 4. ログトレース: 分散システムでのリクエスト追跡