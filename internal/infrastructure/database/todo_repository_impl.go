package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"todoapp-api-golang/internal/domain/entity"
	"todoapp-api-golang/internal/domain/repository"
)

// todoRepositoryImpl は標準のdatabase/sqlパッケージを使用した
// TodoRepositoryインターフェースの具体的実装です
//
// database/sql パッケージの学習ポイント：
// 1. SQL文の直接記述とプリペアードステートメントの理解
// 2. sql.DB, sql.Tx, sql.Rows の適切な使用方法
// 3. NULL値の扱いとsql.NullStringなどの活用
// 4. トランザクション処理の実装
// 5. コネクションプールの仕組み
type todoRepositoryImpl struct {
	// db は標準のdatabase/sqlのDB接続
	// *sql.DB はコネクションプールを管理し、並行安全
	db *sql.DB
}

// NewTodoRepository はtodoRepositoryImplのコンストラクタです
// 標準パッケージを使った依存性注入の実装
func NewTodoRepository(db *sql.DB) repository.TodoRepository {
	return &todoRepositoryImpl{
		db: db,
	}
}

// Create は新しいTodoをデータベースに保存します
// 標準パッケージを使ったINSERT操作の学習
func (r *todoRepositoryImpl) Create(ctx context.Context, todo *entity.Todo) (*entity.Todo, error) {
	// 1. INSERT用のSQL文を定義
	// プリペアードステートメント（?プレースホルダー）でSQLインジェクション対策
	// created_at, updated_atは現在時刻、is_completedはfalseで固定
	query := `
		INSERT INTO todos (title, description, is_completed, created_at, updated_at)
		VALUES (?, ?, false, NOW(), NOW())
	`

	// 2. コンテキスト付きでSQL実行
	// ExecContext はINSERT/UPDATE/DELETE用（結果行を返さない）
	result, err := r.db.ExecContext(ctx, query, todo.Title, todo.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to insert todo: %w", err)
	}

	// 3. 自動生成されたIDを取得
	// LastInsertId() でAUTO_INCREMENTの値を取得
	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get inserted ID: %w", err)
	}

	// 4. IDを設定して作成済みTodoを返却
	todo.ID = int(id)
	todo.IsCompleted = false
	todo.CreatedAt = time.Now()
	todo.UpdatedAt = time.Now()

	return todo, nil
}

// GetByID は主キーによる1件取得を行います
// 標準パッケージを使ったSELECT操作とNULL値の扱い方を学習
func (r *todoRepositoryImpl) GetByID(ctx context.Context, id int) (*entity.Todo, error) {
	// 1. SELECT用のSQL文を定義
	query := `
		SELECT id, title, description, is_completed, created_at, updated_at
		FROM todos
		WHERE id = ?
	`

	// 2. 1行取得用のQueryRowContext を使用
	row := r.db.QueryRowContext(ctx, query, id)

	// 3. 結果を構造体にスキャン
	var todo entity.Todo
	err := row.Scan(
		&todo.ID,
		&todo.Title,
		&todo.Description,
		&todo.IsCompleted,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)

	if err != nil {
		// sql.ErrNoRows は「データが見つからない」を示す標準エラー
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("todo not found")
		}
		return nil, fmt.Errorf("failed to scan todo: %w", err)
	}

	return &todo, nil
}

// GetAll は全件取得を行います
// 標準パッケージを使った複数行取得とRowsの適切な処理を学習
func (r *todoRepositoryImpl) GetAll(ctx context.Context) ([]*entity.Todo, error) {
	// 1. SELECT用のSQL文（作成日時の降順でソート）
	query := `
		SELECT id, title, description, is_completed, created_at, updated_at
		FROM todos
		ORDER BY created_at DESC
	`

	// 2. 複数行取得用のQueryContext を使用
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query todos: %w", err)
	}
	
	// 3. 重要：rowsは必ずClose()する（deferで確実に実行）
	// リソースリーク防止のための必須処理
	defer rows.Close()

	// 4. 結果を格納するスライスを初期化
	var todos []*entity.Todo

	// 5. rows.Next()でループして全ての行を処理
	for rows.Next() {
		var todo entity.Todo
		
		// 各行をScanして構造体に格納
		err := rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Description,
			&todo.IsCompleted,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan todo row: %w", err)
		}
		
		// スライスに追加
		todos = append(todos, &todo)
	}

	// 6. ループ終了後にエラーチェック
	// ネットワークエラーなどでループが中断された場合を検出
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return todos, nil
}

// Update は既存レコードの更新を行います
// 標準パッケージを使ったUPDATE操作と影響行数の確認を学習
func (r *todoRepositoryImpl) Update(ctx context.Context, todo *entity.Todo) (*entity.Todo, error) {
	// 1. UPDATE用のSQL文を定義
	// updated_at は現在時刻で自動更新
	query := `
		UPDATE todos
		SET title = ?, description = ?, is_completed = ?, updated_at = NOW()
		WHERE id = ?
	`

	// 2. UPDATE実行
	result, err := r.db.ExecContext(ctx, query,
		todo.Title,
		todo.Description,
		todo.IsCompleted,
		todo.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update todo: %w", err)
	}

	// 3. 影響を受けた行数を確認
	// RowsAffected()で実際に更新された行数を取得
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	// 4. 行が更新されなかった場合はエラー
	if rowsAffected == 0 {
		return nil, errors.New("todo not found")
	}

	// 5. 更新後のデータを取得して返却
	// updated_at を最新の値にするため再取得
	return r.GetByID(ctx, todo.ID)
}

// Delete は主キーによる削除を行います
// 標準パッケージを使ったDELETE操作を学習
func (r *todoRepositoryImpl) Delete(ctx context.Context, id int) error {
	// 1. DELETE用のSQL文を定義
	query := `DELETE FROM todos WHERE id = ?`

	// 2. DELETE実行
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}

	// 3. 影響を受けた行数を確認
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	// 4. 削除された行がない場合はエラー
	if rowsAffected == 0 {
		return errors.New("todo not found")
	}

	return nil
}

// GetByCompleteStatus は完了状態による検索を行います（将来の拡張用）
// WHERE句を使った条件検索の学習
func (r *todoRepositoryImpl) GetByCompleteStatus(ctx context.Context, isCompleted bool) ([]*entity.Todo, error) {
	query := `
		SELECT id, title, description, is_completed, created_at, updated_at
		FROM todos
		WHERE is_completed = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, isCompleted)
	if err != nil {
		return nil, fmt.Errorf("failed to query todos by status: %w", err)
	}
	defer rows.Close()

	var todos []*entity.Todo
	for rows.Next() {
		var todo entity.Todo
		err := rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Description,
			&todo.IsCompleted,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan todo row: %w", err)
		}
		todos = append(todos, &todo)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return todos, nil
}

// GetWithPagination はページング機能付きの取得を行います（将来の拡張用）
// LIMIT、OFFSET句を使った標準的なページング実装を学習
func (r *todoRepositoryImpl) GetWithPagination(ctx context.Context, offset, limit int) ([]*entity.Todo, int64, error) {
	// 1. 総件数を取得するSQL
	countQuery := `SELECT COUNT(*) FROM todos`
	var total int64
	
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count todos: %w", err)
	}

	// 2. ページング付きでデータを取得するSQL
	dataQuery := `
		SELECT id, title, description, is_completed, created_at, updated_at
		FROM todos
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, dataQuery, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query todos with pagination: %w", err)
	}
	defer rows.Close()

	var todos []*entity.Todo
	for rows.Next() {
		var todo entity.Todo
		err := rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Description,
			&todo.IsCompleted,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan todo row: %w", err)
		}
		todos = append(todos, &todo)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error during rows iteration: %w", err)
	}

	return todos, total, nil
}

// database/sql パッケージの学習ポイント：
//
// 1. コネクション管理：
//    - *sql.DB は複数のgoroutineから安全に使用可能
//    - 内部的にコネクションプールを管理
//    - Close()でリソース解放が必要
//
// 2. SQL実行方法：
//    - QueryContext(): SELECT用（複数行）
//    - QueryRowContext(): SELECT用（1行）
//    - ExecContext(): INSERT/UPDATE/DELETE用
//
// 3. プリペアードステートメント：
//    - ?プレースホルダーでSQLインジェクション防止
//    - データベースでのクエリ最適化
//
// 4. エラーハンドリング：
//    - sql.ErrNoRows での「データなし」判定
//    - fmt.Errorf でのエラーラッピング
//
// 5. リソース管理：
//    - sql.Rows の defer Close()
//    - rows.Err() でのループエラー確認
//
// 6. トランザクション（今回は省略、将来拡張可能）：
//    - db.BeginTx() でのトランザクション開始
//    - tx.Commit() / tx.Rollback() での制御