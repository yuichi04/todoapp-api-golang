#!/bin/bash

# 開発環境セットアップスクリプト
# Air（ホットリロード）を使用した開発サーバーの起動
#
# 使用方法:
# ./scripts/dev.sh
#
# このスクリプトが行うこと:
# 1. Airのインストール確認と自動インストール
# 2. 必要なディレクトリの作成
# 3. 開発サーバーの起動（ホットリロード有効）

set -e  # エラー時にスクリプトを終了

echo "=== Todo API 開発環境セットアップ ==="

# 現在のディレクトリを確認
if [ ! -f "go.mod" ]; then
    echo "エラー: go.modファイルが見つかりません。"
    echo "プロジェクトのルートディレクトリで実行してください。"
    exit 1
fi

# 1. Airのインストール確認
echo "Air（ホットリロードツール）をチェック中..."

if ! command -v air &> /dev/null; then
    echo "Airがインストールされていません。インストールを開始..."
    
    # Goがインストールされているかチェック
    if ! command -v go &> /dev/null; then
        echo "エラー: Goがインストールされていません。"
        echo "先にGoをインストールしてください: https://golang.org/dl/"
        exit 1
    fi
    
    # Airをインストール
    echo "Air（最新版）をインストール中..."
    go install github.com/air-verse/air@latest
    
    # PATHの確認
    if ! command -v air &> /dev/null; then
        echo "警告: airコマンドが見つかりません。"
        echo "以下のいずれかを実行してください:"
        echo "1. export PATH=\$PATH:\$(go env GOPATH)/bin"
        echo "2. シェルを再起動"
        echo ""
        echo "または、以下のコマンドで直接実行："
        echo "\$(go env GOPATH)/bin/air -c .air.toml"
        exit 1
    fi
else
    echo "Air はすでにインストールされています: $(air -v 2>&1 | head -n1)"
fi

# 2. 必要なディレクトリの作成
echo "開発用ディレクトリを作成中..."

# 一時ディレクトリの作成（バイナリファイル用）
mkdir -p tmp

# ログディレクトリの作成（オプション）
mkdir -p logs

# 3. 設定ファイルの確認
if [ ! -f ".air.toml" ]; then
    echo "エラー: .air.tomlファイルが見つかりません。"
    echo "設定ファイルが必要です。"
    exit 1
fi

echo "設定ファイル(.air.toml)を確認: OK"

# 4. 環境変数の設定（開発環境用）
export APP_ENV=development
export SERVER_HOST=localhost
export SERVER_PORT=8080
export DB_DRIVER=mysql
export DB_HOST=localhost
export DB_PORT=3306
export DB_NAME=todoapp
export DB_USER=todouser
export DB_PASSWORD=todopass

echo ""
echo "=== 開発環境設定 ==="
echo "APP_ENV: $APP_ENV"
echo "サーバー: http://$SERVER_HOST:$SERVER_PORT"
echo "データベース: $DB_DRIVER://$DB_USER@$DB_HOST:$DB_PORT/$DB_NAME"
echo ""

# 5. データベースの確認メッセージ
echo "=== 重要: データベースについて ==="
echo "開発を始める前に、MySQLデータベースが起動していることを確認してください。"
echo ""
echo "Dockerを使用する場合:"
echo "  docker-compose up -d mysql"
echo ""
echo "ローカルMySQLを使用する場合:"
echo "  上記の接続情報でデータベースが利用可能であることを確認してください。"
echo ""

# 6. 使用方法の説明
echo "=== Air（ホットリロード）の使用方法 ==="
echo "• ファイルを編集・保存すると自動的にアプリケーションが再起動されます"
echo "• 終了するには Ctrl+C を押してください"
echo "• ログは tmp/build-errors.log に出力されます"
echo ""

# 7. 開発サーバーの起動
echo "開発サーバーを起動します..."
echo "開発サーバーが起動したら http://localhost:8080/health でヘルスチェックを確認できます"
echo ""

# Airの実行
air -c .air.toml

# スクリプト使用時の学習ポイント：
#
# 1. シェルスクリプトの基本:
#    - set -e でエラー時の終了
#    - 条件分岐と環境変数設定
#    - コマンド存在確認
#
# 2. Go開発ツールの管理:
#    - go install による依存ツールのインストール
#    - GOPATH/bin の理解
#
# 3. 開発環境の自動化:
#    - 環境変数による設定管理
#    - 必要なディレクトリの自動作成
#
# 4. エラーハンドリング:
#    - 前提条件の確認
#    - わかりやすいエラーメッセージ
#
# 5. 開発者体験の向上:
#    - セットアップの自動化
#    - 明確な使用手順の提示