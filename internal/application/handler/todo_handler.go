package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"todoapp-api-golang/internal/application/dto"
	"todoapp-api-golang/internal/domain/service"
)

// TodoHandler はTodo関連のHTTPリクエストを処理するハンドラーです
// 標準のnet/httpパッケージを使用してHTTP処理を実装します
//
// net/httpパッケージの学習ポイント：
// 1. http.HandlerFunc の理解
// 2. http.ResponseWriter と *http.Request の使い方
// 3. JSONの手動パース（encoding/json）
// 4. HTTPステータスコードの設定
// 5. URLパスパラメータの取得
type TodoHandler struct {
	// todoService はビジネスロジック処理を担当するドメインサービス
	// 依存性注入によってサービス実装を受け取ります
	todoService *service.TodoService
}

// NewTodoHandler はTodoHandlerのコンストラクタです
// 標準パッケージを使った依存性注入の実装例
func NewTodoHandler(todoService *service.TodoService) *TodoHandler {
	return &TodoHandler{
		todoService: todoService,
	}
}

// CreateTodo は新しいTodoを作成するHTTPハンドラーです
// POST /api/v1/todos へのリクエストを処理します
//
// 標準パッケージでのHTTP処理の学習ポイント：
// 1. http.ResponseWriter での レスポンス書き込み
// 2. json.Decoder での リクエストボディの解析
// 3. Content-Type ヘッダーの設定
// 4. エラーハンドリング パターン
func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	// 1. HTTPメソッドの確認
	// 標準パッケージでは手動でメソッドチェックが必要
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 2. Content-Typeの確認（JSON以外を拒否）
	contentType := r.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		http.Error(w, "Content-Type must be application/json", http.StatusBadRequest)
		return
	}

	// 3. JSONリクエストボディをDTOにデコード
	var req dto.CreateTodoRequest
	
	// json.NewDecoder を使ってストリームからJSONを読み取り
	// これはGinの ShouldBindJSON() と同等の処理を手動で行う
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		// JSONパースエラーの場合は400 Bad Requestを返す
		writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON format", err.Error())
		return
	}

	// 4. 基本的なバリデーション（手動実装）
	if req.Title == "" {
		writeErrorResponse(w, http.StatusBadRequest, "Validation failed", "title is required")
		return
	}
	if len(req.Title) > 100 {
		writeErrorResponse(w, http.StatusBadRequest, "Validation failed", "title must be 100 characters or less")
		return
	}
	if len(req.Description) > 500 {
		writeErrorResponse(w, http.StatusBadRequest, "Validation failed", "description must be 500 characters or less")
		return
	}

	// 5. DTOからエンティティへの変換
	todo := req.ToEntity()

	// 6. ドメインサービスを呼び出してビジネスロジック実行
	createdTodo, err := h.todoService.CreateTodo(r.Context(), todo)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to create todo", err.Error())
		return
	}

	// 7. エンティティからレスポンスDTOへの変換
	response := dto.ToTodoResponse(createdTodo)

	// 8. JSON レスポンスの書き込み
	writeJSONResponse(w, http.StatusCreated, response)
}

// GetTodoByID は指定されたIDのTodoを取得するHTTPハンドラーです
// GET /api/v1/todos/{id} へのリクエストを処理します
//
// URLパスパラメータの取得方法を学習：
// 標準パッケージでは r.URL.Path から手動でパラメータを抽出
func (h *TodoHandler) GetTodoByID(w http.ResponseWriter, r *http.Request) {
	// 1. HTTPメソッドの確認
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 2. URLパスからIDを抽出
	// パスの構造: /api/v1/todos/{id}
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 4 || pathParts[3] == "" {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid URL", "todo ID is required")
		return
	}

	// 3. 文字列を整数に変換
	id, err := strconv.Atoi(pathParts[3])
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid todo ID", "ID must be a number")
		return
	}

	// 4. ドメインサービスでTodo取得
	todo, err := h.todoService.GetTodoByID(r.Context(), id)
	if err != nil {
		// エラーメッセージの内容に応じてHTTPステータスを決定
		if strings.Contains(err.Error(), "not found") {
			writeErrorResponse(w, http.StatusNotFound, "Todo not found", "")
		} else {
			writeErrorResponse(w, http.StatusInternalServerError, "Failed to get todo", err.Error())
		}
		return
	}

	// 5. レスポンス返却
	response := dto.ToTodoResponse(todo)
	writeJSONResponse(w, http.StatusOK, response)
}

// GetAllTodos は全てのTodoを取得するHTTPハンドラーです
// GET /api/v1/todos へのリクエストを処理します
//
// クエリパラメータの処理方法を学習：
// r.URL.Query() を使ってクエリパラメータを取得
func (h *TodoHandler) GetAllTodos(w http.ResponseWriter, r *http.Request) {
	// 1. HTTPメソッドの確認
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 2. クエリパラメータの解析
	query := r.URL.Query()
	
	// ページング用パラメータの取得（将来拡張用）
	page := 1
	if p := query.Get("page"); p != "" {
		if pageNum, err := strconv.Atoi(p); err == nil && pageNum > 0 {
			page = pageNum
		}
	}
	
	limit := 10
	if l := query.Get("limit"); l != "" {
		if limitNum, err := strconv.Atoi(l); err == nil && limitNum > 0 && limitNum <= 100 {
			limit = limitNum
		}
	}

	// 3. ドメインサービスで全Todo取得
	todos, err := h.todoService.GetAllTodos(r.Context())
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to get todos", err.Error())
		return
	}

	// 4. レスポンス生成
	response := dto.ToTodoListResponse(todos, page, limit, len(todos))
	writeJSONResponse(w, http.StatusOK, response)
}

// UpdateTodo は既存のTodoを更新するHTTPハンドラーです
// PUT /api/v1/todos/{id} へのリクエストを処理します
func (h *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	// 1. HTTPメソッドの確認
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 2. Content-Typeの確認
	contentType := r.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		http.Error(w, "Content-Type must be application/json", http.StatusBadRequest)
		return
	}

	// 3. URLパスからIDを抽出
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 4 {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid URL", "todo ID is required")
		return
	}

	id, err := strconv.Atoi(pathParts[3])
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid todo ID", "ID must be a number")
		return
	}

	// 4. リクエストボディの解析
	var req dto.UpdateTodoRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON format", err.Error())
		return
	}

	// 5. 更新対象のTodoを取得
	todo, err := h.todoService.GetTodoByID(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeErrorResponse(w, http.StatusNotFound, "Todo not found", "")
		} else {
			writeErrorResponse(w, http.StatusInternalServerError, "Failed to get todo", err.Error())
		}
		return
	}

	// 6. リクエストの内容を既存Todoに適用（部分更新）
	req.ApplyToEntity(todo)

	// 7. ドメインサービスで更新実行
	updatedTodo, err := h.todoService.UpdateTodo(r.Context(), todo)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to update todo", err.Error())
		return
	}

	// 8. レスポンス返却
	response := dto.ToTodoResponse(updatedTodo)
	writeJSONResponse(w, http.StatusOK, response)
}

// DeleteTodo は指定されたIDのTodoを削除するHTTPハンドラーです
// DELETE /api/v1/todos/{id} へのリクエストを処理します
func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	// 1. HTTPメソッドの確認
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 2. URLパスからIDを抽出
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 4 {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid URL", "todo ID is required")
		return
	}

	id, err := strconv.Atoi(pathParts[3])
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid todo ID", "ID must be a number")
		return
	}

	// 3. ドメインサービスで削除実行
	err = h.todoService.DeleteTodo(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeErrorResponse(w, http.StatusNotFound, "Todo not found", "")
		} else {
			writeErrorResponse(w, http.StatusInternalServerError, "Failed to delete todo", err.Error())
		}
		return
	}

	// 4. 削除成功時は204 No Contentを返却（レスポンスボディなし）
	w.WriteHeader(http.StatusNoContent)
}

// CompleteTodo はTodoを完了状態にするHTTPハンドラーです  
// PATCH /api/v1/todos/{id}/complete へのリクエストを処理します
func (h *TodoHandler) CompleteTodo(w http.ResponseWriter, r *http.Request) {
	// 1. HTTPメソッドの確認
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 2. URLパスからIDを抽出
	// パスの構造: /api/v1/todos/{id}/complete
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 5 || pathParts[4] != "complete" {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid URL", "invalid endpoint")
		return
	}

	id, err := strconv.Atoi(pathParts[3])
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid todo ID", "ID must be a number")
		return
	}

	// 3. ドメインサービスでTodo完了処理
	completedTodo, err := h.todoService.CompleteTodo(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeErrorResponse(w, http.StatusNotFound, "Todo not found", "")
		} else {
			writeErrorResponse(w, http.StatusInternalServerError, "Failed to complete todo", err.Error())
		}
		return
	}

	// 4. レスポンス返却
	response := dto.ToTodoResponse(completedTodo)
	writeJSONResponse(w, http.StatusOK, response)
}

// IncompleteTodo はTodoを未完了状態に戻すHTTPハンドラーです
// PATCH /api/v1/todos/{id}/incomplete へのリクエストを処理します
func (h *TodoHandler) IncompleteTodo(w http.ResponseWriter, r *http.Request) {
	// 1. HTTPメソッドの確認
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 2. URLパスからIDを抽出
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 5 || pathParts[4] != "incomplete" {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid URL", "invalid endpoint")
		return
	}

	id, err := strconv.Atoi(pathParts[3])
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid todo ID", "ID must be a number")
		return
	}

	// 3. ドメインサービスでTodo未完了処理
	incompleteTodo, err := h.todoService.IncompleteTodo(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeErrorResponse(w, http.StatusNotFound, "Todo not found", "")
		} else {
			writeErrorResponse(w, http.StatusInternalServerError, "Failed to mark todo as incomplete", err.Error())
		}
		return
	}

	// 4. レスポンス返却
	response := dto.ToTodoResponse(incompleteTodo)
	writeJSONResponse(w, http.StatusOK, response)
}

// --- ヘルパー関数 ---

// writeJSONResponse はJSONレスポンスを書き込むヘルパー関数です
// 標準パッケージでのJSON出力の学習に重要
func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	// 1. Content-Typeヘッダーを設定
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	
	// 2. ステータスコードを設定
	w.WriteHeader(statusCode)
	
	// 3. JSONエンコードしてレスポンス書き込み
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(data); err != nil {
		// JSON encoding に失敗した場合のフォールバック
		// ただし、この時点では既にステータスコードが送信されているため変更不可
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// writeErrorResponse はエラーレスポンスを書き込むヘルパー関数です
func writeErrorResponse(w http.ResponseWriter, statusCode int, message, details string) {
	errorResponse := dto.ErrorResponse{
		Error:   message,
		Details: details,
	}
	writeJSONResponse(w, statusCode, errorResponse)
}

// 標準パッケージを使ったHTTP処理の学習ポイント：
//
// 1. 低レベルAPI の理解：
//    - http.ResponseWriter と *http.Request の詳細な操作
//    - ヘッダー、ステータスコード、ボディの明示的な制御
//
// 2. 手動でのリクエスト処理：
//    - メソッド判定、パスパラメータ抽出、クエリパラメータ解析
//    - Content-Type チェック、JSONデコード/エンコード
//
// 3. エラーハンドリング：
//    - 各段階でのエラー処理と適切なHTTPステータス返却
//    - ユーザーフレンドリーなエラーメッセージ
//
// 4. RESTful API の手動実装：
//    - 各HTTPメソッドに対応したハンドラー実装
//    - 統一されたレスポンス形式
//
// 5. 保守性とテスタビリティ：
//    - ヘルパー関数による共通処理の抽出
//    - 依存性注入によるテスト容易性