# Todo API アーキテクチャ設計

## 概要

このプロジェクトは、Clean ArchitectureとDomain-Driven Design（DDD）の原則に基づいて設計されています。**Go標準パッケージのみ**を使用した学習目的のプロジェクトとして、保守性、テスト容易性、および拡張性を重視したアーキテクチャパターンを採用しています。

## アーキテクチャパターン

### Clean Architecture
- **依存関係の方向**: 外側の層は内側の層に依存するが、内側の層は外側の層に依存しない
- **ビジネスロジックの独立性**: ドメインレイヤーは外部の詳細（DB、フレームワークなど）から独立
- **テスタビリティ**: 各層を独立してテストできる構造

### レイヤー構成

#### 1. Domain Layer（ドメイン層）
- **Entity**: ビジネスの中核となるデータ構造と振る舞い
- **Repository**: データアクセスの抽象化インターフェース
- **Service**: ドメインロジックとビジネスルール

#### 2. Application Layer（アプリケーション層）
- **Handler**: HTTPリクエストの処理とレスポンス
- **DTO (Data Transfer Object)**: レイヤー間でのデータ転送
- **Middleware**: 横断的関心事（認証、ログ、CORS等）

#### 3. Infrastructure Layer（インフラストラクチャ層）
- **Database**: データベース実装とマイグレーション
- **Web**: HTTPルーティングとサーバー設定

## ディレクトリ構造

```
todoapp-api-golang/
├── cmd/
│   └── api/                    # アプリケーションエントリーポイント
├── internal/
│   ├── domain/
│   │   ├── entity/             # ドメインエンティティ（Todo, User等）
│   │   ├── repository/         # リポジトリインターフェース
│   │   └── service/            # ドメインサービス
│   ├── infrastructure/
│   │   ├── database/           # DB接続、実装
│   │   └── web/                # HTTPサーバー、ルーティング
│   └── application/
│       ├── handler/            # HTTPハンドラー
│       ├── middleware/         # ミドルウェア
│       └── dto/                # データ転送オブジェクト
├── pkg/
│   ├── config/                 # 設定管理
│   └── utils/                  # ユーティリティ関数
├── docs/                       # ドキュメント
├── migrations/                 # データベースマイグレーション
└── CLAUDE.md                   # Claude Code向けガイド
```

## 各層の役割と責任

### Domain Layer
```
internal/domain/
├── entity/
│   ├── todo.go                 # Todoエンティティ
│   └── todo_test.go            # エンティティテスト
├── repository/
│   └── todo_repository.go      # Todo操作インターフェース
└── service/
    ├── todo_service.go         # Todoビジネスロジック
    ├── todo_service_interface.go # サービスインターフェース（テスト用）
    └── todo_service_test.go    # サービステスト
```

**責任**:
- ビジネスルールの実装
- データ構造と振る舞いの定義
- ドメイン固有のバリデーション

### Application Layer
```
internal/application/
├── handler/
│   ├── todo_handler.go         # Todo API ハンドラー
│   └── todo_handler_test.go    # ハンドラーテスト
├── middleware/
│   ├── middleware.go           # ミドルウェア実装（ログ、リクエストID、パニック回復等）
│   └── middleware_test.go      # ミドルウェアテスト
└── dto/
    ├── todo_request.go         # Todoリクエスト用DTO
    ├── todo_response.go        # Todoレスポンス用DTO
    └── todo_dto_test.go        # DTOテスト
```

**責任**:
- HTTPリクエスト/レスポンスの処理
- 入力値のバリデーション
- ドメインサービスのオーケストレーション

### Infrastructure Layer
```
internal/infrastructure/
├── database/
│   ├── connection.go           # DB接続管理
│   ├── todo_repository_impl.go # Todoリポジトリ実装
│   └── todo_repository_impl_test.go # リポジトリ統合テスト
└── web/
    ├── server.go               # HTTPサーバー（標準net/httpパッケージ使用）
    └── routes.go               # ルーティング設定（手動ルーティング）
```

**責任**:
- 外部システムとの統合
- データベースアクセスの実装
- HTTPサーバーの設定

## データフロー

1. **HTTPリクエスト** → Handler
2. **Handler** → DTO変換 → Domain Service
3. **Domain Service** → Repository Interface → Database Implementation
4. **Database** → Entity → Domain Service
5. **Domain Service** → Handler → DTO変換
6. **Handler** → HTTPレスポンス

## 学習ポイント

### 1. 依存性注入（Dependency Injection）
- インターフェースを通じた疎結合
- テスト時のモック化が容易

### 2. 責任の分離（Separation of Concerns）
- 各層が明確な責任を持つ
- 変更が他の層に影響しにくい

### 3. テスト戦略
本プロジェクトでは包括的なテストスイートを実装済み：

- **Unit Test**: 各層を独立してテスト（モック使用）
  - Entity層: ビジネスロジックとバリデーション
  - Service層: ドメインサービスの振る舞い
  - Handler層: HTTPリクエスト/レスポンス処理

- **Integration Test**: 実際のデータベースを使用
  - Repository層: SQLite in-memoryでのCRUD操作テスト
  - トランザクション処理の検証

- **HTTP Test**: 標準`httptest`パッケージを使用
  - ミドルウェアチェーンの動作確認
  - JSON変換の検証
  - エラーハンドリングの確認

## 実装済みAPI

### Todo管理（実装済み）
- `GET /api/v1/todos` - Todo一覧取得（ページング対応）
- `GET /api/v1/todos/{id}` - 特定Todo取得
- `POST /api/v1/todos` - Todo作成
- `PUT /api/v1/todos/{id}` - Todo更新
- `DELETE /api/v1/todos/{id}` - Todo削除
- `POST /api/v1/todos/{id}/complete` - Todo完了状態の切り替え

### システム（実装済み）
- `GET /health` - ヘルスチェック

## 技術スタック（実装済み）

### 言語・フレームワーク
- **Language**: Go 1.21+
- **HTTP Server**: 標準`net/http`パッケージ
- **JSON処理**: 標準`encoding/json`パッケージ
- **ルーティング**: 手動実装（標準パッケージのみ）

### データベース
- **Database**: SQLite（開発・テスト）/ MySQL（本番想定）
- **SQL実行**: 標準`database/sql`パッケージ
- **Database Driver**: 
  - `github.com/mattn/go-sqlite3`（テスト用）
  - `github.com/go-sql-driver/mysql`（本番用）

### テスト
- **Testing Framework**: 標準`testing`パッケージ
- **HTTP Testing**: 標準`net/http/httptest`パッケージ
- **Mock**: 手動実装（外部ライブラリ不使用）
- **Database Testing**: SQLite in-memory

### 開発ツール
- **Hot Reload**: Air
- **Build**: 標準goコマンド

## プロジェクトの特徴

### 学習重視のアプローチ
1. **標準パッケージのみ使用**
   - 外部フレームワーク依存を排除
   - Go言語の基本機能への深い理解
   - プロダクションレディなコード品質

2. **包括的なテストスイート**
   - 全レイヤーをカバーするテスト
   - モック実装による単体テスト
   - 統合テストによる実際のDB操作検証
   - HTTPテストによるエンドツーエンド検証

3. **Clean Architectureの実践**
   - 明確な責任分離
   - 依存性逆転の原則
   - テスタブルな設計

### 実装済み機能
- ✅ 完全なCRUD API
- ✅ ミドルウェア実装（ログ、リクエストID、パニック回復）
- ✅ 包括的テストカバレッジ
- ✅ データベース統合（SQLite/MySQL対応）
- ✅ HTTPルーティングとハンドリング

## 今後の拡張予定

1. **認証・認可機能**
   - JWT認証の実装
   - ロールベースアクセス制御

2. **API ドキュメント**
   - OpenAPI/Swagger仕様書
   - 自動生成ドキュメント

3. **運用面の強化**
   - メトリクス収集
   - ヘルスチェック拡張
   - ログ出力の構造化

4. **デプロイメント**
   - Docker化
   - CI/CD パイプライン

このアーキテクチャにより、学習者はGoの標準パッケージを使ったモダンなAPI設計と実装方法を実践的に学習できます。外部フレームワークに依存しないため、Go言語の本質的な理解が深まります。