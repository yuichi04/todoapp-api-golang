package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	// MySQL ドライバーをインポート
	// _（ブランクインポート）で副作用のみを利用
	// init()関数でdriver登録が実行される
	_ "github.com/go-sql-driver/mysql"

	"todoapp-api-golang/pkg/config"
)

// DatabaseManager は標準のdatabase/sqlを使用してデータベース接続を管理する構造体です
// 標準パッケージを使ったデータベース管理の学習ポイント：
// 1. sql.DB の適切な使用方法
// 2. コネクションプールの設定と管理
// 3. データベースドライバーの登録と使用
// 4. テーブル作成とマイグレーション
// 5. ヘルスチェックとDB接続の確認
type DatabaseManager struct {
	DB     *sql.DB
	config *config.Config
}

// NewDatabaseManager はDatabaseManagerのコンストラクタです
// 標準パッケージを使った依存性注入の実装
func NewDatabaseManager(cfg *config.Config) *DatabaseManager {
	return &DatabaseManager{
		config: cfg,
	}
}

// Connect はデータベースへの接続を確立します
// database/sqlパッケージを使った接続処理の学習
func (dm *DatabaseManager) Connect() error {
	// 1. データベースドライバーの確認
	if dm.config.Database.Driver != "mysql" {
		return fmt.Errorf("unsupported database driver: %s (only mysql supported in standard package version)", dm.config.Database.Driver)
	}

	// 2. データソース名（DSN）の構築
	dsn := dm.config.GetDSN()
	log.Printf("Connecting to database: %s@%s:%d/%s",
		dm.config.Database.User,
		dm.config.Database.Host,
		dm.config.Database.Port,
		dm.config.Database.Name)

	// 3. データベース接続を開く
	// sql.Open() は実際には接続せず、DB構造体を作成するだけ
	// 実際の接続は最初のクエリ実行時に行われる
	db, err := sql.Open(dm.config.Database.Driver, dsn)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	// 4. コネクションプールの設定
	// これらの設定はパフォーマンスとリソース使用量に重要な影響を与える

	// SetMaxOpenConns: 同時に開けるコネクションの最大数
	// 高い値 = 並行性向上、低い値 = DB負荷軽減
	db.SetMaxOpenConns(dm.config.Database.MaxOpenConns)

	// SetMaxIdleConns: アイドル状態で保持するコネクションの最大数
	// アイドル接続があることで新しいリクエストへの応答が高速化
	db.SetMaxIdleConns(dm.config.Database.MaxIdleConns)

	// SetConnMaxLifetime: コネクションの最大生存時間
	// 長時間の接続による問題（タイムアウト等）を防ぐ
	db.SetConnMaxLifetime(time.Duration(dm.config.Database.ConnMaxLifetime) * time.Minute)

	// 5. 接続テスト（重要：実際にDBに接続を試行）
	if err := dm.pingWithTimeout(db, 10*time.Second); err != nil {
		db.Close() // 接続に失敗した場合はリソースを解放
		return fmt.Errorf("database connection test failed: %w", err)
	}

	dm.DB = db
	log.Printf("Successfully connected to MySQL database")
	return nil
}

// CreateTables はテーブルを作成します
// 標準パッケージを使ったDDL（データ定義言語）の実行を学習
func (dm *DatabaseManager) CreateTables() error {
	// todos テーブル作成用のSQL
	// CREATE TABLE IF NOT EXISTS で既存テーブルがある場合はエラーを回避
	createTodosTable := `
		CREATE TABLE IF NOT EXISTS todos (
			id INT AUTO_INCREMENT PRIMARY KEY,
			title VARCHAR(100) NOT NULL,
			description TEXT,
			is_completed BOOLEAN NOT NULL DEFAULT FALSE,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			
			-- インデックスの作成（検索性能向上）
			INDEX idx_is_completed (is_completed),
			INDEX idx_created_at (created_at)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`

	// DDLの実行
	_, err := dm.DB.Exec(createTodosTable)
	if err != nil {
		return fmt.Errorf("failed to create todos table: %w", err)
	}

	log.Println("Database tables created successfully")
	return nil
}

// Close はデータベース接続を閉じます
// リソース管理の重要な学習ポイント
func (dm *DatabaseManager) Close() error {
	if dm.DB == nil {
		return nil
	}

	// Close() は全ての接続プールのコネクションを閉じる
	if err := dm.DB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	log.Println("Database connection closed")
	return nil
}

// pingWithTimeout はタイムアウト付きでデータベースの接続テストを行います
// コンテキストを使ったタイムアウト制御の学習
func (dm *DatabaseManager) pingWithTimeout(db *sql.DB, timeout time.Duration) error {
	// コンテキストにタイムアウトを設定
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel() // 関数終了時に必ずキャンセルを実行

	// PingContext で実際にデータベースに接続を試行
	// これにより sql.Open() では検出できない接続エラーを確認
	return db.PingContext(ctx)
}

// HealthCheck はデータベースの健全性をチェックします
// アプリケーションの監視で使用するヘルスチェック機能
func (dm *DatabaseManager) HealthCheck() error {
	if dm.DB == nil {
		return fmt.Errorf("database connection is nil")
	}

	// 軽量なクエリでDB接続状態を確認
	// SELECT 1 は最も軽量な動作確認用クエリ
	var result int
	err := dm.DB.QueryRow("SELECT 1").Scan(&result)
	if err != nil {
		return fmt.Errorf("health check query failed: %w", err)
	}

	if result != 1 {
		return fmt.Errorf("health check query returned unexpected result: %d", result)
	}

	return nil
}

// GetStats は接続プールの統計情報を返します
// パフォーマンスチューニングと監視に活用
func (dm *DatabaseManager) GetStats() (map[string]interface{}, error) {
	if dm.DB == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	// sql.DB.Stats() で詳細な接続プール情報を取得
	stats := dm.DB.Stats()

	return map[string]interface{}{
		"max_open_connections": stats.MaxOpenConnections,    // 設定された最大オープン接続数
		"open_connections":     stats.OpenConnections,       // 現在のオープン接続数
		"in_use":               stats.InUse,                 // 現在使用中の接続数
		"idle":                 stats.Idle,                  // 現在アイドル状態の接続数
		"wait_count":           stats.WaitCount,             // 接続待ちが発生した回数
		"wait_duration":        stats.WaitDuration.String(), // 接続待ちの累積時間
		"max_idle_closed":      stats.MaxIdleClosed,         // アイドル上限で閉じられた接続数
		"max_idle_time_closed": stats.MaxIdleTimeClosed,     // アイドル時間で閉じられた接続数
		"max_lifetime_closed":  stats.MaxLifetimeClosed,     // 生存時間で閉じられた接続数
	}, nil
}

// ExecuteMigration はマイグレーションSQLを実行します（将来の拡張用）
// バージョン管理されたスキーマ変更の実装例
func (dm *DatabaseManager) ExecuteMigration(migrationSQL string) error {
	if dm.DB == nil {
		return fmt.Errorf("database connection is nil")
	}

	// トランザクション内でマイグレーションを実行
	// 失敗した場合は自動的にロールバック
	tx, err := dm.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// defer でのトランザクション管理
	// パニックやエラーが発生した場合の安全な処理
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // パニックを再発生させる
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// マイグレーションSQL実行
	_, err = tx.Exec(migrationSQL)
	if err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	return nil
}

// database/sql パッケージ使用時のベストプラクティス：
//
// 1. コネクションプール設定：
//    - MaxOpenConns: アプリケーションの並行性に応じて調整
//    - MaxIdleConns: レスポンス性能とメモリ使用量のバランス
//    - ConnMaxLifetime: DB側のタイムアウト設定より短く設定
//
// 2. エラーハンドリング：
//    - sql.ErrNoRows の適切な処理
//    - fmt.Errorf でのエラーコンテキスト追加
//    - リソースリークを防ぐエラー時のクリーンアップ
//
// 3. リソース管理：
//    - sql.Rows の確実なClose()
//    - データベース接続の適切な終了処理
//    - コンテキストを使ったタイムアウト制御
//
// 4. パフォーマンス：
//    - プリペアードステートメントの活用
//    - 適切なインデックス設計
//    - 接続プール統計による監視
//
// 5. セキュリティ：
//    - SQLインジェクション対策（プリペアードステートメント）
//    - 接続文字列の適切な管理
//    - 最小権限でのDB接続
