# Todo API - Go言語学習用プロジェクト

このプロジェクトは、Go言語でREST APIを開発するための教材として設計されています。Clean ArchitectureとDomain-Driven Designの原則に基づいて実装されており、モダンなAPI開発のベストプラクティスを学習できます。

## 🎯 学習目標

- Go言語でのWebAPI開発
- Clean Architectureの実装
- RESTful API設計
- データベース操作（database/sql）
- 依存性注入（DI）
- エラーハンドリング
- テスト駆動開発（TDD）

## 🏗️ アーキテクチャ

```
cmd/
├── api/              # アプリケーションエントリーポイント
internal/
├── domain/           # ドメイン層
│   ├── entity/       # エンティティ
│   ├── repository/   # リポジトリインターフェース  
│   └── service/      # ドメインサービス
├── application/      # アプリケーション層
│   ├── dto/          # データ転送オブジェクト
│   ├── handler/      # HTTPハンドラー
│   └── middleware/   # ミドルウェア
└── infrastructure/   # インフラストラクチャ層
    ├── database/     # データベース実装
    └── web/          # Webサーバー設定
pkg/
├── config/           # 設定管理
└── utils/            # ユーティリティ
```

詳細なアーキテクチャについては [docs/architecture.md](docs/architecture.md) を参照してください。

## 🚀 クイックスタート

### 前提条件

#### Docker使用時（推奨）
- Docker 20.10以上
- Docker Compose 2.0以上
- Git

#### ローカル環境使用時
- Go 1.21以上
- MySQL 8.0以上 または PostgreSQL 13以上 または SQLite3
- Git

### セットアップ

#### Option 1: Docker を使用（推奨）

1. **リポジトリのクローン**
```bash
git clone https://github.com/your-username/todoapp-api-golang.git
cd todoapp-api-golang
```

2. **Docker環境でのセットアップと起動**
```bash
# 初回セットアップ
make docker-setup

# アプリケーション起動
make docker-start
```

サーバーが起動します：
- API: `http://localhost:8080`
- phpMyAdmin: `http://localhost:8081`

#### Option 2: ローカル環境での実行

1. **リポジトリのクローン**
```bash
git clone https://github.com/your-username/todoapp-api-golang.git
cd todoapp-api-golang
```

2. **依存関係のインストール**
```bash
go mod tidy
```

3. **環境設定**
```bash
cp .env.example .env
# .env ファイルを編集してデータベース設定を行う
```

4. **データベースの準備**

**MySQL の場合:**
```sql
CREATE DATABASE todoapp CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

**PostgreSQL の場合:**
```sql
CREATE DATABASE todoapp;
```

**SQLite の場合:**
```bash
# 特別な準備は不要（ファイルが自動作成されます）
```

5. **アプリケーション実行**
```bash
go run cmd/api/main.go
```

サーバーが `http://localhost:8080` で起動します。

### 動作確認

```bash
# ヘルスチェック
curl http://localhost:8080/health

# Todo作成
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Content-Type: application/json" \
  -d '{"title":"サンプルタスク","description":"APIテスト用のタスクです"}'

# Todo一覧取得
curl http://localhost:8080/api/v1/todos
```

## 📋 API仕様

### エンドポイント一覧

| メソッド | エンドポイント | 説明 |
|---------|---------------|------|
| GET | `/health` | ヘルスチェック |
| GET | `/api/v1/todos` | Todo一覧取得 |
| POST | `/api/v1/todos` | Todo作成 |
| GET | `/api/v1/todos/:id` | Todo詳細取得 |
| PUT | `/api/v1/todos/:id` | Todo更新 |
| DELETE | `/api/v1/todos/:id` | Todo削除 |
| PATCH | `/api/v1/todos/:id/complete` | Todo完了 |
| PATCH | `/api/v1/todos/:id/incomplete` | Todo未完了 |

### リクエスト・レスポンス例

**Todo作成**
```bash
POST /api/v1/todos
Content-Type: application/json

{
  "title": "買い物リスト作成",
  "description": "明日の夕食の材料をリストアップする"
}
```

**レスポンス**
```json
{
  "id": 1,
  "title": "買い物リスト作成",
  "description": "明日の夕食の材料をリストアップする",
  "is_completed": false,
  "created_at": "2023-01-01T10:00:00Z",
  "updated_at": "2023-01-01T10:00:00Z"
}
```

## 🐳 Docker使用方法

### 基本コマンド

```bash
# 初回セットアップ
make docker-setup

# アプリケーション起動
make docker-start

# ログ確認
make docker-logs

# アプリケーション停止
make docker-stop

# 再起動
make docker-restart

# データとイメージのクリーンアップ
make docker-clean

# 完全リセット（データ削除）
make docker-reset
```

### Docker Compose直接操作

```bash
# ビルドして起動
docker-compose --env-file .env.docker up -d --build

# ログ表示
docker-compose --env-file .env.docker logs -f

# 停止
docker-compose --env-file .env.docker down

# データも含めて完全削除
docker-compose --env-file .env.docker down -v
```

### アクセス方法

- **Todo API**: http://localhost:8080
- **phpMyAdmin**: http://localhost:8081
  - サーバー: mysql
  - ユーザー名: root
  - パスワード: rootpassword

## 🔧 開発

### Makefileコマンド

```bash
# ヘルプ表示
make help

# 開発環境セットアップ
make setup

# アプリケーション実行
make run

# テスト実行
make test

# ビルド
make build

# フォーマット + 検査
make lint
```

### テスト実行

```bash
# 全テスト実行
make test
# または
go test ./...

# カバレッジ付きテスト
make test-coverage
# または
go test -cover ./...

# 特定パッケージのテスト
go test ./internal/domain/service/
```

### コードフォーマット

```bash
# フォーマット + 検査
make lint

# または個別に
make fmt    # フォーマット
make vet    # 検査
```

### ビルド

```bash
# 開発用ビルド
make build
# または
go build -o todoapp cmd/api/main.go

# 本番用ビルド（最適化）
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o todoapp cmd/api/main.go
```

## 🛠️ 設定

### 環境変数

| 変数名 | 説明 | デフォルト値 |
|-------|------|------------|
| `APP_ENV` | 実行環境 | `development` |
| `SERVER_PORT` | サーバーポート | `8080` |
| `DB_DRIVER` | DBドライバー | `mysql` |
| `DB_HOST` | DBホスト | `localhost` |
| `DB_PORT` | DBポート | `3306` |
| `DB_NAME` | DB名 | `todoapp` |
| `DB_USER` | DBユーザー | `root` |
| `DB_PASSWORD` | DBパスワード | 空文字 |

詳細は `.env.example` を参照してください。

## 📚 学習ガイド

### 段階的な学習プロセス

1. **基本理解** - エンティティとドメインモデル
2. **データアクセス** - リポジトリパターン
3. **ビジネスロジック** - ドメインサービス
4. **HTTP処理** - ハンドラーとDTO
5. **インフラ** - データベース接続とサーバー設定
6. **テスト** - 単体テストと統合テスト

### コードの読み方

各ファイルには学習者向けの詳細なコメントが記載されています：

- **なぜそうするのか** - 設計判断の理由
- **どうやって動くのか** - 実装の仕組み
- **何に注意するのか** - ベストプラクティスと落とし穴

## 🤝 コントリビューション

このプロジェクトは学習用教材のため、以下の点でのコントリビューションを歓迎します：

- より分かりやすい説明の追加
- サンプルコードの改善
- テストケースの追加
- ドキュメントの改善

## 📝 ライセンス

このプロジェクトは学習目的で作成されています。

## 🔗 参考資料

- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Domain-Driven Design](https://martinfowler.com/bliki/DomainDrivenDesign.html)
- [Go Web Programming](https://golang.org/doc/)
- [Go標準HTTP処理](https://pkg.go.dev/net/http)
- [Go Database/SQL パッケージ](https://pkg.go.dev/database/sql)