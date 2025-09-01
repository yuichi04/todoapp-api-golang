package middleware

import (
	"net/http"
)

// CORSMiddleware は標準パッケージを使用してCORS（Cross-Origin Resource Sharing）を実装するミドルウェアです
//
// 標準パッケージでのミドルウェアパターンの学習ポイント：
// 1. http.HandlerFunc を返すファクトリー関数パターン
// 2. クロージャーを使った設定の保持
// 3. next ハンドラーの呼び出しチェーン
// 4. リクエスト・レスポンスの前処理・後処理
//
// CORSの役割：
// 1. 異なるオリジン（ドメイン、プロトコル、ポート）からのHTTPリクエストを制御
// 2. ブラウザのセキュリティポリシー（Same-Origin Policy）を緩和
// 3. プリフライトリクエスト（OPTIONS）の適切な処理
type CORSConfig struct {
	// AllowedOrigins は許可するオリジンのリスト
	// "*" ですべてのオリジンを許可（開発環境用）
	AllowedOrigins []string

	// AllowedMethods は許可するHTTPメソッドのリスト
	AllowedMethods []string

	// AllowedHeaders は許可するリクエストヘッダーのリスト
	AllowedHeaders []string

	// AllowCredentials は認証情報を含むリクエストを許可するか
	AllowCredentials bool

	// MaxAge はプリフライトリクエストの結果をキャッシュする時間（秒）
	MaxAge int
}

// DefaultCORSConfig は開発環境用のデフォルト設定を返します
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodPatch,
			http.MethodOptions,
		},
		AllowedHeaders: []string{
			"Origin",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"X-CSRF-Token",
			"Authorization",
			"Accept",
			"Cache-Control",
			"X-Requested-With",
		},
		AllowCredentials: false,
		MaxAge:          86400, // 24時間
	}
}

// CORSMiddleware は設定可能なCORSミドルウェアを作成します
// ファクトリー関数パターン：設定を受け取ってミドルウェア関数を生成
func CORSMiddleware(config CORSConfig) func(http.Handler) http.Handler {
	// クロージャーで設定を保持し、実際のミドルウェア関数を返す
	return func(next http.Handler) http.Handler {
		// http.HandlerFunc は func(ResponseWriter, *Request) を http.Handler に変換
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. オリジンの確認と設定
			origin := r.Header.Get("Origin")
			if isOriginAllowed(origin, config.AllowedOrigins) {
				if len(config.AllowedOrigins) == 1 && config.AllowedOrigins[0] == "*" {
					// ワイルドカードの場合
					w.Header().Set("Access-Control-Allow-Origin", "*")
				} else {
					// 具体的なオリジンの場合
					w.Header().Set("Access-Control-Allow-Origin", origin)
				}
			}

			// 2. 基本的なCORSヘッダーを設定
			w.Header().Set("Access-Control-Allow-Methods", joinStrings(config.AllowedMethods, ", "))
			w.Header().Set("Access-Control-Allow-Headers", joinStrings(config.AllowedHeaders, ", "))
			
			// 3. 認証情報の許可設定
			if config.AllowCredentials {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			} else {
				w.Header().Set("Access-Control-Allow-Credentials", "false")
			}

			// 4. プリフライト結果のキャッシュ時間
			if config.MaxAge > 0 {
				w.Header().Set("Access-Control-Max-Age", intToString(config.MaxAge))
			}

			// 5. プリフライトリクエスト（OPTIONS）の処理
			// ブラウザが実際のリクエスト前に送信する事前チェックリクエスト
			if r.Method == http.MethodOptions {
				// プリフライトリクエストには200 OKで即座に応答
				// 実際のハンドラー処理は実行しない
				w.WriteHeader(http.StatusOK)
				return
			}

			// 6. 次のミドルウェアまたはハンドラーに処理を移す
			// これがミドルウェアチェーンの核心概念
			next.ServeHTTP(w, r)
		})
	}
}

// SimpleCORSMiddleware はシンプルなCORSミドルウェアです（学習用）
// より簡素な実装でミドルウェアの基本概念を理解
func SimpleCORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 開発環境用の緩い設定
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// プリフライトリクエストの処理
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// 次のハンドラーを呼び出し
		next.ServeHTTP(w, r)
	})
}

// --- ヘルパー関数 ---

// isOriginAllowed は指定されたオリジンが許可リストに含まれているかチェックします
func isOriginAllowed(origin string, allowedOrigins []string) bool {
	if len(allowedOrigins) == 0 {
		return false
	}

	// ワイルドカード（*）の場合は全て許可
	for _, allowed := range allowedOrigins {
		if allowed == "*" {
			return true
		}
		if allowed == origin {
			return true
		}
	}

	return false
}

// joinStrings は文字列スライスを指定されたセパレータで結合します
// strings.Join の代替として学習用に実装
func joinStrings(strs []string, separator string) string {
	if len(strs) == 0 {
		return ""
	}

	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += separator + strs[i]
	}

	return result
}

// intToString は整数を文字列に変換します
// strconv.Itoa の代替として学習用に実装
func intToString(i int) string {
	if i == 0 {
		return "0"
	}

	// 簡単な整数→文字列変換（正の整数のみ）
	digits := []byte{}
	for i > 0 {
		digits = append([]byte{byte(i%10) + '0'}, digits...)
		i /= 10
	}

	return string(digits)
}

// 標準パッケージでのミドルウェアパターンの学習ポイント：
//
// 1. ミドルウェアの基本構造：
//    - func(http.Handler) http.Handler 型
//    - next.ServeHTTP(w, r) での処理継続
//    - 前処理・後処理の実装
//
// 2. ファクトリーパターン：
//    - 設定を受け取ってミドルウェアを生成
//    - クロージャーでの設定保持
//    - 再利用可能な設計
//
// 3. チェーンパターン：
//    - 複数のミドルウェアを組み合わせ
//    - リクエスト処理の前後での共通処理
//    - 処理の順序制御
//
// 4. HTTPヘッダー操作：
//    - ResponseWriter.Header() の使用
//    - ブラウザ向けの適切なヘッダー設定
//
// 5. エラーハンドリング：
//    - 早期リターンでの処理中断
//    - 適切なHTTPステータスコード設定
//
// ミドルウェアの使用例：
// ```go
// mux := http.NewServeMux()
// mux.Handle("/api/", CORSMiddleware(DefaultCORSConfig())(apiHandler))
// mux.Handle("/public/", SimpleCORSMiddleware(publicHandler))
// ```