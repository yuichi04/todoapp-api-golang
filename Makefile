# Makefile for Todo API with Standard Packages
# プロジェクトの一般的なタスクを簡素化するためのファイル
# Air（ホットリロード）による開発効率化機能を追加

.PHONY: help setup run build test clean docker-setup docker-start docker-stop docker-logs docker-clean dev-hot install-air

# デフォルトターゲット
help: ## このヘルプメッセージを表示
	@echo "利用可能なコマンド:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

# 開発環境
setup: ## 開発環境のセットアップ
	go mod tidy
	cp .env.example .env

run: ## アプリケーションの実行（開発モード）
	go run cmd/api/main.go

dev-hot: install-air ## ホットリロード付き開発サーバー起動（Air使用）
	@echo "ホットリロード開発サーバーを起動中..."
	@echo "ファイルを編集すると自動的に再起動されます"
	@mkdir -p tmp logs
	@APP_ENV=development \
	 SERVER_HOST=localhost \
	 SERVER_PORT=8080 \
	 DB_DRIVER=mysql \
	 DB_HOST=localhost \
	 DB_PORT=3306 \
	 DB_NAME=todoapp \
	 DB_USER=todouser \
	 DB_PASSWORD=todopass \
	 air -c .air.toml

install-air: ## Air（ホットリロードツール）をインストール
	@command -v air >/dev/null 2>&1 || { \
		echo "Airをインストール中..."; \
		go install github.com/air-verse/air@latest; \
	}

build: ## アプリケーションのビルド
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o todoapp cmd/api/main.go

test: ## テストの実行
	go test ./...

test-coverage: ## カバレッジ付きテスト
	go test -cover ./...

clean: ## ビルド成果物のクリーンアップ
	rm -f todoapp
	go clean

# フォーマットとコード検査
fmt: ## コードフォーマット
	go fmt ./...

vet: ## コード検査
	go vet ./...

lint: fmt vet ## フォーマットと検査

# Docker関連
docker-setup: ## Docker環境の初期セットアップ
	./scripts/docker-setup.sh setup

docker-start: ## Dockerでアプリケーション起動
	./scripts/docker-setup.sh start

docker-stop: ## Dockerアプリケーション停止
	./scripts/docker-setup.sh stop

docker-restart: ## Dockerアプリケーション再起動
	./scripts/docker-setup.sh restart

docker-logs: ## Dockerログ表示
	./scripts/docker-setup.sh logs

docker-clean: ## Dockerクリーンアップ
	./scripts/docker-setup.sh clean

docker-reset: ## Docker完全リセット
	./scripts/docker-setup.sh reset

# 開発用ショートカット
dev: setup run ## 開発環境セットアップ＋実行
dev-reload: dev-hot ## ホットリロード開発（エイリアス）

docker-dev: docker-setup docker-start ## Docker開発環境セットアップ＋実行