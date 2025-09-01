package web

import (
	"net/http"
	"strings"

	"todoapp-api-golang/internal/application/handler"
	"todoapp-api-golang/internal/application/middleware"
)

// Router は標準パッケージを使用したHTTPルーティングを管理する構造体です
//
// 標準パッケージでのルーティングの学習ポイント：
// 1. http.ServeMux の基本的な使用方法
// 2. 手動でのパスマッチングとパラメータ抽出
// 3. HTTPメソッドの手動判定
// 4. ミドルウェアチェーンの構築
// 5. RESTful URLパターンの実装
type Router struct {
	mux         *http.ServeMux
	todoHandler *handler.TodoHandler
}

// NewRouter はRouterのコンストラクタです
func NewRouter(todoHandler *handler.TodoHandler) *Router {
	return &Router{
		mux:         http.NewServeMux(),
		todoHandler: todoHandler,
	}
}

// SetupRoutes はHTTPルーティングを設定します
// 標準パッケージでRESTful APIの設計原則を学習
func (router *Router) SetupRoutes() http.Handler {
	// 1. ヘルスチェックエンドポイント
	// システムの稼働状態を確認するためのシンプルなエンドポイント
	router.mux.HandleFunc("/health", router.healthCheckHandler)

	// 2. API v1のルートハンドラー
	// /api/v1/* へのすべてのリクエストを単一のハンドラーで処理
	// 標準パッケージでは詳細なパスマッチングを手動で実装
	router.mux.HandleFunc("/api/v1/", router.apiV1Handler)

	// 3. ミドルウェアチェーンの構築
	// 複数のミドルウェアを組み合わせてリクエスト処理を強化
	finalHandler := middleware.ChainMiddleware(
		middleware.RecoveryMiddleware,                              // パニック回復
		middleware.LoggingMiddleware,                               // アクセスログ
		middleware.SimpleCORSMiddleware,                           // CORS対応
		middleware.RequestIDMiddleware,                             // リクエストID付与
	)(router.mux)

	return finalHandler
}

// healthCheckHandler はヘルスチェックエンドポイントのハンドラーです
// GET /health への対応
func (router *Router) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// HTTPメソッドの確認
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// シンプルなJSONレスポンス
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	// 手動でJSONを構築（encoding/jsonを使わない学習用）
	response := `{
		"status": "ok",
		"message": "Todo API is running",
		"version": "1.0.0"
	}`
	w.Write([]byte(response))
}

// apiV1Handler は /api/v1/* へのすべてのリクエストを処理するメインハンドラーです
// 標準パッケージでの手動ルーティングの実装例
func (router *Router) apiV1Handler(w http.ResponseWriter, r *http.Request) {
	// URLパスから /api/v1/ プレフィックスを除去
	path := strings.TrimPrefix(r.URL.Path, "/api/v1")
	path = strings.Trim(path, "/")

	// パスを "/" で分割してセグメント化
	segments := strings.Split(path, "/")
	
	// 空のパスや無効なパスの処理
	if len(segments) == 0 || segments[0] == "" {
		http.NotFound(w, r)
		return
	}

	// リソースタイプによる分岐
	switch segments[0] {
	case "todos":
		router.handleTodosRoutes(w, r, segments[1:])
	default:
		http.NotFound(w, r)
	}
}

// handleTodosRoutes はTodoリソースへのルーティングを処理します
// RESTful APIパターンの手動実装
//
// 対応するエンドポイント：
// GET    /api/v1/todos           -> 一覧取得
// POST   /api/v1/todos           -> 新規作成
// GET    /api/v1/todos/{id}      -> 詳細取得
// PUT    /api/v1/todos/{id}      -> 更新
// DELETE /api/v1/todos/{id}      -> 削除
// PATCH  /api/v1/todos/{id}/complete   -> 完了
// PATCH  /api/v1/todos/{id}/incomplete -> 未完了
func (router *Router) handleTodosRoutes(w http.ResponseWriter, r *http.Request, segments []string) {
	switch len(segments) {
	case 0:
		// /api/v1/todos
		router.handleTodoCollection(w, r)
	case 1:
		// /api/v1/todos/{id}
		router.handleTodoItem(w, r, segments[0])
	case 2:
		// /api/v1/todos/{id}/{action}
		router.handleTodoAction(w, r, segments[0], segments[1])
	default:
		http.NotFound(w, r)
	}
}

// handleTodoCollection はTodoコレクションへの操作を処理します
// /api/v1/todos へのリクエスト
func (router *Router) handleTodoCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// GET /api/v1/todos -> 全Todo取得
		router.todoHandler.GetAllTodos(w, r)
	case http.MethodPost:
		// POST /api/v1/todos -> 新Todo作成
		router.todoHandler.CreateTodo(w, r)
	default:
		// サポートされていないHTTPメソッド
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleTodoItem は個別Todoアイテムへの操作を処理します
// /api/v1/todos/{id} へのリクエスト
func (router *Router) handleTodoItem(w http.ResponseWriter, r *http.Request, id string) {
	// IDの基本的な検証（空文字チェック）
	if id == "" {
		http.Error(w, "Todo ID is required", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// GET /api/v1/todos/{id} -> Todo詳細取得
		router.todoHandler.GetTodoByID(w, r)
	case http.MethodPut:
		// PUT /api/v1/todos/{id} -> Todo更新
		router.todoHandler.UpdateTodo(w, r)
	case http.MethodDelete:
		// DELETE /api/v1/todos/{id} -> Todo削除
		router.todoHandler.DeleteTodo(w, r)
	default:
		// サポートされていないHTTPメソッド
		w.Header().Set("Allow", "GET, PUT, DELETE")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleTodoAction は特定のTodoに対するアクションを処理します
// /api/v1/todos/{id}/{action} へのリクエスト
func (router *Router) handleTodoAction(w http.ResponseWriter, r *http.Request, id, action string) {
	// IDの基本的な検証
	if id == "" {
		http.Error(w, "Todo ID is required", http.StatusBadRequest)
		return
	}

	// PATCHメソッドのみサポート
	if r.Method != http.MethodPatch {
		w.Header().Set("Allow", "PATCH")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// アクションタイプによる分岐
	switch action {
	case "complete":
		// PATCH /api/v1/todos/{id}/complete -> Todo完了
		router.todoHandler.CompleteTodo(w, r)
	case "incomplete":
		// PATCH /api/v1/todos/{id}/incomplete -> Todo未完了
		router.todoHandler.IncompleteTodo(w, r)
	default:
		http.NotFound(w, r)
	}
}

// GetMux はhttp.ServeMuxを返します（テスト等で使用）
func (router *Router) GetMux() *http.ServeMux {
	return router.mux
}

// 標準パッケージでのルーティング学習のポイント：
//
// 1. ServeMux の基本：
//    - http.NewServeMux() での作成
//    - HandleFunc() でのハンドラー登録
//    - パターンマッチングの制限と回避方法
//
// 2. 手動ルーティング：
//    - URL パスの解析と分割
//    - strings パッケージの活用
//    - セグメントベースのルーティング
//
// 3. RESTful 設計：
//    - リソース指向のURL構造
//    - HTTPメソッドによる操作の分離
//    - エラーレスポンスの統一
//
// 4. ミドルウェアパターン：
//    - func(http.Handler) http.Handler 型の活用
//    - チェーン構築による機能組み合わせ
//    - 横断的関心事の分離
//
// 5. エラーハンドリング：
//    - 適切なHTTPステータスコードの設定
//    - Allow ヘッダーでのメソッド通知
//    - 一貫性のあるエラーレスポンス
//
// 標準パッケージでの制限と対策：
// - パスパラメータの自動抽出がない → 手動パース
// - HTTPメソッドの自動判定がない → 手動チェック
// - ミドルウェアの標準実装がない → 自作ミドルウェア
// - 複雑なルーティングルールがない → 単純化または手動実装
//
// これらの制限により、Goのnet/httpパッケージの基本概念を
// より深く理解することができます。