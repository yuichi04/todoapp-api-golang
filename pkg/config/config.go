package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config はアプリケーション全体の設定を管理する構造体です
// 設定管理の役割：
// 1. 環境変数から設定値を読み込み
// 2. デフォルト値の設定
// 3. 設定値の型変換とバリデーション
// 4. 環境別設定（開発、テスト、本番）の統一管理
type Config struct {
	// Server はHTTPサーバー関連の設定
	Server ServerConfig `json:"server"`
	
	// Database はデータベース接続関連の設定
	Database DatabaseConfig `json:"database"`
	
	// App はアプリケーション固有の設定
	App AppConfig `json:"app"`
}

// ServerConfig はHTTPサーバーの設定を管理します
type ServerConfig struct {
	// Port はHTTPサーバーが使用するポート番号
	Port int `json:"port"`
	
	// Host はHTTPサーバーがバインドするホスト名/IPアドレス
	Host string `json:"host"`
	
	// ReadTimeout は読み取りタイムアウト（秒）
	ReadTimeout int `json:"read_timeout"`
	
	// WriteTimeout は書き込みタイムアウト（秒）
	WriteTimeout int `json:"write_timeout"`
}

// DatabaseConfig はデータベース接続の設定を管理します
type DatabaseConfig struct {
	// Driver はデータベースドライバー名（mysql, postgres等）
	Driver string `json:"driver"`
	
	// Host はデータベースサーバーのホスト名
	Host string `json:"host"`
	
	// Port はデータベースサーバーのポート番号
	Port int `json:"port"`
	
	// Name はデータベース名
	Name string `json:"name"`
	
	// User はデータベース接続ユーザー名
	User string `json:"user"`
	
	// Password はデータベース接続パスワード
	Password string `json:"password"`
	
	// SSLMode はSSL接続モード（postgres用）
	SSLMode string `json:"ssl_mode"`
	
	// MaxOpenConns は最大オープン接続数
	MaxOpenConns int `json:"max_open_conns"`
	
	// MaxIdleConns は最大アイドル接続数
	MaxIdleConns int `json:"max_idle_conns"`
	
	// ConnMaxLifetime は接続の最大生存時間（分）
	ConnMaxLifetime int `json:"conn_max_lifetime"`
}

// AppConfig はアプリケーション固有の設定を管理します
type AppConfig struct {
	// Environment は実行環境（development, production, test）
	Environment string `json:"environment"`
	
	// LogLevel はログレベル（debug, info, warn, error）
	LogLevel string `json:"log_level"`
	
	// Version はアプリケーションバージョン
	Version string `json:"version"`
}

// Load は環境変数から設定を読み込んでConfig構造体を作成します
// 12-Factor Appの原則に従い、設定は環境変数から読み込みます
func Load() (*Config, error) {
	config := &Config{
		// サーバー設定の読み込み
		Server: ServerConfig{
			Port:         getEnvAsInt("SERVER_PORT", 8080),           // デフォルト: 8080
			Host:         getEnv("SERVER_HOST", "0.0.0.0"),          // デフォルト: 全IPでバインド
			ReadTimeout:  getEnvAsInt("SERVER_READ_TIMEOUT", 30),    // デフォルト: 30秒
			WriteTimeout: getEnvAsInt("SERVER_WRITE_TIMEOUT", 30),   // デフォルト: 30秒
		},
		
		// データベース設定の読み込み
		Database: DatabaseConfig{
			Driver:          getEnv("DB_DRIVER", "mysql"),                    // デフォルト: MySQL
			Host:            getEnv("DB_HOST", "localhost"),                  // デフォルト: localhost
			Port:            getEnvAsInt("DB_PORT", 3306),                   // デフォルト: MySQL標準ポート
			Name:            getEnv("DB_NAME", "todoapp"),                   // デフォルト: todoapp
			User:            getEnv("DB_USER", "root"),                      // デフォルト: root
			Password:        getEnv("DB_PASSWORD", ""),                      // デフォルト: パスワードなし
			SSLMode:         getEnv("DB_SSL_MODE", "disable"),               // デフォルト: SSL無効
			MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 10),          // デフォルト: 10接続
			MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 5),           // デフォルト: 5接続
			ConnMaxLifetime: getEnvAsInt("DB_CONN_MAX_LIFETIME", 60),       // デフォルト: 60分
		},
		
		// アプリケーション設定の読み込み
		App: AppConfig{
			Environment: getEnv("APP_ENV", "development"),      // デフォルト: 開発環境
			LogLevel:    getEnv("LOG_LEVEL", "info"),          // デフォルト: infoレベル
			Version:     getEnv("APP_VERSION", "1.0.0"),       // デフォルト: 1.0.0
		},
	}
	
	// 設定値のバリデーション
	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("config validation error: %w", err)
	}
	
	return config, nil
}

// validate は設定値の妥当性をチェックします
func (c *Config) validate() error {
	// サーバーポートの範囲チェック
	if c.Server.Port < 1 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d (must be 1-65535)", c.Server.Port)
	}
	
	// データベース名の必須チェック
	if c.Database.Name == "" {
		return fmt.Errorf("database name is required")
	}
	
	// 環境の値チェック
	if c.App.Environment != "development" && 
	   c.App.Environment != "production" && 
	   c.App.Environment != "test" {
		return fmt.Errorf("invalid environment: %s (must be development, production, or test)", c.App.Environment)
	}
	
	// ログレベルの値チェック
	if c.App.LogLevel != "debug" && 
	   c.App.LogLevel != "info" && 
	   c.App.LogLevel != "warn" && 
	   c.App.LogLevel != "error" {
		return fmt.Errorf("invalid log level: %s (must be debug, info, warn, or error)", c.App.LogLevel)
	}
	
	return nil
}

// GetDSN はデータベース接続文字列（DSN: Data Source Name）を生成します
// データベースドライバーに応じて適切な接続文字列を返します
func (c *Config) GetDSN() string {
	switch c.Database.Driver {
	case "mysql":
		// MySQL用DSN形式: user:password@tcp(host:port)/dbname?parseTime=true
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4",
			c.Database.User,
			c.Database.Password,
			c.Database.Host,
			c.Database.Port,
			c.Database.Name,
		)
	case "postgres":
		// PostgreSQL用DSN形式: host=localhost port=5432 user=gorm dbname=gorm password=gorm sslmode=disable
		return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
			c.Database.Host,
			c.Database.Port,
			c.Database.User,
			c.Database.Name,
			c.Database.Password,
			c.Database.SSLMode,
		)
	case "sqlite":
		// SQLite用DSN（開発・テスト環境用）
		return c.Database.Name + ".db"
	default:
		// デフォルトはMySQL形式
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4",
			c.Database.User,
			c.Database.Password,
			c.Database.Host,
			c.Database.Port,
			c.Database.Name,
		)
	}
}

// IsProduction は本番環境かどうかを判定します
func (c *Config) IsProduction() bool {
	return c.App.Environment == "production"
}

// IsDevelopment は開発環境かどうかを判定します
func (c *Config) IsDevelopment() bool {
	return c.App.Environment == "development"
}

// IsTest はテスト環境かどうかを判定します
func (c *Config) IsTest() bool {
	return c.App.Environment == "test"
}

// --- ヘルパー関数 ---

// getEnv は環境変数を取得し、存在しない場合はデフォルト値を返します
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt は環境変数を整数として取得し、存在しない場合や変換に失敗した場合はデフォルト値を返します
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvAsBool は環境変数をbool値として取得します（将来の拡張用）
func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// 設定管理のベストプラクティス：
//
// 1. 環境変数の活用: 12-Factor Appの原則に従った設定管理
// 2. デフォルト値: 設定漏れによる起動失敗を防ぐ
// 3. バリデーション: 不正な設定値による実行時エラーを防ぐ
// 4. 型安全性: 文字列以外の型（int, bool等）の適切な変換
// 5. セキュリティ: 機密情報（パスワード等）のログ出力回避
// 6. 文書化: 各設定項目の説明とデフォルト値の明記
// 7. 環境別設定: 開発、テスト、本番環境の適切な分離