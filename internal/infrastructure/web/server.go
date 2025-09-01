package web

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"todoapp-api-golang/pkg/config"
)

// Server は標準パッケージを使用してHTTPサーバーを管理する構造体です
//
// 標準パッケージでのHTTPサーバー管理の学習ポイント：
// 1. http.Server の詳細設定
// 2. グレースフルシャットダウンの実装
// 3. シグナルハンドリング（os/signal）
// 4. context パッケージによるタイムアウト制御
// 5. サーバー設定のベストプラクティス
type Server struct {
	httpServer *http.Server
	config     *config.Config
	router     *Router
}

// NewServer はServerのコンストラクタです
func NewServer(cfg *config.Config, router *Router) *Server {
	return &Server{
		config: cfg,
		router: router,
	}
}

// Start はHTTPサーバーを起動します
// 標準パッケージでの本格的なサーバー実装を学習
func (s *Server) Start() error {
	// 1. HTTP サーバーの詳細設定
	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port),
		Handler: s.router.SetupRoutes(), // ルーティング設定を取得
		
		// タイムアウト設定（セキュリティとパフォーマンス対策）
		ReadTimeout:  time.Duration(s.config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(s.config.Server.WriteTimeout) * time.Second,
		IdleTimeout:  60 * time.Second, // Keep-Alive接続のタイムアウト
		
		// ヘッダーサイズ制限（DoS攻撃対策）
		MaxHeaderBytes: 1 << 20, // 1MB
		
		// エラーログの設定
		ErrorLog: log.New(os.Stderr, "SERVER ERROR: ", log.LstdFlags|log.Lshortfile),
	}

	// 2. グレースフルシャットダウンの準備
	// 別のgoroutineでシグナル監視を開始
	go s.gracefulShutdown()

	// 3. サーバー起動ログ
	log.Printf("Starting HTTP server on %s (environment: %s)", 
		s.httpServer.Addr, s.config.App.Environment)

	// 4. HTTPSまたはHTTPでの起動
	// 本番環境ではHTTPS、開発環境ではHTTPを使用
	var err error
	if s.shouldUseHTTPS() {
		// HTTPS での起動（証明書が必要）
		certFile := s.getCertFile()
		keyFile := s.getKeyFile()
		log.Printf("Starting HTTPS server with cert: %s", certFile)
		err = s.httpServer.ListenAndServeTLS(certFile, keyFile)
	} else {
		// HTTP での起動
		log.Println("Starting HTTP server (development mode)")
		err = s.httpServer.ListenAndServe()
	}

	// 5. サーバー終了処理
	// http.ErrServerClosed は正常なシャットダウン時に発生する
	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server failed to start: %w", err)
	}

	log.Println("Server stopped")
	return nil
}

// Stop はHTTPサーバーを停止します
// 標準パッケージでのグレースフルシャットダウンの実装
func (s *Server) Stop(ctx context.Context) error {
	if s.httpServer == nil {
		return nil
	}

	log.Println("Shutting down HTTP server...")

	// Shutdown() は新規接続を拒否し、既存接続の完了を待つ
	// contextのタイムアウトで強制終了のタイミングを制御
	return s.httpServer.Shutdown(ctx)
}

// gracefulShutdown はシステムシグナルを監視してグレースフルシャットダウンを実行します
// 標準パッケージでのシグナルハンドリングを学習
func (s *Server) gracefulShutdown() {
	// 1. シグナルを受信するチャンネルを作成
	sigChan := make(chan os.Signal, 1)
	
	// 2. 監視するシグナルを登録
	// SIGINT: 割り込みシグナル（Ctrl+C）
	// SIGTERM: 終了シグナル（docker stop、killコマンド等）
	signal.Notify(sigChan, 
		syscall.SIGINT,  // 2
		syscall.SIGTERM, // 15
	)
	
	// 3. シグナル受信を待機（ブロッキング）
	sig := <-sigChan
	log.Printf("Received signal: %v", sig)

	// 4. シャットダウンのタイムアウト設定
	// 30秒以内に既存のリクエスト処理を完了させる
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 5. グレースフルシャットダウンの実行
	if err := s.Stop(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
		os.Exit(1)
	}

	log.Println("Server shutdown completed")
	os.Exit(0)
}

// shouldUseHTTPS はHTTPSを使用すべきかを判定します
func (s *Server) shouldUseHTTPS() bool {
	// 本番環境かつ証明書ファイルが存在する場合のみHTTPS
	return s.config.IsProduction() && s.hasCertificateFiles()
}

// hasCertificateFiles は証明書ファイルが存在するかチェックします
func (s *Server) hasCertificateFiles() bool {
	certFile := s.getCertFile()
	keyFile := s.getKeyFile()
	
	// 両方のファイルが存在することを確認
	if _, err := os.Stat(certFile); os.IsNotExist(err) {
		return false
	}
	if _, err := os.Stat(keyFile); os.IsNotExist(err) {
		return false
	}
	
	return true
}

// getCertFile は証明書ファイルのパスを返します
func (s *Server) getCertFile() string {
	// 環境変数から取得、なければデフォルトパス
	if cert := os.Getenv("TLS_CERT_FILE"); cert != "" {
		return cert
	}
	return "./certs/server.crt"
}

// getKeyFile は秘密鍵ファイルのパスを返します
func (s *Server) getKeyFile() string {
	// 環境変数から取得、なければデフォルトパス
	if key := os.Getenv("TLS_KEY_FILE"); key != "" {
		return key
	}
	return "./certs/server.key"
}

// GetAddr はサーバーのアドレスを返します（テスト用）
func (s *Server) GetAddr() string {
	if s.httpServer != nil {
		return s.httpServer.Addr
	}
	return fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)
}

// GetHandler はサーバーのハンドラーを返します（テスト用）
func (s *Server) GetHandler() http.Handler {
	if s.httpServer != nil {
		return s.httpServer.Handler
	}
	return s.router.SetupRoutes()
}

// IsRunning はサーバーが動作中かどうかを返します
func (s *Server) IsRunning() bool {
	return s.httpServer != nil
}

// 標準パッケージでのHTTPサーバー実装の学習ポイント：
//
// 1. http.Server の詳細設定：
//    - Addr: サーバーアドレス（host:port）
//    - Handler: ルートハンドラー
//    - タイムアウト設定（Read/Write/Idle）
//    - MaxHeaderBytes: セキュリティ対策
//
// 2. グレースフルシャットダウン：
//    - signal.Notify() でのシグナルキャッチ
//    - context.WithTimeout() でのタイムアウト制御
//    - Shutdown() での既存接続完了待ち
//
// 3. HTTPS サポート：
//    - 証明書ファイルの管理
//    - 環境別の設定（HTTP/HTTPS）
//    - セキュリティベストプラクティス
//
// 4. エラーハンドリング：
//    - http.ErrServerClosed の適切な処理
//    - ログ出力によるデバッグ支援
//    - エラー時の適切な終了処理
//
// 5. 設定管理：
//    - 環境変数による動的設定
//    - デフォルト値の提供
//    - 環境別設定の切り替え
//
// 実運用での考慮事項：
// - リバースプロキシ（nginx）での運用
// - ヘルスチェックエンドポイントの提供
// - メトリクス収集（Prometheus等）
// - ログ集約システム（fluentd、Logstash等）
// - サーキットブレーカーパターンの実装