# Go標準パッケージで学ぶバックエンドAPI開発 - 写経チュートリアル

## 🎯 このチュートリアルの目標

このチュートリアルでは、Go言語の標準パッケージのみを使用してTodo APIを一から構築し、モダンなバックエンド開発の基礎を学習します。フレームワークに頼らず、Goの本質的な機能を理解することで、より深いバックエンド知識を身につけることができます。

### なぜ標準パッケージのみなのか？

**🟢 学習効果が高い理由:**
- Go言語の核となる概念を深く理解できる
- フレームワーク特有の「魔法」に惑わされない
- どのGoプロジェクトでも応用できる基礎力が身につく
- パフォーマンスとメモリ使用量を最適化できる

**🔴 フレームワーク依存の問題点:**
```go
// ❌ 悪い例：フレームワーク依存のコード
func CreateTodo(c *gin.Context) {
    // Ginのマジックメソッドに依存
    // 内部動作が理解しにくい
    var todo Todo
    c.ShouldBindJSON(&todo)  // 何が起きているか不明
}

// ✅ 良い例：標準パッケージでの明示的な処理
func CreateTodo(w http.ResponseWriter, r *http.Request) {
    // 各ステップが明確で理解しやすい
    var todo Todo
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&todo); err != nil {
        // エラーハンドリングも自分で制御
        writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON", err.Error())
        return
    }
}
```

## 📚 学習内容と重要な概念

### 🏗️ アーキテクチャ設計
- **Clean Architecture**: ビジネスロジックを外部依存から分離
- **依存関係逆転の原則**: インターフェースを使った疎結合設計
- **レイヤー分離**: Domain, Application, Infrastructure の責任分離

### 🌐 HTTP/Web開発
- **標準net/httpパッケージ**によるサーバー構築
- **手動ルーティング**の実装と理解
- **JSON API**の設計原則とベストプラクティス
- **RESTful API**の設計パターン

### 🗄️ データベース設計
- **database/sqlパッケージ**による生のSQL操作
- **Repository パターン**による抽象化
- **プリペアードステートメント**によるセキュリティ対策
- **トランザクション処理**の理解

### 🧪 テスト戦略
- **テスト駆動開発（TDD）**の実践
- **テーブル駆動テスト**パターン
- **モック実装**による単体テスト
- **統合テスト**による品質保証

### 🔒 セキュリティとベストプラクティス
- **SQLインジェクション対策**
- **エラーハンドリング**パターン
- **リソース管理**（メモリリーク防止）
- **Graceful Shutdown**の実装

### 💡 重要な専門用語

| 用語 | 説明 | 重要度 |
|------|------|--------|
| **Clean Architecture** | ビジネスロジックを外部依存から分離するアーキテクチャパターン | ⭐⭐⭐ |
| **Repository Pattern** | データアクセスを抽象化し、ビジネスロジックから分離するパターン | ⭐⭐⭐ |
| **DTO (Data Transfer Object)** | レイヤー間でデータを転送するためのオブジェクト | ⭐⭐⭐ |
| **依存性注入 (DI)** | 依存関係を外部から注入することで疎結合を実現する手法 | ⭐⭐⭐ |
| **プリペアードステートメント** | SQLインジェクションを防ぐためのSQL実行方法 | ⭐⭐⭐ |
| **ミドルウェア** | HTTPリクエスト処理の前後に共通処理を挟み込むパターン | ⭐⭐ |
| **コンテキスト (context.Context)** | リクエストスコープの値やキャンセレーション信号を伝播 | ⭐⭐ |
| **Graceful Shutdown** | 実行中の処理を適切に終了させてからサーバーを停止する手法 | ⭐⭐ |

## 🛠 前提条件

- Go 1.21+ がインストールされている
- SQLiteまたはMySQLの基本知識
- HTTPとREST APIの基本概念
- 基本的なGo言語の文法知識

## 📖 チュートリアル構成

### [Chapter 1: プロジェクト構成とClean Architecture](#chapter-1)
- プロジェクト構造の理解
- Clean Architectureの概念
- ディレクトリ構成の意味

### [Chapter 2: ドメイン層の実装](#chapter-2)
- エンティティの設計
- ビジネスロジックの実装
- バリデーション処理

### [Chapter 3: インフラストラクチャ層の実装](#chapter-3)
- データベース接続
- リポジトリパターンの実装
- SQL操作の実装

### [Chapter 4: アプリケーション層の実装](#chapter-4)
- DTOの設計
- HTTPハンドラーの実装
- リクエスト・レスポンス処理

### [Chapter 5: ミドルウェアの実装](#chapter-5)
- ログ出力ミドルウェア
- リクエストIDの生成
- パニック回復処理

### [Chapter 6: テストの実装](#chapter-6)
- ユニットテストの書き方
- モックの実装
- 統合テストの作成

### [Chapter 7: サーバーの起動と統合](#chapter-7)
- HTTPサーバーの設定
- ルーティングの実装
- 依存性注入

---

## Chapter 1: プロジェクト構成とClean Architecture

### 1.1 プロジェクト構造の理解

まず、プロジェクトの全体構造を把握しましょう。以下のディレクトリ構造を作成してください：

```
todoapp-api-golang/
├── cmd/
│   └── api/
│       └── main.go                 # アプリケーションのエントリーポイント
├── internal/
│   ├── domain/                     # ドメイン層：ビジネスロジック
│   │   ├── entity/
│   │   │   ├── todo.go
│   │   │   └── todo_test.go
│   │   ├── repository/
│   │   │   └── todo_repository.go
│   │   └── service/
│   │       ├── todo_service.go
│   │       ├── todo_service_interface.go
│   │       └── todo_service_test.go
│   ├── application/                # アプリケーション層：HTTP処理
│   │   ├── dto/
│   │   │   ├── todo_request.go
│   │   │   ├── todo_response.go
│   │   │   └── todo_dto_test.go
│   │   ├── handler/
│   │   │   ├── todo_handler.go
│   │   │   └── todo_handler_test.go
│   │   └── middleware/
│   │       ├── middleware.go
│   │       └── middleware_test.go
│   └── infrastructure/             # インフラストラクチャ層：外部依存
│       ├── database/
│       │   ├── connection.go
│       │   ├── todo_repository_impl.go
│       │   └── todo_repository_impl_test.go
│       └── web/
│           ├── server.go
│           └── routes.go
├── pkg/
│   └── config/
│       └── config.go
├── go.mod
├── go.sum
├── .air.toml                       # ホットリロード設定
└── CLAUDE.md                       # Claude Code向けガイド
```

### 1.2 Clean Architectureの原則と実装パターン

Clean Architectureでは以下の原則を守ります：

#### 🎯 核となる4つの原則

1. **依存関係の方向**: 外側の層は内側の層に依存するが、逆は禁止
2. **ドメインの独立性**: ビジネスロジックは外部の詳細から独立
3. **インターフェースの活用**: 抽象化によって疎結合を実現
4. **単一責任の原則**: 各層は明確に定義された責任を持つ

```
┌─────────────────────────────────────┐
│           Infrastructure            │  ← 最外層（データベース、HTTP等）
│  ┌─────────────────────────────────┐ │
│  │          Application           │ │  ← アプリケーション層（ハンドラー、DTO）
│  │  ┌─────────────────────────────┐ │ │
│  │  │           Domain            │ │ │  ← ドメイン層（エンティティ、サービス）
│  │  │                             │ │ │  ← 最内層（ビジネスロジック）
│  │  └─────────────────────────────┘ │ │
│  └─────────────────────────────────┘ │
└─────────────────────────────────────┘
```

#### ❌ よくある設計上の間違い

```go
// 🚫 アンチパターン：レイヤー間の不適切な依存関係
package entity

import (
    "database/sql"  // ❌ ドメイン層がインフラ層に依存
    "net/http"      // ❌ ドメイン層がアプリケーション層に依存
)

type Todo struct {
    ID          int
    Title       string
    db          *sql.DB           // ❌ エンティティがDBに直接依存
    httpRequest *http.Request     // ❌ エンティティがHTTPに依存
}

// ❌ ドメインエンティティにインフラ層の処理を書いてしまう
func (t *Todo) Save() error {
    query := "INSERT INTO todos (title) VALUES (?)"
    _, err := t.db.Exec(query, t.Title)  // ❌ SQL処理がエンティティに混入
    return err
}
```

```go
// ✅ 正しいパターン：適切なレイヤー分離
package entity

// ✅ エンティティは純粋なドメインオブジェクト
// 外部依存を持たない
type Todo struct {
    ID          int
    Title       string
    Description string
    IsCompleted bool
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

// ✅ ビジネスロジックのみに集中
func (t *Todo) IsValid() bool {
    return len(strings.TrimSpace(t.Title)) > 0 && len(t.Title) <= 100
}

func (t *Todo) MarkAsCompleted() {
    t.IsCompleted = true
    t.UpdatedAt = time.Now()
}
```

#### 🔄 データフローの理解

正しいClean Architectureでのデータフローは以下のようになります：

```
1. HTTPリクエスト → Infrastructure Layer (Web)
2. Infrastructure → Application Layer (Handler)
3. Handler → DTO → Domain Layer (Service)
4. Service → Repository Interface (Domain)
5. Repository Interface ← Repository Implementation (Infrastructure)
6. データベース操作
7. Repository Implementation → Repository Interface
8. Repository Interface → Service
9. Service → Handler
10. Handler → DTO → Infrastructure Layer
11. Infrastructure → HTTPレスポンス
```

#### 🎨 レイヤー責任の明確化

| レイヤー | 責任 | やってはいけないこと |
|----------|------|---------------------|
| **Domain** | ビジネスルール、ドメインロジック | HTTP処理、SQL処理、JSON処理 |
| **Application** | ユースケース実行、DTOマッピング | ビジネスルール実装、データベース直接操作 |
| **Infrastructure** | 外部システムとの通信 | ビジネスロジック実装 |

### 1.3 プロジェクトの初期化

`go.mod`ファイルを作成して、プロジェクトを初期化します：

```bash
go mod init todoapp-api-golang
```

必要な依存関係を追加：

```bash
go get github.com/mattn/go-sqlite3  # SQLite driver（テスト用）
go get github.com/go-sql-driver/mysql  # MySQL driver（本番用）
```

---

## Chapter 2: ドメイン層の実装

ドメイン層から実装を始めます。これがアプリケーションの心臓部となります。

### 2.1 Todoエンティティの実装

`internal/domain/entity/todo.go`を作成：

```go
package entity

import (
    "fmt"
    "strings"
    "time"
)

// Todo はTodoアイテムを表現するドメインエンティティです
// エンティティの役割：
// 1. ビジネスデータの構造定義
// 2. ドメイン固有の振る舞いの実装
// 3. ビジネスルールの検証
type Todo struct {
    // ID は一意識別子（データベースで自動生成）
    ID int `json:"id"`
    
    // Title はTodoのタイトル（必須、1-100文字）
    Title string `json:"title"`
    
    // Description はTodoの詳細説明（任意、最大500文字）
    Description string `json:"description"`
    
    // IsCompleted は完了状態（デフォルト: false）
    IsCompleted bool `json:"is_completed"`
    
    // CreatedAt は作成日時
    CreatedAt time.Time `json:"created_at"`
    
    // UpdatedAt は最終更新日時
    UpdatedAt time.Time `json:"updated_at"`
}

// IsValid はTodoエンティティの妥当性を検証します
// ビジネスルールの実装例
func (t *Todo) IsValid() bool {
    // タイトルの妥当性チェック
    trimmed := strings.TrimSpace(t.Title)
    if len(trimmed) == 0 || len(trimmed) > 100 {
        return false
    }
    
    // 説明文の長さチェック（任意項目だが、ある場合は500文字まで）
    if len(t.Description) > 500 {
        return false
    }
    
    return true
}

// MarkAsCompleted はTodoを完了状態にします
// ドメイン固有の操作を encapsulation
func (t *Todo) MarkAsCompleted() {
    t.IsCompleted = true
    t.UpdatedAt = time.Now()
}

// MarkAsIncomplete はTodoを未完了状態にします
func (t *Todo) MarkAsIncomplete() {
    t.IsCompleted = false
    t.UpdatedAt = time.Now()
}

// String はTodoの文字列表現を返します
// デバッグやログ出力時に便利
func (t *Todo) String() string {
    status := "未完了"
    if t.IsCompleted {
        status = "完了"
    }
    return fmt.Sprintf("Todo{ID: %d, Title: %s, Status: %s}", t.ID, t.Title, status)
}

// TableName はデータベースのテーブル名を返します
// ORMを使わない場合でも、テーブル名の管理に便利
func (t *Todo) TableName() string {
    return "todos"
}
```

#### 🎓 エンティティ実装の学習ポイント

**✅ 良い実装パターン:**
- エンティティは**データ**と**振る舞い**を持つ
- ビジネスルールはエンティティ内で検証
- メソッドはドメイン用語を使用（MarkAsCompleted等）
- 外部依存を持たない純粋なGoオブジェクト

**❌ よくある間違い:**

```go
// 🚫 アンチパターン：神オブジェクト（God Object）
type Todo struct {
    ID          int
    Title       string
    Description string
    // ❌ エンティティに過剰な責任を持たせる
    HTTPClient  *http.Client    // HTTP処理の責任
    Database    *sql.DB         // DB操作の責任  
    Logger      *log.Logger     // ログ出力の責任
    Validator   interface{}     // バリデーションの責任
}

// ❌ エンティティメソッドが複数の責任を持つ
func (t *Todo) SaveAndNotifyAndLog() error {
    // データベース保存
    if err := t.Database.Exec("INSERT..."); err != nil {
        return err
    }
    
    // 外部API通知  
    if err := t.HTTPClient.Post("..."); err != nil {
        return err
    }
    
    // ログ出力
    t.Logger.Printf("Todo saved: %d", t.ID)
    
    return nil
}
```

```go
// ✅ 良いパターン：単一責任のエンティティ
type Todo struct {
    ID          int
    Title       string  
    Description string
    IsCompleted bool
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

// ✅ 各メソッドが単一の責任を持つ
func (t *Todo) IsValid() bool {
    // バリデーションのみに集中
    return len(strings.TrimSpace(t.Title)) > 0 && len(t.Title) <= 100
}

func (t *Todo) MarkAsCompleted() {
    // 状態変更のみに集中
    t.IsCompleted = true
    t.UpdatedAt = time.Now()
}

func (t *Todo) GetPriority() string {
    // ビジネスロジックのみに集中
    if strings.Contains(strings.ToLower(t.Title), "urgent") {
        return "high"
    }
    return "normal"
}
```

#### 🔍 Goでのエンティティ設計のベストプラクティス

1. **構造体のフィールド設計**
```go
// ✅ 適切なフィールド設計
type Todo struct {
    ID          int       `json:"id" db:"id"`                    // 一意識別子
    Title       string    `json:"title" db:"title"`              // 必須フィールド
    Description string    `json:"description" db:"description"`  // 任意フィールド  
    IsCompleted bool      `json:"is_completed" db:"is_completed"`// 状態フィールド
    CreatedAt   time.Time `json:"created_at" db:"created_at"`    // 作成日時
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`    // 更新日時
}
```

2. **メソッドの命名規則**
```go
// ✅ 意図が明確なメソッド名
func (t *Todo) IsValid() bool         // 状態確認
func (t *Todo) MarkAsCompleted()      // アクション実行  
func (t *Todo) GetDisplayTitle()      // 計算されたプロパティ
func (t *Todo) CanBeDeleted() bool    // 権限チェック
```

### 2.2 リポジトリインターフェースの定義

`internal/domain/repository/todo_repository.go`を作成：

```go
package repository

import (
    "context"
    "todoapp-api-golang/internal/domain/entity"
)

// TodoRepository はTodoのデータアクセスを抽象化するインターフェースです
// インターフェースの役割：
// 1. データアクセスの抽象化
// 2. ドメイン層とインフラ層の分離
// 3. テスト時のモック実装を可能にする
// 4. データベース実装の詳細を隠蔽
type TodoRepository interface {
    // Create は新しいTodoをデータストアに保存します
    // context.Context はキャンセレーション、タイムアウト、値の伝播に使用
    Create(ctx context.Context, todo *entity.Todo) (*entity.Todo, error)
    
    // GetByID は指定されたIDのTodoを取得します
    // 存在しない場合はnil, errorを返します
    GetByID(ctx context.Context, id int) (*entity.Todo, error)
    
    // GetAll は全てのTodoを取得します
    // 大量データの場合は将来的にページング対応を検討
    GetAll(ctx context.Context) ([]*entity.Todo, error)
    
    // Update は既存のTodoを更新します
    // 存在しないIDの場合はerrorを返します
    Update(ctx context.Context, todo *entity.Todo) (*entity.Todo, error)
    
    // Delete は指定されたIDのTodoを削除します
    // 存在しないIDでもエラーにしない（冪等性）
    Delete(ctx context.Context, id int) error
}

// インターフェースの設計原則：
// 1. 小さく、焦点を絞ったインターフェース
// 2. 実装の詳細ではなく、振る舞いに着目
// 3. context.Contextを第一引数に取る（Go慣例）
// 4. エラーハンドリングを明示的に行う
```

#### 🎓 Repository パターンの学習ポイント

**🎯 Repository パターンの目的:**
- **データアクセスの抽象化**: SQLやNoSQLの詳細を隠蔽
- **テスタビリティの向上**: モック実装を容易にする  
- **技術の交換可能性**: データベースの変更に強い設計
- **ビジネスロジックの分離**: データ取得とビジネスルールを分離

**❌ よくある Repository の設計ミス:**

```go
// 🚫 アンチパターン：Generic Repository
type GenericRepository interface {
    Save(entity interface{}) error           // ❌ 型安全性がない
    FindById(id interface{}) interface{}     // ❌ 何でも取れてしまう
    FindAll() []interface{}                  // ❌ 型情報が失われる
    Delete(entity interface{}) error         // ❌ エラーの温床
}

// 🚫 アンチパターン：SQL漏洩Repository  
type TodoRepository interface {
    ExecuteSQL(query string, args ...interface{}) error  // ❌ SQLが外に漏れる
    FindBySQL(query string) ([]*Todo, error)            // ❌ 実装詳細が公開
}

// 🚫 アンチパターン：肥大化Repository
type TodoRepository interface {
    Create(todo *Todo) error
    Update(todo *Todo) error
    Delete(id int) error
    // ❌ 以下は別のRepositoryの責任
    CreateUser(user *User) error
    SendEmail(email *Email) error  
    CalculateStatistics() (*Stats, error)
}
```

```go
// ✅ 良いパターン：特化型Repository
type TodoRepository interface {
    // ✅ 型安全で明確なメソッド
    Create(ctx context.Context, todo *entity.Todo) (*entity.Todo, error)
    GetByID(ctx context.Context, id int) (*entity.Todo, error)
    GetAll(ctx context.Context) ([]*entity.Todo, error)
    Update(ctx context.Context, todo *entity.Todo) (*entity.Todo, error)
    Delete(ctx context.Context, id int) error
    
    // ✅ ドメイン特有のクエリメソッド（必要に応じて）
    GetCompletedTodos(ctx context.Context) ([]*entity.Todo, error)
    GetTodosByDateRange(ctx context.Context, from, to time.Time) ([]*entity.Todo, error)
}
```

**🔧 context.Context の重要性:**

```go
// context.Context が提供する機能
func ExampleContextUsage(ctx context.Context) {
    // 1. キャンセレーション
    select {
    case <-ctx.Done():
        return ctx.Err() // タイムアウトやキャンセル
    default:
        // 処理続行
    }
    
    // 2. タイムアウト設定
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    
    // 3. 値の伝播（使用は慎重に）
    userID := ctx.Value("userID")
}
```

**📝 Repository実装時の注意点:**

| ✅ 良い実践 | ❌ 避けるべき実践 |
|------------|------------------|
| 型安全なメソッドシグネチャ | `interface{}` の多用 |
| 明確なメソッド名 | 曖昧な命名（`Find`, `Get`の混在） |
| context.Context の使用 | グローバル変数の使用 |  
| 単一責任の原則 | 複数ドメインの混在 |
| エラーハンドリングの統一 | エラーの隠蔽や無視 |

### 2.3 ドメインサービスの実装

`internal/domain/service/todo_service_interface.go`を作成：

```go
package service

import (
    "context"
    "todoapp-api-golang/internal/domain/entity"
)

// TodoServiceInterface は Todo サービスのインターフェースです
// テスタビリティ向上のため、ハンドラー層のテストでモック実装を使用できます
type TodoServiceInterface interface {
    // CreateTodo は新しいTodoを作成します
    CreateTodo(ctx context.Context, todo *entity.Todo) (*entity.Todo, error)
    
    // GetTodoByID は指定されたIDのTodoを取得します
    GetTodoByID(ctx context.Context, id int) (*entity.Todo, error)
    
    // GetAllTodos は全てのTodoを取得します
    GetAllTodos(ctx context.Context) ([]*entity.Todo, error)
    
    // UpdateTodo は既存のTodoを更新します
    UpdateTodo(ctx context.Context, todo *entity.Todo) (*entity.Todo, error)
    
    // DeleteTodo は指定されたIDのTodoを削除します
    DeleteTodo(ctx context.Context, id int) error
    
    // CompleteTodo はTodoを完了状態にします
    CompleteTodo(ctx context.Context, id int) (*entity.Todo, error)
    
    // IncompleteTodo はTodoを未完了状態にします
    IncompleteTodo(ctx context.Context, id int) (*entity.Todo, error)
}
```

`internal/domain/service/todo_service.go`を作成：

```go
package service

import (
    "context"
    "fmt"
    "strings"
    
    "todoapp-api-golang/internal/domain/entity"
    "todoapp-api-golang/internal/domain/repository"
)

// TodoService はTodo関連のビジネスロジックを実装します
// サービス層の役割：
// 1. ビジネスロジックの実装
// 2. 複数のリポジトリの協調
// 3. トランザクション管理
// 4. ドメインルールの適用
type TodoService struct {
    todoRepository repository.TodoRepository
}

// NewTodoService はTodoServiceのコンストラクタです
// 依存性注入パターンの実装
func NewTodoService(todoRepository repository.TodoRepository) *TodoService {
    return &TodoService{
        todoRepository: todoRepository,
    }
}

// CreateTodo は新しいTodoを作成します
func (s *TodoService) CreateTodo(ctx context.Context, todo *entity.Todo) (*entity.Todo, error) {
    // 1. ビジネスルールの検証
    if todo == nil {
        return nil, fmt.Errorf("todo cannot be nil")
    }
    
    // タイトルの正規化（前後の空白を除去）
    todo.Title = strings.TrimSpace(todo.Title)
    
    // 2. エンティティレベルのバリデーション
    if !todo.IsValid() {
        return nil, fmt.Errorf("invalid todo: title must be 1-100 characters")
    }
    
    // 3. リポジトリ経由でデータ保存
    createdTodo, err := s.todoRepository.Create(ctx, todo)
    if err != nil {
        return nil, fmt.Errorf("failed to create todo: %w", err)
    }
    
    return createdTodo, nil
}

// GetTodoByID は指定されたIDのTodoを取得します
func (s *TodoService) GetTodoByID(ctx context.Context, id int) (*entity.Todo, error) {
    // IDの妥当性チェック（ビジネスルール）
    if id <= 0 {
        return nil, fmt.Errorf("invalid id: %d", id)
    }
    
    todo, err := s.todoRepository.GetByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("failed to get todo by id %d: %w", id, err)
    }
    
    return todo, nil
}

// GetAllTodos は全てのTodoを取得します
func (s *TodoService) GetAllTodos(ctx context.Context) ([]*entity.Todo, error) {
    todos, err := s.todoRepository.GetAll(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to get all todos: %w", err)
    }
    
    return todos, nil
}

// UpdateTodo は既存のTodoを更新します
func (s *TodoService) UpdateTodo(ctx context.Context, todo *entity.Todo) (*entity.Todo, error) {
    // 1. 基本的なバリデーション
    if todo == nil {
        return nil, fmt.Errorf("todo cannot be nil")
    }
    
    if todo.ID <= 0 {
        return nil, fmt.Errorf("invalid todo id: %d", todo.ID)
    }
    
    // 2. 既存データの存在確認
    existingTodo, err := s.todoRepository.GetByID(ctx, todo.ID)
    if err != nil {
        return nil, fmt.Errorf("todo not found: %w", err)
    }
    if existingTodo == nil {
        return nil, fmt.Errorf("todo with id %d not found", todo.ID)
    }
    
    // 3. タイトルの正規化
    todo.Title = strings.TrimSpace(todo.Title)
    
    // 4. ビジネスルールの検証
    if !todo.IsValid() {
        return nil, fmt.Errorf("invalid todo: title must be 1-100 characters")
    }
    
    // 5. 更新実行
    updatedTodo, err := s.todoRepository.Update(ctx, todo)
    if err != nil {
        return nil, fmt.Errorf("failed to update todo: %w", err)
    }
    
    return updatedTodo, nil
}

// DeleteTodo は指定されたIDのTodoを削除します
func (s *TodoService) DeleteTodo(ctx context.Context, id int) error {
    // IDの妥当性チェック
    if id <= 0 {
        return fmt.Errorf("invalid id: %d", id)
    }
    
    err := s.todoRepository.Delete(ctx, id)
    if err != nil {
        return fmt.Errorf("failed to delete todo: %w", err)
    }
    
    return nil
}

// CompleteTodo はTodoを完了状態にします
func (s *TodoService) CompleteTodo(ctx context.Context, id int) (*entity.Todo, error) {
    // 1. 既存Todoの取得
    todo, err := s.GetTodoByID(ctx, id)
    if err != nil {
        return nil, err // エラーメッセージは GetTodoByID で設定済み
    }
    
    // 2. 既に完了している場合はそのまま返す（冪等性）
    if todo.IsCompleted {
        return todo, nil
    }
    
    // 3. 完了状態に変更
    todo.MarkAsCompleted()
    
    // 4. 更新実行
    updatedTodo, err := s.todoRepository.Update(ctx, todo)
    if err != nil {
        return nil, fmt.Errorf("failed to complete todo: %w", err)
    }
    
    return updatedTodo, nil
}

// IncompleteTodo はTodoを未完了状態にします
func (s *TodoService) IncompleteTodo(ctx context.Context, id int) (*entity.Todo, error) {
    // 1. 既存Todoの取得
    todo, err := s.GetTodoByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // 2. 既に未完了の場合はそのまま返す（冪等性）
    if !todo.IsCompleted {
        return todo, nil
    }
    
    // 3. 未完了状態に変更
    todo.MarkAsIncomplete()
    
    // 4. 更新実行
    updatedTodo, err := s.todoRepository.Update(ctx, todo)
    if err != nil {
        return nil, fmt.Errorf("failed to incomplete todo: %w", err)
    }
    
    return updatedTodo, nil
}

// コンパイル時インターフェース実装確認
// この行により、TodoService が TodoServiceInterface を実装していることを
// コンパイル時に確認できます
var _ TodoServiceInterface = (*TodoService)(nil)
```

#### 🎓 ドメインサービス実装の学習ポイント

**🎯 ドメインサービスの役割:**
- **複雑なビジネスロジック**: 単一エンティティでは表現できないロジック
- **複数エンティティの協調**: エンティティ間の相互作用を管理
- **トランザクション境界**: データ整合性を保つための処理単位
- **ビジネスルールの実装**: ドメインエキスパートが定義したルール

**❌ よくあるサービス層の設計ミス:**

```go
// 🚫 アンチパターン：貧弱ドメインモデル（Anemic Domain Model）
type TodoService struct {
    todoRepository TodoRepository
}

// ❌ サービスが全ての処理を行い、エンティティは単なるデータホルダー
func (s *TodoService) CompleteTodo(ctx context.Context, id int) error {
    // ❌ ビジネスロジックがサービスに集中
    todo, err := s.todoRepository.GetByID(ctx, id)
    if err != nil {
        return err
    }
    
    // ❌ エンティティのメソッドを使わずに直接フィールドを操作
    todo.IsCompleted = true
    todo.UpdatedAt = time.Now()
    
    // ❌ バリデーションもサービス側で実装
    if len(todo.Title) == 0 {
        return errors.New("title is required")
    }
    
    return s.todoRepository.Update(ctx, todo)
}
```

```go
// 🚫 アンチパターン：God Service（何でもサービス）
type TodoService struct {
    todoRepository TodoRepository
    userRepository UserRepository
    emailService   EmailService
    fileService    FileService
    // ❌ 過剰な依存関係
}

func (s *TodoService) ProcessTodoCreation(ctx context.Context, todo *Todo) error {
    // ❌ 複数の責任を一つのメソッドで処理
    
    // Todo作成
    if err := s.todoRepository.Create(ctx, todo); err != nil {
        return err
    }
    
    // ユーザー更新
    user, _ := s.userRepository.GetByID(ctx, todo.UserID)
    user.TodoCount++
    s.userRepository.Update(ctx, user)
    
    // メール送信
    s.emailService.SendTodoCreatedEmail(user.Email, todo)
    
    // ファイル作成
    s.fileService.CreateTodoBackup(todo)
    
    // ❌ トランザクション管理も不適切
    return nil
}
```

```go
// ✅ 良いパターン：豊富なドメインモデル（Rich Domain Model）
type TodoService struct {
    todoRepository TodoRepository
}

func (s *TodoService) CompleteTodo(ctx context.Context, id int) (*entity.Todo, error) {
    // 1. エンティティ取得
    todo, err := s.todoRepository.GetByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("failed to get todo: %w", err)
    }
    
    if todo == nil {
        return nil, fmt.Errorf("todo not found: %d", id)
    }
    
    // 2. ✅ エンティティのメソッドを使用（ビジネスロジックは適切な場所に）
    if todo.IsCompleted {
        return todo, nil // 冪等性を保つ
    }
    
    // 3. ✅ ドメインオブジェクトのメソッドを活用
    todo.MarkAsCompleted()
    
    // 4. ✅ エンティティレベルのバリデーション
    if !todo.IsValid() {
        return nil, fmt.Errorf("invalid todo state")
    }
    
    // 5. 永続化
    updatedTodo, err := s.todoRepository.Update(ctx, todo)
    if err != nil {
        return nil, fmt.Errorf("failed to update todo: %w", err)
    }
    
    return updatedTodo, nil
}
```

**⚡ エラーハンドリングのベストプラクティス:**

```go
// ✅ 適切なエラーハンドリング
func (s *TodoService) CreateTodo(ctx context.Context, todo *entity.Todo) (*entity.Todo, error) {
    // 1. バリデーション用のエラー
    if todo == nil {
        return nil, fmt.Errorf("todo cannot be nil")
    }
    
    todo.Title = strings.TrimSpace(todo.Title)
    
    // 2. ビジネスルール違反のエラー  
    if !todo.IsValid() {
        return nil, fmt.Errorf("validation failed: title must be 1-100 characters")
    }
    
    // 3. インフラストラクチャ層のエラーをラップ
    createdTodo, err := s.todoRepository.Create(ctx, todo)
    if err != nil {
        return nil, fmt.Errorf("failed to create todo: %w", err)
    }
    
    return createdTodo, nil
}
```

**🔄 冪等性の重要性:**

```go
// ✅ 冪等性を保つ実装
func (s *TodoService) CompleteTodo(ctx context.Context, id int) (*entity.Todo, error) {
    todo, err := s.GetTodoByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // ✅ 既に完了済みの場合は何もしない（冪等性）
    if todo.IsCompleted {
        return todo, nil  // 同じリクエストを複数回実行しても安全
    }
    
    todo.MarkAsCompleted()
    return s.todoRepository.Update(ctx, todo)
}
```

**📏 適切なサービス粒度の判断:**

| ✅ サービスレイヤーで扱うべき | ❌ サービスレイヤーで扱うべきでない |
|----------------------------|-----------------------------------|
| 複数エンティティの協調処理 | 単一エンティティの状態変更 |
| トランザクション境界の管理 | 単純なCRUD操作 |
| 複雑なビジネスルール | データフォーマットの変換 |
| 外部システムとの整合性 | HTTP特有の処理 |

---

## Chapter 3: インフラストラクチャ層の実装

次に、外部システムとの接続を担当するインフラストラクチャ層を実装します。

### 3.1 データベース接続の設定

`internal/infrastructure/database/connection.go`を作成：

```go
package database

import (
    "database/sql"
    "fmt"
    "log"
    
    _ "github.com/go-sql-driver/mysql"
    _ "github.com/mattn/go-sqlite3"
)

// DatabaseManager はデータベース接続を管理する構造体です
// 標準パッケージでのデータベース管理の学習ポイント：
// 1. sql.DB の使用方法
// 2. 接続プールの管理
// 3. トランザクションの扱い
// 4. プリペアードステートメントの活用
type DatabaseManager struct {
    DB *sql.DB
}

// NewDatabaseManager はDatabaseManagerのコンストラクタです
func NewDatabaseManager() *DatabaseManager {
    return &DatabaseManager{}
}

// ConnectMySQL はMySQLデータベースに接続します
// 本番環境での使用を想定
func (dm *DatabaseManager) ConnectMySQL(dsn string) error {
    // sql.Open は接続プールを作成（実際の接続はまだ行われない）
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return fmt.Errorf("failed to open mysql connection: %w", err)
    }
    
    // 実際の接続をテストする
    if err := db.Ping(); err != nil {
        return fmt.Errorf("failed to ping mysql database: %w", err)
    }
    
    // 接続プールの設定
    db.SetMaxOpenConns(25)     // 最大接続数
    db.SetMaxIdleConns(25)     // アイドル接続数
    db.SetConnMaxLifetime(300) // 接続の最大生存時間（秒）
    
    dm.DB = db
    log.Printf("Successfully connected to MySQL database")
    return nil
}

// ConnectSQLite はSQLiteデータベースに接続します
// 開発・テスト環境での使用を想定
func (dm *DatabaseManager) ConnectSQLite(filepath string) error {
    // SQLiteは軽量なファイルベースDB
    db, err := sql.Open("sqlite3", filepath)
    if err != nil {
        return fmt.Errorf("failed to open sqlite connection: %w", err)
    }
    
    if err := db.Ping(); err != nil {
        return fmt.Errorf("failed to ping sqlite database: %w", err)
    }
    
    dm.DB = db
    log.Printf("Successfully connected to SQLite database: %s", filepath)
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
        )
    `
    
    _, err := dm.DB.Exec(createTodosTable)
    if err != nil {
        return fmt.Errorf("failed to create todos table: %w", err)
    }
    
    log.Println("Database tables created successfully")
    return nil
}

// Close はデータベース接続を閉じます
// アプリケーション終了時のクリーンアップ
func (dm *DatabaseManager) Close() error {
    if dm.DB != nil {
        return dm.DB.Close()
    }
    return nil
}

// GetDB はデータベース接続を返します
// 他のコンポーネントがDBアクセスする際に使用
func (dm *DatabaseManager) GetDB() *sql.DB {
    return dm.DB
}
```

### 3.2 リポジトリ実装

`internal/infrastructure/database/todo_repository_impl.go`を作成：

```go
package database

import (
    "context"
    "database/sql"
    "fmt"
    "time"
    
    "todoapp-api-golang/internal/domain/entity"
    "todoapp-api-golang/internal/domain/repository"
)

// todoRepositoryImpl はTodoRepositoryインターフェースの実装です
// 標準パッケージでのCRUD操作実装の学習ポイント：
// 1. database/sql パッケージの使用方法
// 2. プリペアードステートメントでのSQLインジェクション対策
// 3. エラーハンドリングのパターン
// 4. トランザクション処理
type todoRepositoryImpl struct {
    db *sql.DB
}

// NewTodoRepository はtodoRepositoryImplのコンストラクタです
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
        VALUES (?, ?, false, datetime('now'), datetime('now'))
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

// GetByID は指定されたIDのTodoを取得します
// 標準パッケージを使ったSELECT操作の学習
func (r *todoRepositoryImpl) GetByID(ctx context.Context, id int) (*entity.Todo, error) {
    // IDの妥当性チェック
    if id <= 0 {
        return nil, fmt.Errorf("invalid id: %d", id)
    }
    
    // 1. SELECT用のSQL文を定義
    query := `
        SELECT id, title, description, is_completed, created_at, updated_at
        FROM todos
        WHERE id = ?
    `
    
    // 2. QueryRowContext で単一行取得
    // QueryRowContext は1行だけ返すクエリ用
    row := r.db.QueryRowContext(ctx, query, id)
    
    // 3. 結果をスキャンしてTodoエンティティに変換
    todo := &entity.Todo{}
    var createdAt, updatedAt string
    
    err := row.Scan(
        &todo.ID,
        &todo.Title,
        &todo.Description,
        &todo.IsCompleted,
        &createdAt,
        &updatedAt,
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil // 見つからない場合はnilを返す（エラーではない）
        }
        return nil, fmt.Errorf("failed to scan todo: %w", err)
    }
    
    // 4. 時刻文字列をtime.Timeに変換
    todo.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt)
    if err != nil {
        todo.CreatedAt = time.Now() // パースエラー時はデフォルト値
    }
    
    todo.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", updatedAt)
    if err != nil {
        todo.UpdatedAt = time.Now()
    }
    
    return todo, nil
}

// GetAll は全てのTodoを取得します
// 標準パッケージを使った複数行SELECT操作の学習
func (r *todoRepositoryImpl) GetAll(ctx context.Context) ([]*entity.Todo, error) {
    // 1. 全件取得のSQL（作成日時の降順）
    query := `
        SELECT id, title, description, is_completed, created_at, updated_at
        FROM todos
        ORDER BY created_at DESC
    `
    
    // 2. QueryContext で複数行取得
    // QueryContext は複数行を返すクエリ用
    rows, err := r.db.QueryContext(ctx, query)
    if err != nil {
        return nil, fmt.Errorf("failed to query todos: %w", err)
    }
    defer rows.Close() // 必ずCloseを呼ぶ（リソースリーク防止）
    
    // 3. 結果をスライスに変換
    var todos []*entity.Todo
    
    for rows.Next() {
        todo := &entity.Todo{}
        var createdAt, updatedAt string
        
        err := rows.Scan(
            &todo.ID,
            &todo.Title,
            &todo.Description,
            &todo.IsCompleted,
            &createdAt,
            &updatedAt,
        )
        
        if err != nil {
            return nil, fmt.Errorf("failed to scan todo: %w", err)
        }
        
        // 時刻変換
        todo.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
        todo.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)
        
        todos = append(todos, todo)
    }
    
    // 4. イテレーション中のエラーをチェック
    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("error occurred during rows iteration: %w", err)
    }
    
    return todos, nil
}

// Update は既存のTodoを更新します
// 標準パッケージを使ったUPDATE操作と影響行数の確認を学習
func (r *todoRepositoryImpl) Update(ctx context.Context, todo *entity.Todo) (*entity.Todo, error) {
    // 1. UPDATE用のSQL文を定義
    // updated_at は現在時刻で自動更新
    query := `
        UPDATE todos
        SET title = ?, description = ?, is_completed = ?, updated_at = datetime('now')
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
    
    // 3. 影響を受けた行数をチェック
    // RowsAffected で実際に更新された行数を確認
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return nil, fmt.Errorf("failed to get affected rows: %w", err)
    }
    
    if rowsAffected == 0 {
        return nil, fmt.Errorf("todo with id %d not found", todo.ID)
    }
    
    // 4. 更新後のデータを取得して返す
    // updated_atが自動更新されているため、最新データを取得
    updatedTodo, err := r.GetByID(ctx, todo.ID)
    if err != nil {
        return nil, fmt.Errorf("failed to get updated todo: %w", err)
    }
    
    return updatedTodo, nil
}

// Delete は指定されたIDのTodoを削除します
// 標準パッケージを使ったDELETE操作の学習
func (r *todoRepositoryImpl) Delete(ctx context.Context, id int) error {
    // IDの妥当性チェック
    if id <= 0 {
        return fmt.Errorf("invalid id: %d", id)
    }
    
    // 1. DELETE用のSQL文を定義
    query := `DELETE FROM todos WHERE id = ?`
    
    // 2. DELETE実行
    result, err := r.db.ExecContext(ctx, query, id)
    if err != nil {
        return fmt.Errorf("failed to delete todo: %w", err)
    }
    
    // 3. 影響を受けた行数をチェック（オプション）
    // 削除操作は冪等性を保つため、存在しないIDでもエラーにしない場合もある
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get affected rows: %w", err)
    }
    
    // ログ出力のみでエラーにはしない（冪等性のため）
    if rowsAffected == 0 {
        // log.Printf("No todo found with id %d for deletion", id)
    }
    
    return nil
}

// コンパイル時インターフェース実装確認
var _ repository.TodoRepository = (*todoRepositoryImpl)(nil)
```

#### 🎓 database/sql パッケージの学習ポイント

**🎯 標準パッケージを使う理由:**
- **軽量性**: 外部依存が最小限で高速
- **制御性**: SQL処理の詳細を完全にコントロール
- **学習効果**: データベース操作の仕組みを深く理解
- **パフォーマンス**: オーバーヘッドが少ない

**❌ よくあるSQL操作の間違い:**

```go
// 🚫 危険なパターン：SQLインジェクション脆弱性
func (r *badRepository) GetByID(ctx context.Context, id string) (*entity.Todo, error) {
    // 直接文字列結合 - 絶対にやってはいけない
    query := fmt.Sprintf("SELECT * FROM todos WHERE id = %s", id)
    row := r.db.QueryRowContext(ctx, query) // SQLインジェクション脆弱性
}

// ✅ 正しいパターン：プリペアードステートメント
func (r *todoRepository) GetByID(ctx context.Context, id int) (*entity.Todo, error) {
    query := "SELECT * FROM todos WHERE id = ?"
    row := r.db.QueryRowContext(ctx, query, id) // 安全
}

// 🚫 リソースリーク：rows.Close()を忘れる
func (r *badRepository) GetAll(ctx context.Context) ([]*entity.Todo, error) {
    rows, err := r.db.QueryContext(ctx, "SELECT * FROM todos")
    if err != nil {
        return nil, err
    }
    // defer rows.Close() を忘れる -> コネクションリーク
    
    var todos []*entity.Todo
    for rows.Next() {
        // ... スキャン処理
    }
    return todos, nil
}

// ✅ 正しいパターン：確実なリソース解放
func (r *todoRepository) GetAll(ctx context.Context) ([]*entity.Todo, error) {
    rows, err := r.db.QueryContext(ctx, "SELECT * FROM todos")
    if err != nil {
        return nil, err
    }
    defer rows.Close() // 必須：リソースリーク防止
}

// 🚫 エラーハンドリングの不備
func (r *badRepository) Update(ctx context.Context, todo *entity.Todo) error {
    result, _ := r.db.ExecContext(ctx, query, todo.Title, todo.ID) // エラー無視
    // 影響行数チェックなし - 更新されたかわからない
    return nil
}

// ✅ 正しいパターン：適切なエラーハンドリング
func (r *todoRepository) Update(ctx context.Context, todo *entity.Todo) (*entity.Todo, error) {
    result, err := r.db.ExecContext(ctx, query, todo.Title, todo.ID)
    if err != nil {
        return nil, fmt.Errorf("failed to update: %w", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return nil, fmt.Errorf("failed to check affected rows: %w", err)
    }
    
    if rowsAffected == 0 {
        return nil, errors.New("todo not found")
    }
    
    return r.GetByID(ctx, todo.ID)
}
```

**📚 技術用語解説：**

| 用語 | 意味 | 重要度 |
|------|------|-------|
| SQLインジェクション | 悪意のあるSQL文を注入して不正な操作を行う攻撃手法 | ★★★ |
| プリペアードステートメント | SQL文とデータを分離してセキュリティを確保する仕組み | ★★★ |
| リソースリーク | メモリやコネクションが適切に解放されずに残る問題 | ★★★ |
| Context | タイムアウトやキャンセル処理を制御するGo標準の仕組み | ★★☆ |
| RowsAffected | UPDATE/DELETE文で実際に影響を受けた行数 | ★★☆ |

```go
// ✅ 安全で適切な実装パターン
func (r *todoRepositoryImpl) GetByTitle(ctx context.Context, title string) (*entity.Todo, error) {
    // ✅ プリペアードステートメントでSQLインジェクション対策
    query := `SELECT id, title, description, is_completed, created_at, updated_at 
              FROM todos WHERE title = ?`
    
    row := r.db.QueryRowContext(ctx, query, title) // ✅ context付きでタイムアウト対応
    
    todo := &entity.Todo{}
    var createdAt, updatedAt string
    
    err := row.Scan(
        &todo.ID,
        &todo.Title,
        &todo.Description, 
        &todo.IsCompleted,
        &createdAt,
        &updatedAt,
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil // ✅ 見つからない場合とエラーを区別
        }
        return nil, fmt.Errorf("failed to scan todo: %w", err) // ✅ エラー詳細を保持
    }
    
    // ✅ 時刻変換処理も適切に
    todo.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
    todo.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)
    
    return todo, nil
}

// ✅ 適切なリソース管理
func (r *todoRepositoryImpl) GetAll(ctx context.Context) ([]*entity.Todo, error) {
    query := `SELECT id, title, description, is_completed, created_at, updated_at
              FROM todos ORDER BY created_at DESC`
    
    rows, err := r.db.QueryContext(ctx, query) // ✅ context付き
    if err != nil {
        return nil, fmt.Errorf("failed to query todos: %w", err)
    }
    defer rows.Close() // ✅ 必ずリソースをクリーンアップ
    
    var todos []*entity.Todo
    
    for rows.Next() {
        todo := &entity.Todo{}
        var createdAt, updatedAt string
        
        err := rows.Scan(
            &todo.ID,
            &todo.Title,
            &todo.Description,
            &todo.IsCompleted,
            &createdAt,
            &updatedAt,
        )
        
        if err != nil {
            return nil, fmt.Errorf("failed to scan todo: %w", err)
        }
        
        // 時刻変換
        todo.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
        todo.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)
        
        todos = append(todos, todo)
    }
    
    // ✅ イテレーション中のエラーもチェック
    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("error occurred during rows iteration: %w", err)
    }
    
    return todos, nil
}
```

**🔐 セキュリティのベストプラクティス:**

```go
// ✅ プリペアードステートメントの正しい使用
func (r *todoRepositoryImpl) SearchTodos(ctx context.Context, searchTerm string, status *bool) ([]*entity.Todo, error) {
    // 動的クエリもプリペアードステートメントで安全に
    var args []interface{}
    query := `SELECT id, title, description, is_completed, created_at, updated_at 
              FROM todos WHERE 1=1`
    
    if searchTerm != "" {
        query += " AND title LIKE ?"
        args = append(args, "%"+searchTerm+"%")
    }
    
    if status != nil {
        query += " AND is_completed = ?"  
        args = append(args, *status)
    }
    
    query += " ORDER BY created_at DESC"
    
    rows, err := r.db.QueryContext(ctx, query, args...)
    // ... 以下省略
}
```

**⚡ パフォーマンス最適化のコツ:**

| 最適化項目 | ❌ 避けるべき | ✅ 推奨される |
|------------|--------------|---------------|
| **SELECT文** | `SELECT *` | 必要なカラムのみ指定 |
| **インデックス** | インデックスなし | WHERE句のカラムにインデックス |
| **接続管理** | 毎回新しい接続 | 接続プールの活用 |
| **トランザクション** | 長時間のトランザクション | 短時間での commit/rollback |
| **プリペアード文** | 毎回SQL文を準備 | 使い回し可能な文の準備 |

**🔧 エラーハンドリングのパターン:**

```go
// ✅ 適切なエラー分類と処理
func (r *todoRepositoryImpl) Update(ctx context.Context, todo *entity.Todo) (*entity.Todo, error) {
    query := `UPDATE todos 
              SET title = ?, description = ?, is_completed = ?, updated_at = datetime('now')
              WHERE id = ?`
    
    result, err := r.db.ExecContext(ctx, query,
        todo.Title,
        todo.Description, 
        todo.IsCompleted,
        todo.ID,
    )
    
    if err != nil {
        // ✅ データベースエラーを適切にラップ
        return nil, fmt.Errorf("failed to update todo: %w", err)
    }
    
    // ✅ 影響行数のチェック
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return nil, fmt.Errorf("failed to get affected rows: %w", err)
    }
    
    if rowsAffected == 0 {
        // ✅ ビジネスロジック的なエラーを区別
        return nil, fmt.Errorf("todo with id %d not found", todo.ID)
    }
    
    // 更新後のデータを取得
    return r.GetByID(ctx, todo.ID)
}
```

---

## Chapter 4: アプリケーション層の実装

アプリケーション層では、HTTPリクエストの処理とレスポンスの生成を行います。

### 4.1 DTO（Data Transfer Object）の実装

`internal/application/dto/todo_request.go`を作成：

```go
package dto

// CreateTodoRequest はTodo作成時のHTTPリクエストボディを表すDTO（Data Transfer Object）です
// DTOの役割：
// 1. HTTPリクエスト/レスポンスの構造を定義
// 2. 外部システム（クライアント）とのデータ交換フォーマット
// 3. ドメインエンティティとの変換（マッピング）
// 4. 入力値の基本的なバリデーション
type CreateTodoRequest struct {
    // Title はTodoのタイトル（必須項目）
    // `json:"title"` : JSONキー名を指定（Goの命名規則と異なる場合に使用）
    // `binding:"required"` : Ginフレームワークのバリデーションタグ（必須チェック）
    // `validate:"required,min=1,max=100"` : validator パッケージのバリデーション
    Title string `json:"title" binding:"required" validate:"required,min=1,max=100"`
    
    // Description はTodoの詳細説明（任意項目）
    // バリデーションは設定していませんが、長さ制限を設ける場合は
    // `validate:"max=500"` などを追加できます
    Description string `json:"description" validate:"max=500"`
}

// UpdateTodoRequest はTodo更新時のHTTPリクエストボディを表すDTOです
// 作成時とは異なり、全てのフィールドが任意更新可能な設計にしています
// （部分更新：PATCHメソッドの考え方）
type UpdateTodoRequest struct {
    // Title の更新（任意）
    // ポインタ型 (*string) を使用することで、フィールドが送信されたかどうかを判別可能
    // nil の場合は更新しない、値がある場合は更新する
    Title *string `json:"title,omitempty" validate:"omitempty,min=1,max=100"`
    
    // Description の更新（任意）
    Description *string `json:"description,omitempty" validate:"omitempty,max=500"`
    
    // IsCompleted の更新（任意）
    // bool のポインタ型で、完了状態の変更を任意にします
    IsCompleted *bool `json:"is_completed,omitempty"`
}

// CompleteTodoRequest はTodo完了/未完了切り替え専用のリクエストです
// シンプルなアクション用のDTOとして定義
type CompleteTodoRequest struct {
    // IsCompleted で完了状態を指定
    // true: 完了, false: 未完了
    IsCompleted bool `json:"is_completed" binding:"required"`
}

// TodoListRequest はTodo一覧取得時のクエリパラメータを表すDTOです
// 将来的な拡張（ページング、フィルタリング、ソート）を想定した構造
type TodoListRequest struct {
    // ページング関連（将来的な拡張用）
    // Page は取得するページ番号（1から開始）
    Page int `form:"page" validate:"min=1"`
    
    // Limit は1ページあたりの取得件数
    Limit int `form:"limit" validate:"min=1,max=100"`
    
    // フィルタリング関連（将来的な拡張用）
    // IsCompleted で完了状態によるフィルタ（任意）
    // nil の場合は全て、true/false で絞り込み
    IsCompleted *bool `form:"is_completed"`
    
    // ソート関連（将来的な拡張用）
    // SortBy はソートするフィールド名
    SortBy string `form:"sort_by" validate:"omitempty,oneof=id title created_at updated_at"`
    
    // SortOrder はソート順序（asc/desc）
    SortOrder string `form:"sort_order" validate:"omitempty,oneof=asc desc"`
}

// DTOのバリデーションルール解説：
//
// binding:"required" - Ginのバリデーション（必須）
// validate:"required" - validator パッケージのバリデーション（必須）
// validate:"min=1,max=100" - 最小1文字、最大100文字
// validate:"omitempty" - 空の場合はバリデーションをスキップ
// validate:"oneof=asc desc" - 指定した値のいずれかのみ許可
// json:"field_name,omitempty" - 空の場合はJSONに含めない
// form:"field_name" - URLクエリパラメータやフォームデータのキー名
```

`internal/application/dto/todo_response.go`を作成：

```go
package dto

import (
    "time"
    
    "todoapp-api-golang/internal/domain/entity"
)

// TodoResponse はTodo情報をクライアントに返すためのレスポンスDTOです
// レスポンスDTOの役割：
// 1. 外部に公開する情報の制御（セキュリティ）
// 2. クライアントに最適化されたデータ構造の提供
// 3. APIのバージョニング対応
// 4. 内部実装の隠蔽（エンティティの変更がAPIに影響しないようにする）
type TodoResponse struct {
    // ID はTodoの一意識別子
    ID int `json:"id"`
    
    // Title はTodoのタイトル
    Title string `json:"title"`
    
    // Description はTodoの詳細説明
    Description string `json:"description"`
    
    // IsCompleted はTodoの完了状態
    IsCompleted bool `json:"is_completed"`
    
    // CreatedAt は作成日時（RFC3339形式でJSONシリアライズ）
    CreatedAt time.Time `json:"created_at"`
    
    // UpdatedAt は最終更新日時
    UpdatedAt time.Time `json:"updated_at"`
}

// TodoListResponse はTodo一覧取得時のレスポンスDTOです
// 将来的なページング情報なども含められる構造にしています
type TodoListResponse struct {
    // Todos はTodoのリスト
    Todos []TodoResponse `json:"todos"`
    
    // Meta はメタ情報（ページング等）
    Meta ListMetaResponse `json:"meta"`
}

// ListMetaResponse は一覧取得時のメタ情報を表すDTOです
// ページング情報や総件数など、一覧表示に必要な付加情報を含みます
type ListMetaResponse struct {
    // Total は総件数
    Total int `json:"total"`
    
    // Page は現在のページ番号
    Page int `json:"page"`
    
    // Limit は1ページあたりの表示件数
    Limit int `json:"limit"`
    
    // TotalPages は総ページ数
    TotalPages int `json:"total_pages"`
}

// ErrorResponse はエラー発生時のレスポンスDTOです
// 統一的なエラーレスポンス形式を提供します
type ErrorResponse struct {
    // Error はエラーメッセージ
    Error string `json:"error"`
    
    // Code はエラーコード（任意、アプリケーション固有のコード）
    Code string `json:"code,omitempty"`
    
    // Details は詳細情報（バリデーションエラー等）
    Details interface{} `json:"details,omitempty"`
}

// ValidationErrorResponse はバリデーションエラー専用のレスポンスDTOです
type ValidationErrorResponse struct {
    // Error は基本エラーメッセージ
    Error string `json:"error"`
    
    // ValidationErrors はフィールド別のバリデーションエラー
    ValidationErrors []FieldError `json:"validation_errors"`
}

// FieldError はフィールド単位のバリデーションエラー情報です
type FieldError struct {
    // Field はエラーが発生したフィールド名
    Field string `json:"field"`
    
    // Message はエラーメッセージ
    Message string `json:"message"`
    
    // Value は入力された値（セキュリティ上問題ない場合のみ）
    Value interface{} `json:"value,omitempty"`
}

// --- 変換関数（Mapper functions） ---

// ToTodoResponse はEntityをResponseDTOに変換します
// エンティティ → レスポンスDTO の変換ロジック
func ToTodoResponse(todo *entity.Todo) TodoResponse {
    return TodoResponse{
        ID:          todo.ID,
        Title:       todo.Title,
        Description: todo.Description,
        IsCompleted: todo.IsCompleted,
        CreatedAt:   todo.CreatedAt,
        UpdatedAt:   todo.UpdatedAt,
    }
}

// ToTodoListResponse はEntity配列をResponseDTOに変換します
func ToTodoListResponse(todos []*entity.Todo, page, limit, total int) TodoListResponse {
    // Entity配列を Response配列に変換
    todoResponses := make([]TodoResponse, len(todos))
    for i, todo := range todos {
        todoResponses[i] = ToTodoResponse(todo)
    }
    
    // ページ数の計算
    totalPages := total / limit
    if total%limit != 0 {
        totalPages++
    }
    
    return TodoListResponse{
        Todos: todoResponses,
        Meta: ListMetaResponse{
            Total:      total,
            Page:       page,
            Limit:      limit,
            TotalPages: totalPages,
        },
    }
}

// ToEntity はリクエストDTOをEntityに変換します（Create用）
func (req CreateTodoRequest) ToEntity() *entity.Todo {
    return &entity.Todo{
        Title:       req.Title,
        Description: req.Description,
        // IsCompleted は新規作成時は常にfalse（デフォルト値）
        IsCompleted: false,
    }
}

// ApplyToEntity は更新リクエストDTOを既存Entityに適用します（Update用）
// nil チェックを行い、送信されたフィールドのみを更新します
func (req UpdateTodoRequest) ApplyToEntity(todo *entity.Todo) {
    // タイトルが送信された場合のみ更新
    if req.Title != nil {
        todo.Title = *req.Title
    }
    
    // 説明が送信された場合のみ更新
    if req.Description != nil {
        todo.Description = *req.Description
    }
    
    // 完了状態が送信された場合のみ更新
    if req.IsCompleted != nil {
        todo.IsCompleted = *req.IsCompleted
    }
}

// DTOパターンの利点：
// 1. セキュリティ: 内部IDやパスワードなど、外部に公開したくない情報を隠蔽
// 2. 進化性: APIの変更を内部実装の変更から分離
// 3. 最適化: クライアント要件に合わせたデータ構造の提供
// 4. バリデーション: 入力値の検証と制御
// 5. ドキュメント化: APIドキュメント生成のための明確な構造定義
```

### 4.2 HTTPハンドラーの実装

`internal/application/handler/todo_handler.go`を作成：

```go
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
    todoService service.TodoServiceInterface
}

// NewTodoHandler はTodoHandlerのコンストラクタです
// 標準パッケージを使った依存性注入の実装例
func NewTodoHandler(todoService service.TodoServiceInterface) *TodoHandler {
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
    
    // 2. Content-Type の確認
    contentType := r.Header.Get("Content-Type")
    if !strings.Contains(contentType, "application/json") {
        writeErrorResponse(w, http.StatusBadRequest, "Content-Type must be application/json", "")
        return
    }
    
    // 3. リクエストボディの解析
    var req dto.CreateTodoRequest
    decoder := json.NewDecoder(r.Body)
    decoder.DisallowUnknownFields() // 未知のフィールドを拒否
    
    if err := decoder.Decode(&req); err != nil {
        writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON format", err.Error())
        return
    }
    
    // 4. DTOからEntityへの変換
    todoEntity := req.ToEntity()
    
    // 5. ビジネスロジック実行
    createdTodo, err := h.todoService.CreateTodo(r.Context(), todoEntity)
    if err != nil {
        // エラーメッセージから適切なHTTPステータスを判定
        if strings.Contains(err.Error(), "invalid todo") {
            writeErrorResponse(w, http.StatusBadRequest, "Validation error", err.Error())
        } else {
            writeErrorResponse(w, http.StatusInternalServerError, "Failed to create todo", err.Error())
        }
        return
    }
    
    // 6. EntityからDTOへの変換とレスポンス
    response := dto.ToTodoResponse(createdTodo)
    writeJSONResponse(w, http.StatusCreated, response)
}

// GetAllTodos は全てのTodoを取得するHTTPハンドラーです
// GET /api/v1/todos へのリクエストを処理します
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
    if pageStr := query.Get("page"); pageStr != "" {
        if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
            page = p
        }
    }
    
    limit := 10
    if limitStr := query.Get("limit"); limitStr != "" {
        if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
            limit = l
        }
    }
    
    // 3. ビジネスロジック実行
    todos, err := h.todoService.GetAllTodos(r.Context())
    if err != nil {
        writeErrorResponse(w, http.StatusInternalServerError, "Failed to get todos", err.Error())
        return
    }
    
    // 4. ページング処理（簡易実装）
    total := len(todos)
    start := (page - 1) * limit
    end := start + limit
    
    if start >= total {
        todos = []*entity.Todo{} // 空のスライス
    } else {
        if end > total {
            end = total
        }
        todos = todos[start:end]
    }
    
    // 5. レスポンス作成
    response := dto.ToTodoListResponse(todos, page, limit, int64(total))
    writeJSONResponse(w, http.StatusOK, response)
}

// GetTodoByID は指定されたIDのTodoを取得するHTTPハンドラーです
// GET /api/v1/todos/{id} へのリクエストを処理します
func (h *TodoHandler) GetTodoByID(w http.ResponseWriter, r *http.Request) {
    // 1. HTTPメソッドの確認
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    // 2. URLパスからIDを抽出
    // 標準net/httpでは手動でパスパラメータを解析
    path := strings.TrimPrefix(r.URL.Path, "/api/v1/todos/")
    idStr := strings.Split(path, "/")[0] // パスの最初の部分がID
    
    // 3. IDの妥当性チェック
    id, err := strconv.Atoi(idStr)
    if err != nil || id <= 0 {
        writeErrorResponse(w, http.StatusBadRequest, "Invalid todo ID", "")
        return
    }
    
    // 4. ビジネスロジック実行
    todo, err := h.todoService.GetTodoByID(r.Context(), id)
    if err != nil {
        if strings.Contains(err.Error(), "not found") {
            writeErrorResponse(w, http.StatusNotFound, "Todo not found", err.Error())
        } else {
            writeErrorResponse(w, http.StatusInternalServerError, "Failed to get todo", err.Error())
        }
        return
    }
    
    // 5. レスポンス作成
    response := dto.ToTodoResponse(todo)
    writeJSONResponse(w, http.StatusOK, response)
}

// ヘルパー関数：JSONレスポンスを書き込む
func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    
    if err := json.NewEncoder(w).Encode(data); err != nil {
        // JSONエンコードに失敗した場合
        http.Error(w, "Failed to encode response", http.StatusInternalServerError)
    }
}

// ヘルパー関数：エラーレスポンスを書き込む  
func writeErrorResponse(w http.ResponseWriter, statusCode int, message, detail string) {
    errorResponse := dto.ErrorResponse{
        Error:   message,
        Message: detail,
        Code:    statusCode,
    }
    writeJSONResponse(w, statusCode, errorResponse)
}
```

#### 🎓 net/httpパッケージの学習ポイント

**🎯 標準パッケージを使う理由:**
- **軽量性**: 外部フレームワーク依存がゼロ
- **透明性**: HTTP処理の詳細を完全に制御
- **学習効果**: ウェブサーバーの仕組みを深く理解
- **パフォーマンス**: 余分なオーバーヘッドがない

**❌ よくあるHTTP処理の間違い:**

```go
// 🚫 危険なパターン：エラーハンドリングの不備
func badHandler(w http.ResponseWriter, r *http.Request) {
    var req CreateTodoRequest
    json.NewDecoder(r.Body).Decode(&req) // エラー無視
    
    // バリデーションなし
    todo := req.ToEntity()
    
    // 実行結果の確認なし
    createdTodo, _ := service.CreateTodo(todo)
    
    // Content-Type設定忘れ
    json.NewEncoder(w).Encode(createdTodo)
}

// 🚫 HTTPメソッドチェック忘れ
func badMethodHandler(w http.ResponseWriter, r *http.Request) {
    // どんなHTTPメソッドでも受け付けてしまう
    var req CreateTodoRequest
    json.NewDecoder(r.Body).Decode(&req)
    // ...
}

// 🚫 リクエストボディのサイズ制限なし
func badSizeHandler(w http.ResponseWriter, r *http.Request) {
    // 巨大なリクエストでDDoS攻撃される可能性
    var req CreateTodoRequest
    json.NewDecoder(r.Body).Decode(&req)
    // ...
}

// ✅ 正しいパターン：適切なエラーハンドリング
func goodHandler(w http.ResponseWriter, r *http.Request) {
    // HTTPメソッドチェック
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    // Content-Typeチェック
    if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
        writeErrorResponse(w, http.StatusBadRequest, "Invalid Content-Type", "")
        return
    }
    
    // リクエストサイズ制限
    r.Body = http.MaxBytesReader(w, r.Body, 1048576) // 1MB制限
    
    // JSONデコード with エラーハンドリング
    var req CreateTodoRequest
    decoder := json.NewDecoder(r.Body)
    decoder.DisallowUnknownFields()
    
    if err := decoder.Decode(&req); err != nil {
        writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON", err.Error())
        return
    }
    
    // バリデーション実行
    if err := req.Validate(); err != nil {
        writeErrorResponse(w, http.StatusBadRequest, "Validation failed", err.Error())
        return
    }
    
    // ビジネスロジック実行 with エラーハンドリング
    todo, err := h.service.CreateTodo(r.Context(), req.ToEntity())
    if err != nil {
        writeErrorResponse(w, http.StatusInternalServerError, "Creation failed", err.Error())
        return
    }
    
    // 適切なレスポンス
    writeJSONResponse(w, http.StatusCreated, dto.ToTodoResponse(todo))
}
```

**📚 技術用語解説：**

| 用語 | 意味 | 重要度 |
|------|------|-------|
| http.HandlerFunc | HTTPハンドラー関数の型定義 | ★★★ |
| http.ResponseWriter | HTTPレスポンスを書き込むインターフェース | ★★★ |
| Content-Type | リクエスト/レスポンスのデータ形式を示すヘッダー | ★★☆ |
| json.Decoder | JSONデータを構造体にデコードする標準機能 | ★★☆ |
| http.MaxBytesReader | リクエストサイズを制限するセキュリティ機能 | ★★☆ |

---

## Chapter 5: ミドルウェアとルーティング

### 5.1 ミドルウェアの実装

`internal/application/middleware/middleware.go`を作成：

```go
package middleware

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "runtime/debug"
    "time"
    
    "github.com/google/uuid"
)

// MiddlewareChain はミドルウェアチェーンを管理する構造体です
// 標準net/httpでのミドルウェアパターンを学習します
type MiddlewareChain struct {
    middlewares []Middleware
}

// Middleware はミドルウェア関数の型定義です
// 標準的なミドルウェアパターンの実装
type Middleware func(http.Handler) http.Handler

// NewMiddlewareChain は新しいミドルウェアチェーンを作成します
func NewMiddlewareChain(middlewares ...Middleware) *MiddlewareChain {
    return &MiddlewareChain{
        middlewares: middlewares,
    }
}

// Then はハンドラーにミドルウェアチェーンを適用します
// 逆順で適用されるため、最初に登録したミドルウェアが最外層になります
func (mc *MiddlewareChain) Then(handler http.Handler) http.Handler {
    // チェーンを逆順で適用
    for i := len(mc.middlewares) - 1; i >= 0; i-- {
        handler = mc.middlewares[i](handler)
    }
    return handler
}

// LoggingMiddleware はリクエスト/レスポンスのログを出力するミドルウェアです
// 標準logパッケージを使用したロギング実装を学習
func LoggingMiddleware() Middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            
            // レスポンスライターをラップして情報を取得
            wrapped := &responseWriter{
                ResponseWriter: w,
                statusCode:    http.StatusOK, // デフォルトは200
            }
            
            // リクエスト開始ログ
            log.Printf("[INFO] %s %s started", r.Method, r.RequestURI)
            
            // 次のハンドラーを実行
            next.ServeHTTP(wrapped, r)
            
            // レスポンス完了ログ
            duration := time.Since(start)
            log.Printf("[INFO] %s %s completed in %v - Status: %d",
                r.Method, r.RequestURI, duration, wrapped.statusCode)
        })
    }
}

// RequestIDMiddleware はリクエストにユニークなIDを付与するミドルウェアです
// トレーシングとデバッグのために使用します
func RequestIDMiddleware() Middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // ユニークなリクエストIDを生成
            requestID := uuid.New().String()
            
            // コンテキストにリクエストIDを追加
            ctx := context.WithValue(r.Context(), "request_id", requestID)
            
            // レスポンスヘッダーにリクエストIDを追加
            w.Header().Set("X-Request-ID", requestID)
            
            // 更新されたコンテキストでリクエストを続行
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

// RecoveryMiddleware はパニックから回復し、500エラーを返すミドルウェアです
// アプリケーションの安定性確保のための重要な機能
func RecoveryMiddleware() Middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            defer func() {
                if err := recover(); err != nil {
                    // スタックトレースを取得
                    stack := debug.Stack()
                    
                    // エラーログ出力
                    log.Printf("[ERROR] Panic recovered: %v\nStack trace:\n%s", err, stack)
                    
                    // クライアントには一般的なエラーメッセージを返す
                    http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                }
            }()
            
            // 次のハンドラーを実行
            next.ServeHTTP(w, r)
        })
    }
}

// CORSMiddleware はCross-Origin Resource Sharing設定を行うミドルウェアです
// フロントエンドとの連携に必要なセキュリティ設定
func CORSMiddleware() Middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // CORS ヘッダーを設定
            w.Header().Set("Access-Control-Allow-Origin", "*")
            w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
            w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
            w.Header().Set("Access-Control-Max-Age", "3600")
            
            // プリフライトリクエストの処理
            if r.Method == http.MethodOptions {
                w.WriteHeader(http.StatusNoContent)
                return
            }
            
            // 次のハンドラーを実行
            next.ServeHTTP(w, r)
        })
    }
}

// responseWriter はhttp.ResponseWriterをラップしてステータスコードを記録します
// ロギングミドルウェアで使用するヘルパー構造体
type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

// WriteHeader はステータスコードを記録してから元のWriteHeaderを呼び出します
func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}

// Write はレスポンスボディを書き込みます
// 既にヘッダーが書き込まれていない場合、200ステータスで書き込みます
func (rw *responseWriter) Write(b []byte) (int, error) {
    if rw.statusCode == 0 {
        rw.statusCode = http.StatusOK
    }
    return rw.ResponseWriter.Write(b)
}
```

### 5.2 ルーティングの実装

`internal/infrastructure/web/routes.go`を作成：

```go
package web

import (
    "net/http"
    "strings"
    
    "todoapp-api-golang/internal/application/handler"
    "todoapp-api-golang/internal/application/middleware"
)

// Router は標準net/httpを使用したルーターです
// フレームワークを使わずに手動でルーティングを実装する学習
type Router struct {
    todoHandler *handler.TodoHandler
    middleware  *middleware.MiddlewareChain
}

// NewRouter は新しいルーターを作成します
func NewRouter(todoHandler *handler.TodoHandler) *Router {
    // ミドルウェアチェーンを作成
    // 実行順序：Recovery → CORS → Logging → RequestID → Handler
    middlewareChain := middleware.NewMiddlewareChain(
        middleware.RecoveryMiddleware(),
        middleware.CORSMiddleware(), 
        middleware.LoggingMiddleware(),
        middleware.RequestIDMiddleware(),
    )
    
    return &Router{
        todoHandler: todoHandler,
        middleware:  middlewareChain,
    }
}

// ServeHTTP はhttp.Handlerインターフェースを実装します
// 標準net/httpでの手動ルーティング実装を学習
func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // パスを正規化（末尾スラッシュを削除）
    path := strings.TrimSuffix(r.URL.Path, "/")
    if path == "" {
        path = "/"
    }
    
    // ルーティングロジック
    switch {
    case path == "/health":
        rt.middleware.Then(http.HandlerFunc(rt.healthCheck)).ServeHTTP(w, r)
        
    case path == "/api/v1/todos":
        rt.handleTodosCollection(w, r)
        
    case strings.HasPrefix(path, "/api/v1/todos/") && len(strings.Split(path, "/")) == 5:
        rt.handleTodosItem(w, r)
        
    default:
        // 404 Not Found
        rt.middleware.Then(http.HandlerFunc(rt.notFound)).ServeHTTP(w, r)
    }
}

// handleTodosCollection は /api/v1/todos コレクション操作を処理します
func (rt *Router) handleTodosCollection(w http.ResponseWriter, r *http.Request) {
    var handler http.HandlerFunc
    
    switch r.Method {
    case http.MethodGet:
        handler = rt.todoHandler.GetAllTodos
    case http.MethodPost:
        handler = rt.todoHandler.CreateTodo
    default:
        handler = rt.methodNotAllowed
    }
    
    rt.middleware.Then(handler).ServeHTTP(w, r)
}

// handleTodosItem は /api/v1/todos/{id} 個別操作を処理します
func (rt *Router) handleTodosItem(w http.ResponseWriter, r *http.Request) {
    var handler http.HandlerFunc
    
    switch r.Method {
    case http.MethodGet:
        handler = rt.todoHandler.GetTodoByID
    case http.MethodPut:
        handler = rt.todoHandler.UpdateTodo
    case http.MethodDelete:
        handler = rt.todoHandler.DeleteTodo
    default:
        handler = rt.methodNotAllowed
    }
    
    rt.middleware.Then(handler).ServeHTTP(w, r)
}

// healthCheck はヘルスチェックエンドポイントです
func (rt *Router) healthCheck(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"status":"healthy","timestamp":"` + time.Now().Format(time.RFC3339) + `"}`))
}

// notFound は404エラーを返します
func (rt *Router) notFound(w http.ResponseWriter, r *http.Request) {
    http.Error(w, "Not Found", http.StatusNotFound)
}

// methodNotAllowed は405エラーを返します
func (rt *Router) methodNotAllowed(w http.ResponseWriter, r *http.Request) {
    http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}

// extractIDFromPath はURLパスからIDを抽出するヘルパー関数です
func extractIDFromPath(path string) (int, error) {
    parts := strings.Split(strings.Trim(path, "/"), "/")
    if len(parts) < 4 {
        return 0, fmt.Errorf("invalid path")
    }
    
    idStr := parts[3] // /api/v1/todos/{id} の{id}部分
    id, err := strconv.Atoi(idStr)
    if err != nil || id <= 0 {
        return 0, fmt.Errorf("invalid ID: %s", idStr)
    }
    
    return id, nil
}
    // 1. HTTPメソッドの確認
    if r.Method != http.MethodPut {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    // 2. URLパスからIDを抽出
    id, err := extractIDFromPath(r.URL.Path)
    if err != nil {
        writeErrorResponse(w, http.StatusBadRequest, "Invalid todo ID", err.Error())
        return
    }
    
    // 3. Content-Type の確認
    contentType := r.Header.Get("Content-Type")
    if !strings.Contains(contentType, "application/json") {
        writeErrorResponse(w, http.StatusBadRequest, "Content-Type must be application/json", "")
        return
    }
    
    // 4. リクエストボディの解析
    var req dto.UpdateTodoRequest
    decoder := json.NewDecoder(r.Body)
    decoder.DisallowUnknownFields()
    
    if err := decoder.Decode(&req); err != nil {
        writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON format", err.Error())
        return
    }
    
    // 5. 既存Todoの取得
    existingTodo, err := h.todoService.GetTodoByID(r.Context(), id)
    if err != nil {
        writeErrorResponse(w, http.StatusInternalServerError, "Failed to get todo", err.Error())
        return
    }
    
    if existingTodo == nil {
        writeErrorResponse(w, http.StatusNotFound, "Todo not found", "")
        return
    }
    
    // 6. 更新内容を既存Todoに適用
    req.ApplyToEntity(existingTodo)
    
    // 7. 更新実行
    updatedTodo, err := h.todoService.UpdateTodo(r.Context(), existingTodo)
    if err != nil {
        if strings.Contains(err.Error(), "invalid todo") {
            writeErrorResponse(w, http.StatusBadRequest, "Validation error", err.Error())
        } else {
            writeErrorResponse(w, http.StatusInternalServerError, "Failed to update todo", err.Error())
        }
        return
    }
    
    // 8. EntityからDTOへの変換とレスポンス
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
    id, err := extractIDFromPath(r.URL.Path)
    if err != nil {
        writeErrorResponse(w, http.StatusBadRequest, "Invalid todo ID", err.Error())
        return
    }
    
    // 3. 削除実行
    err = h.todoService.DeleteTodo(r.Context(), id)
    if err != nil {
        if strings.Contains(err.Error(), "invalid id") {
            writeErrorResponse(w, http.StatusBadRequest, "Invalid ID", err.Error())
        } else {
            writeErrorResponse(w, http.StatusInternalServerError, "Failed to delete todo", err.Error())
        }
        return
    }
    
    // 4. 削除成功レスポンス（204 No Content）
    w.WriteHeader(http.StatusNoContent)
}

// CompleteTodo はTodoの完了状態を切り替えるHTTPハンドラーです
// POST /api/v1/todos/{id}/complete へのリクエストを処理します
func (h *TodoHandler) CompleteTodo(w http.ResponseWriter, r *http.Request) {
    // 1. HTTPメソッドの確認
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    // 2. URLパスからIDを抽出
    id, err := extractIDFromPath(r.URL.Path)
    if err != nil {
        writeErrorResponse(w, http.StatusBadRequest, "Invalid todo ID", err.Error())
        return
    }
    
    // 3. 完了処理実行
    completedTodo, err := h.todoService.CompleteTodo(r.Context(), id)
    if err != nil {
        if strings.Contains(err.Error(), "invalid id") || strings.Contains(err.Error(), "not found") {
            writeErrorResponse(w, http.StatusNotFound, "Todo not found", err.Error())
        } else {
            writeErrorResponse(w, http.StatusInternalServerError, "Failed to complete todo", err.Error())
        }
        return
    }
    
    // 4. EntityからDTOへの変換とレスポンス
    response := dto.ToTodoResponse(completedTodo)
    writeJSONResponse(w, http.StatusOK, response)
}

// --- ヘルパー関数 ---

// extractIDFromPath はURLパスからIDを抽出します
// 例: "/api/v1/todos/123" から "123" を抽出
func extractIDFromPath(path string) (int, error) {
    // パスを "/" で分割
    segments := strings.Split(strings.Trim(path, "/"), "/")
    
    // 最低限必要なセグメント数をチェック
    if len(segments) < 4 { // ["api", "v1", "todos", "id"]
        return 0, fmt.Errorf("invalid path format")
    }
    
    // 最後のセグメントがIDの場合
    if segments[len(segments)-1] != "complete" {
        return strconv.Atoi(segments[len(segments)-1])
    }
    
    // "complete" の前のセグメントがIDの場合
    if len(segments) >= 5 {
        return strconv.Atoi(segments[len(segments)-2])
    }
    
    return 0, fmt.Errorf("invalid path format")
}

// writeJSONResponse はJSONレスポンスを書き込みます
func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
    // Content-Type ヘッダーの設定
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    
    // JSON エンコード
    if err := json.NewEncoder(w).Encode(data); err != nil {
        // JSONエンコードエラーの場合は500エラー
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    }
}

// writeErrorResponse はエラーレスポンスを書き込みます
func writeErrorResponse(w http.ResponseWriter, statusCode int, message, details string) {
    errorResponse := dto.ErrorResponse{
        Error:   message,
        Details: details,
    }
    
    writeJSONResponse(w, statusCode, errorResponse)
}
```

**学習ポイント:**
- 標準`net/http`パッケージでのHTTPハンドラー実装
- JSON処理（`encoding/json`）
- エラーハンドリングとHTTPステータスコード
- URLパラメータの手動解析

---

## Chapter 5: ミドルウェアの実装

`internal/application/middleware/middleware.go`を作成：

```go
package middleware

import (
    "fmt"
    "log"
    "net/http"
    "runtime/debug"
    "strconv"
    "time"
)

// Middleware は http.Handler を受け取り、http.Handler を返す関数型です
// ミドルウェアパターンの標準的な実装
type Middleware func(http.Handler) http.Handler

// ChainMiddleware は複数のミドルウェアを chain します
// 標準パッケージでのミドルウェアチェーン実装の学習
func ChainMiddleware(middlewares ...Middleware) Middleware {
    return func(next http.Handler) http.Handler {
        // 逆順で適用することで、指定した順序で実行される
        for i := len(middlewares) - 1; i >= 0; i-- {
            next = middlewares[i](next)
        }
        return next
    }
}

// LoggingMiddleware はHTTPリクエストをログ出力するミドルウェアです
// 標準パッケージでのHTTPログ実装
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // ResponseWriter をラップして情報を取得
        recorder := NewResponseRecorder(w)
        
        // 次のハンドラーを実行
        next.ServeHTTP(recorder, r)
        
        // ログ出力
        duration := time.Since(start)
        log.Printf("%s %s %s %d %d %v",
            r.RemoteAddr,
            r.Method,
            r.URL.Path,
            recorder.statusCode,
            recorder.responseSize,
            duration,
        )
    })
}

// DetailedLoggingMiddleware は詳細なログを出力するミドルウェアです
func DetailedLoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // リクエスト詳細をログ出力
        log.Printf("→ %s %s %s", r.Method, r.URL.Path, r.Proto)
        for key, values := range r.Header {
            for _, value := range values {
                log.Printf("  %s: %s", key, value)
            }
        }
        
        recorder := NewResponseRecorder(w)
        next.ServeHTTP(recorder, r)
        
        // レスポンス詳細をログ出力
        duration := time.Since(start)
        log.Printf("← %s %s %d %d %v",
            r.Method,
            r.URL.Path,
            recorder.statusCode,
            recorder.responseSize,
            duration,
        )
        
        for key, values := range recorder.Header() {
            for _, value := range values {
                log.Printf("  %s: %s", key, value)
            }
        }
    })
}

// RequestIDMiddleware はリクエストにユニークなIDを付与するミドルウェアです
func RequestIDMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 既存のリクエストIDをチェック
        requestID := r.Header.Get("X-Request-ID")
        
        // リクエストIDが無い場合は生成
        if requestID == "" {
            requestID = generateRequestID()
        }
        
        // レスポンスヘッダーにリクエストIDを設定
        w.Header().Set("X-Request-ID", requestID)
        
        // ログ出力
        log.Printf("Request ID: %s - %s %s", requestID, r.Method, r.URL.Path)
        
        next.ServeHTTP(w, r)
    })
}

// RecoveryMiddleware はパニックを回復するミドルウェアです
// アプリケーションのクラッシュを防止
func RecoveryMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                // パニック情報をログ出力
                log.Printf("PANIC: %v", err)
                log.Printf("Request: %s %s", r.Method, r.URL.Path)
                log.Printf("Stack trace:\n%s", debug.Stack())
                
                // クライアントには500エラーを返す
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            }
        }()
        
        next.ServeHTTP(w, r)
    })
}

// ResponseRecorder はhttp.ResponseWriterをラップしてレスポンス情報を記録します
type ResponseRecorder struct {
    http.ResponseWriter
    statusCode   int
    responseSize int
}

// NewResponseRecorder はResponseRecorderのコンストラクタです
func NewResponseRecorder(w http.ResponseWriter) *ResponseRecorder {
    return &ResponseRecorder{
        ResponseWriter: w,
        statusCode:     http.StatusOK, // デフォルトは200
        responseSize:   0,
    }
}

// WriteHeader はステータスコードを記録します
func (r *ResponseRecorder) WriteHeader(statusCode int) {
    r.statusCode = statusCode
    r.ResponseWriter.WriteHeader(statusCode)
}

// Write はレスポンスボディを書き込み、サイズを記録します
func (r *ResponseRecorder) Write(data []byte) (int, error) {
    size, err := r.ResponseWriter.Write(data)
    r.responseSize += size
    return size, err
}

// generateRequestID はユニークなリクエストIDを生成します
func generateRequestID() string {
    // 簡単な実装：現在時刻のナノ秒を使用
    // 実際のプロダクションではUUIDなどを使用することが推奨
    return fmt.Sprintf("req_%d", time.Now().UnixNano())
}
```

**学習ポイント:**
- ミドルウェアパターンの実装
- `http.ResponseWriter`のラッピング
- パニック回復とログ出力
- チェーンパターンによる組み合わせ

---

## Chapter 6: テストの実装

包括的なテストスイートは既に実装済みですが、主要なテストパターンを理解しましょう。

### 6.1 エンティティテスト（例）

```go
// internal/domain/entity/todo_test.go から抜粋
func TestTodo_IsValid(t *testing.T) {
    tests := []struct {
        name   string
        todo   Todo
        expect bool
    }{
        {
            name: "有効なTodo",
            todo: Todo{
                Title:       "有効なタイトル",
                Description: "有効な説明文",
                IsCompleted: false,
            },
            expect: true,
        },
        // ... 他のテストケース
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := tt.todo.IsValid(); got != tt.expect {
                t.Errorf("IsValid() = %v, want %v", got, tt.expect)
            }
        })
    }
}
```

### 6.2 モックの実装例

```go
// サービステスト用のモック実装例
type MockTodoRepository struct {
    todos       map[int]*entity.Todo
    nextID      int
    shouldError bool
    errorMsg    string
}

func (m *MockTodoRepository) Create(ctx context.Context, todo *entity.Todo) (*entity.Todo, error) {
    if m.shouldError {
        return nil, errors.New(m.errorMsg)
    }
    
    m.nextID++
    todo.ID = m.nextID
    todo.CreatedAt = time.Now()
    todo.UpdatedAt = time.Now()
    
    todoToSave := *todo
    m.todos[todo.ID] = &todoToSave
    
    return todo, nil
}
```

---

## Chapter 7: サーバーの起動と統合

### 7.1 サーバー設定

`internal/infrastructure/web/server.go`を作成：

```go
package web

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "time"
)

// Server はHTTPサーバーを管理する構造体です
type Server struct {
    httpServer *http.Server
    addr       string
}

// NewServer はサーバーのコンストラクタです
func NewServer(addr string, handler http.Handler) *Server {
    return &Server{
        httpServer: &http.Server{
            Addr:         addr,
            Handler:      handler,
            ReadTimeout:  15 * time.Second,
            WriteTimeout: 15 * time.Second,
            IdleTimeout:  60 * time.Second,
        },
        addr: addr,
    }
}

// Start はサーバーを開始します
func (s *Server) Start() error {
    log.Printf("Starting server on %s", s.addr)
    return s.httpServer.ListenAndServe()
}

// Shutdown はサーバーをシャットダウンします
func (s *Server) Shutdown(ctx context.Context) error {
    log.Println("Shutting down server...")
    return s.httpServer.Shutdown(ctx)
}
```

### 7.2 ルーティング設定

`internal/infrastructure/web/routes.go`を作成：

```go
package web

import (
    "net/http"
    
    "todoapp-api-golang/internal/application/handler"
    "todoapp-api-golang/internal/application/middleware"
)

// SetupRoutes はルーティングを設定します
func SetupRoutes(todoHandler *handler.TodoHandler) http.Handler {
    mux := http.NewServeMux()
    
    // ヘルスチェック
    mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    })
    
    // Todo API routes
    mux.HandleFunc("/api/v1/todos", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            todoHandler.GetAllTodos(w, r)
        case http.MethodPost:
            todoHandler.CreateTodo(w, r)
        default:
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    })
    
    // 個別Todo操作
    mux.HandleFunc("/api/v1/todos/", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            todoHandler.GetTodoByID(w, r)
        case http.MethodPut:
            todoHandler.UpdateTodo(w, r)
        case http.MethodDelete:
            todoHandler.DeleteTodo(w, r)
        case http.MethodPost:
            if strings.HasSuffix(r.URL.Path, "/complete") {
                todoHandler.CompleteTodo(w, r)
            } else {
                http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            }
        default:
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    })
    
    // ミドルウェアチェーンの適用
    handler := middleware.ChainMiddleware(
        middleware.RecoveryMiddleware,
        middleware.RequestIDMiddleware,
        middleware.LoggingMiddleware,
    )(mux)
    
    return handler
}
```

### 7.3 メインアプリケーション

`cmd/api/main.go`を作成：

```go
package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"
    
    "todoapp-api-golang/internal/application/handler"
    "todoapp-api-golang/internal/domain/service"
    "todoapp-api-golang/internal/infrastructure/database"
    "todoapp-api-golang/internal/infrastructure/web"
    "todoapp-api-golang/pkg/config"
)

func main() {
    // 1. 設定の読み込み
    cfg := config.Load()
    
    // 2. データベース接続
    dbManager := database.NewDatabaseManager()
    if err := dbManager.ConnectSQLite(cfg.DatabaseURL); err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer dbManager.Close()
    
    // 3. テーブル作成
    if err := dbManager.CreateTables(); err != nil {
        log.Fatalf("Failed to create tables: %v", err)
    }
    
    // 4. 依存性注入
    todoRepo := database.NewTodoRepository(dbManager.GetDB())
    todoService := service.NewTodoService(todoRepo)
    todoHandler := handler.NewTodoHandler(todoService)
    
    // 5. ルーティング設定
    router := web.SetupRoutes(todoHandler)
    
    // 6. サーバー起動
    server := web.NewServer(cfg.ServerAddress, router)
    
    // 7. Graceful shutdown の設定
    go func() {
        if err := server.Start(); err != nil {
            log.Printf("Server error: %v", err)
        }
    }()
    
    // 8. シグナル待ち
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    // 9. シャットダウン
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    if err := server.Shutdown(ctx); err != nil {
        log.Printf("Server shutdown error: %v", err)
    } else {
        log.Println("Server shutdown complete")
    }
}
```

---

## 🎯 学習の進め方

### 1. 段階的な実装
1. **Chapter 2から順番に実装**
2. **各段階でテストを実行**
3. **コンパイルエラーを一つずつ解決**

### 2. 理解度チェック
- 各章の「学習ポイント」を理解できているか確認
- コードの意味を説明できるか
- なぜその設計にしたかを理解しているか

### 3. 実験と改善
- コードを変更して動作を確認
- エラーハンドリングを追加
- 新しい機能を実装

### 4. テスト駆動開発の実践
- テストを先に書く
- 実装後にテストを実行
- リファクタリングでコードを改善

---

## 🔧 動作確認方法

### 1. プロジェクトの起動
```bash
# ホットリロード環境での開発
air

# または直接実行
go run cmd/api/main.go
```

### 2. API動作確認
```bash
# Todo作成
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Content-Type: application/json" \
  -d '{"title":"学習用Todo","description":"Go APIの学習"}'

# Todo一覧取得
curl http://localhost:8080/api/v1/todos

# 特定Todo取得
curl http://localhost:8080/api/v1/todos/1
```

### 3. テスト実行
```bash
# 全テスト実行
go test ./...

# 詳細表示
go test ./... -v

# カバレッジ
go test ./... -cover
```

---

## 🚀 次のステップと発展課題

### レベル1: 基本機能の拡張
1. **バリデーション強化**
   - 文字数制限（タイトル100文字、説明500文字）
   - 必須フィールドのチェック
   - 特殊文字のサニタイゼーション

2. **エラーハンドリング改善**
   - カスタムエラータイプの実装
   - エラーコードの統一化
   - ログレベルの適切な設定

### レベル2: 高度な機能実装
1. **認証・認可**
   - JWT認証の実装
   - ミドルウェアでの認証チェック
   - ユーザー管理機能

2. **パフォーマンス最適化**
   - データベース接続プールの調整
   - キャッシュ機能の追加
   - リクエストレート制限

### レベル3: 本格的な運用機能
1. **監視・メトリクス**
   - ヘルスチェック機能の拡張
   - メトリクス収集とログ出力
   - 構造化ログ（JSON形式）の実装

2. **デプロイメント準備**
   - Docker化
   - 環境変数による設定管理
   - CI/CDパイプラインの構築

### 学習リソース
- [Go公式ドキュメント](https://golang.org/doc/)
- [Clean Architecture書籍](https://www.amazon.com/Clean-Architecture-Craftsmans-Software-Structure/dp/0134494164)
- [Go標準パッケージドキュメント](https://pkg.go.dev/std)

---

## 💡 追加の学習ポイント

### セキュリティ
```go
// 良い例：SQLインジェクション対策
query := "SELECT * FROM todos WHERE id = ?"
row := db.QueryRow(query, todoID)

// 悪い例：脆弱性あり
query := fmt.Sprintf("SELECT * FROM todos WHERE id = %s", todoID)
```

### エラーハンドリングのベストプラクティス
```go
// 良い例：適切なエラーラッピング
if err != nil {
    return nil, fmt.Errorf("failed to create todo: %w", err)
}

// 悪い例：エラー情報の損失
if err != nil {
    return nil, errors.New("something went wrong")
}
```

### リソース管理
```go
// 良い例：確実なリソース解放
rows, err := db.Query(query)
if err != nil {
    return err
}
defer rows.Close() // 必須

// 悪い例：リソースリーク
rows, _ := db.Query(query)
// Close()を忘れると接続が残り続ける
```

---

このチュートリアルを通じて、Go標準パッケージによるバックエンドAPI開発の基礎から応用まで、実践的に学習できます。写経によって手を動かしながら、Clean ArchitectureとGoのベストプラクティスを身につけてください。

実際の開発では、このプロジェクトをベースに機能を拡張し、より複雑なビジネスロジックやパフォーマンス要件に対応していくことができます。標準パッケージでの実装を通じて、Goの本質的な理解を深めることが、長期的な開発スキル向上に繋がります。