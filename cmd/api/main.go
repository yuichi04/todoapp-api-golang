package main

import (
	"log"

	"todoapp-api-golang/internal/application/handler"
	"todoapp-api-golang/internal/domain/service"
	"todoapp-api-golang/internal/infrastructure/database"
	"todoapp-api-golang/internal/infrastructure/web"
	"todoapp-api-golang/pkg/config"
)

// main はアプリケーションのエントリーポイント（開始点）です
// 標準パッケージを使用したアプリケーション構築の学習ポイント：
// 1. 依存関係の手動管理と依存性注入
// 2. Clean Architectureの層構造の実装
// 3. データベース接続とテーブル作成の管理
// 4. エラーハンドリングとログ出力
// 5. アプリケーションライフサイクルの管理
func main() {
	// アプリケーション初期化の開始ログ
	log.Println("Starting Todo API application with standard packages...")

	// 1. 設定の読み込み
	// 環境変数から設定値を読み込み、デフォルト値で補完
	cfg, err := config.Load()
	if err != nil {
		// 設定読み込みに失敗した場合はアプリケーションを停止
		// log.Fatal()は log.Print() の後に os.Exit(1) を実行
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 設定内容のログ出力（本番環境では機密情報を除外すること）
	log.Printf("Configuration loaded - Environment: %s, Port: %d, DB Driver: %s", 
		cfg.App.Environment, cfg.Server.Port, cfg.Database.Driver)

	// 2. データベース接続の確立
	// 標準パッケージを使用したデータベースマネージャーの作成と接続
	dbManager := database.NewDatabaseManager(cfg)
	if err := dbManager.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// アプリケーション終了時のクリーンアップ処理
	// defer文により、main関数終了時に自動実行される
	defer func() {
		if err := dbManager.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	// 3. データベーステーブルの作成
	// 開発環境では自動テーブル作成、本番環境では手動マイグレーション推奨
	if !cfg.IsProduction() {
		if err := dbManager.CreateTables(); err != nil {
			log.Fatalf("Failed to create database tables: %v", err)
		}
	} else {
		log.Println("Production mode: skipping automatic table creation")
		log.Println("Please ensure database schema is properly migrated")
	}

	// 4. 依存性注入による各層の構築
	// Clean Architectureの依存関係の流れ：
	// main -> Handler -> Service -> Repository -> Database
	
	// 4-1. リポジトリ層（データアクセス）の初期化
	// 標準のdatabase/sqlパッケージを使用したリポジトリ実装
	todoRepo := database.NewTodoRepository(dbManager.DB)
	
	// 4-2. ドメインサービス層（ビジネスロジック）の初期化
	// リポジトリをサービスに注入
	todoService := service.NewTodoService(todoRepo)
	
	// 4-3. ハンドラー層（HTTP処理）の初期化  
	// サービスをハンドラーに注入
	todoHandler := handler.NewTodoHandler(todoService)

	// 4-4. ルーティング層の初期化
	// 標準パッケージを使用したルーター作成
	router := web.NewRouter(todoHandler)

	// 4-5. HTTPサーバー層の初期化
	server := web.NewServer(cfg, router)

	// 5. データベース接続の健全性チェック
	// アプリケーション起動前の最終確認
	if err := dbManager.HealthCheck(); err != nil {
		log.Fatalf("Database health check failed: %v", err)
	}

	// 6. 接続プール統計情報の出力（デバッグ用）
	if !cfg.IsProduction() {
		if stats, err := dbManager.GetStats(); err == nil {
			log.Printf("Database connection pool stats: %+v", stats)
		}
	}

	// 7. アプリケーション起動の完了ログ
	log.Printf("Todo API is ready to serve requests")
	log.Printf("Server will start on: http://%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Health check endpoint: http://%s:%d/health", cfg.Server.Host, cfg.Server.Port)
	log.Printf("API base URL: http://%s:%d/api/v1", cfg.Server.Host, cfg.Server.Port)

	// 8. HTTPサーバーの起動
	// Start()は内部でグレースフルシャットダウンを処理
	// ブロッキング関数のため、ここでアプリケーションが待機状態になる
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// 標準パッケージを使用したアプリケーション構築の学習ポイント：
//
// 1. 手動依存性注入：
//    - 各層のインスタンス生成を明示的に管理
//    - 依存関係の方向性を意識した構築順序
//    - インターフェースを通じた疎結合の実現
//
// 2. エラーハンドリング：
//    - 各段階でのエラーチェックと適切な対応
//    - log.Fatalf() による致命的エラーの処理
//    - defer による確実なリソース解放
//
// 3. 設定管理：
//    - 環境変数ベースの設定読み込み
//    - 環境別の処理分岐
//    - デフォルト値による堅牢性確保
//
// 4. データベース管理：
//    - 接続確立、テーブル作成、健全性チェック
//    - 接続プールの統計情報活用
//    - プロダクション環境での注意事項
//
// 5. アプリケーションライフサイクル：
//    - 初期化順序の重要性
//    - グレースフルシャットダウンの準備
//    - 適切なログ出力による可視性確保
//
// フレームワークを使わない利点：
// - 各コンポーネントの役割と責任の明確化
// - 依存関係の完全な理解
// - カスタマイズの自由度
// - デバッグ時の透明性
// - パフォーマンスの最適化可能性
//
// 学習効果：
// - Goの標準パッケージの深い理解
// - Clean Architectureの実装体験
// - HTTPサーバーの低レベル実装知識
// - データベース操作の基本原理
// - エラーハンドリングとログ設計