# Todo API アーキテクチャ設計

## 概要

このプロジェクトは、Clean ArchitectureとDomain-Driven Design（DDD）の原則に基づいて設計されています。学習目的として、保守性、テスト容易性、および拡張性を重視したアーキテクチャパターンを採用しています。

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
│   └── user.go                 # Userエンティティ
├── repository/
│   ├── todo_repository.go      # Todo操作インターフェース
│   └── user_repository.go      # User操作インターフェース
└── service/
    └── todo_service.go         # Todoビジネスロジック
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
│   └── health_handler.go       # ヘルスチェックハンドラー
├── middleware/
│   ├── cors.go                 # CORS設定
│   ├── logger.go               # ログ出力
│   └── auth.go                 # 認証処理
└── dto/
    ├── todo_request.go         # Todoリクエスト用DTO
    └── todo_response.go        # Todoレスポンス用DTO
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
│   └── user_repository_impl.go # Userリポジトリ実装
└── web/
    ├── server.go               # HTTPサーバー
    └── routes.go               # ルーティング設定
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
- **Unit Test**: 各層を独立してテスト
- **Integration Test**: DB込みのテスト
- **E2E Test**: API全体のテスト

## 実装予定のAPI

### Todo管理
- `GET /todos` - Todo一覧取得
- `GET /todos/{id}` - 特定Todo取得
- `POST /todos` - Todo作成
- `PUT /todos/{id}` - Todo更新
- `DELETE /todos/{id}` - Todo削除

### システム
- `GET /health` - ヘルスチェック

## 技術スタック（予定）

- **Language**: Go 1.21+
- **Web Framework**: Gin または Echo
- **Database**: PostgreSQL または MySQL
- **ORM**: GORM
- **Migration**: golang-migrate
- **Testing**: testify
- **Documentation**: Swagger/OpenAPI

## 今後の拡張予定

1. **認証・認可機能**
2. **ログ出力の強化**
3. **メトリクス収集**
4. **Docker化**
5. **CI/CD パイプライン**

このアーキテクチャにより、学習者はモダンなGo APIアプリケーションの設計パターンと実装方法を段階的に学習できます。