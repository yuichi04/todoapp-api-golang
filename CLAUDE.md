# CLAUDE.md

このファイルは、このリポジトリでコードを扱う際のClaude Code (claude.ai/code) への指針を提供します。

## プロジェクト概要

このプロジェクトはAPI開発学習のための教科書として設計されたGo言語ベースのTodo APIプロジェクトです。タスク管理システム実装を理解するための教育リソースとして機能します。

## 一般的なコマンド

### プロジェクトセットアップ
```bash
# Goモジュールの初期化
go mod init todoapp-api-golang

# 依存関係のインストール
go mod tidy
```

### 開発

#### 通常の開発
```bash
# アプリケーションの実行
go run ./cmd/api

# アプリケーションのビルド
go build -o todoapp ./cmd/api

# ビルドしたバイナリの実行
./todoapp
```

#### ホットリロード開発（推奨）
```bash
# Air（ホットリロードツール）を使用した開発サーバー起動
# ファイル変更時に自動再起動される
make dev-hot

# または直接Airを使用
air -c .air.toml

# または開発スクリプトを使用（Linux/macOS）
./scripts/dev.sh

# Windows環境
scripts\dev.bat
```

### テスト
```bash
# すべてのテストを実行
go test ./...

# カバレッジ付きでテストを実行
go test -cover ./...

# 特定のテストを実行
go test -run TestFunctionName ./path/to/package
```

### コード品質
```bash
# コードのフォーマット
go fmt ./...

# 潜在的な問題のチェック
go vet ./...
```

## 教育的コンテキスト

このプロジェクトは以下の学習リソースとして構成されています：
- RESTful API設計と実装
- データベース操作とモデリング
- GoでのHTTPリクエスト処理
- タスク/Todo管理システムアーキテクチャ
- Go Web開発のベストプラクティス

## Air（ホットリロード）について

このプロジェクトは開発効率化のためにAir（Live reload for Go apps）を使用してホットリロード機能を提供します。

### Airの特徴
- ファイル変更を自動検出
- 自動ビルドと再起動
- カスタマイズ可能な設定ファイル（.air.toml）
- 開発時のみ使用（本番環境では不要）

### 設定ファイル
- `.air.toml` - Air設定ファイル
- `scripts/dev.sh` - Linux/macOS用開発スクリプト
- `scripts/dev.bat` - Windows用開発スクリプト

### 対応環境
- Linux
- macOS
- Windows

## 学習予定トピック

- 標準パッケージでのHTTPサーバー実装
- タスク管理のCRUD操作
- Clean Architecture実装
- database/sqlパッケージの使用
- エラーハンドリング
- HTTPミドルウェア
- Docker化
- ホットリロード開発環境