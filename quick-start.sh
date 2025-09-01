#!/bin/bash

# Todo API クイックスタートスクリプト
# 初回セットアップから動作確認まで自動実行

set -e

echo "🚀 Todo API クイックスタートを開始します..."

# Docker環境の確認
if command -v docker &> /dev/null && command -v docker-compose &> /dev/null; then
    echo "✅ Docker環境が利用可能です"
    
    # Docker環境でセットアップ
    echo "📦 Docker環境でセットアップ中..."
    make docker-setup
    
    echo "🔄 アプリケーションを起動中..."
    make docker-start
    
    echo "⏳ アプリケーションの起動を待機中..."
    sleep 10
    
    echo "🧪 動作確認を実行中..."
    
    # ヘルスチェック
    echo "  ヘルスチェック..."
    curl -f http://localhost:8080/health > /dev/null 2>&1 && echo "  ✅ ヘルスチェック成功" || echo "  ❌ ヘルスチェック失敗"
    
    # サンプルTodo作成
    echo "  サンプルTodo作成..."
    RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/todos \
        -H "Content-Type: application/json" \
        -d '{"title":"Docker環境テスト","description":"自動作成されたサンプルタスク"}')
    
    if echo "$RESPONSE" | grep -q '"id"'; then
        echo "  ✅ Todo作成成功"
        TODO_ID=$(echo "$RESPONSE" | grep -o '"id":[0-9]*' | cut -d':' -f2)
        echo "  📝 作成されたTodo ID: $TODO_ID"
    else
        echo "  ❌ Todo作成失敗"
    fi
    
    # Todo一覧取得
    echo "  Todo一覧取得..."
    curl -f http://localhost:8080/api/v1/todos > /dev/null 2>&1 && echo "  ✅ Todo一覧取得成功" || echo "  ❌ Todo一覧取得失敗"
    
    echo ""
    echo "🎉 セットアップ完了！"
    echo ""
    echo "📍 アクセス情報:"
    echo "  • API: http://localhost:8080"
    echo "  • phpMyAdmin: http://localhost:8081"
    echo ""
    echo "🔧 基本操作:"
    echo "  • ログ確認: make docker-logs"
    echo "  • 停止: make docker-stop"
    echo "  • 再起動: make docker-restart"
    echo ""
    
else
    echo "⚠️  Docker環境が利用できません"
    echo "ローカル環境でのセットアップを実行します..."
    
    # Go環境の確認
    if ! command -v go &> /dev/null; then
        echo "❌ Go環境がインストールされていません"
        echo "Go 1.21以上をインストールしてください"
        exit 1
    fi
    
    echo "✅ Go環境が利用可能です"
    
    # ローカル環境でセットアップ
    echo "📦 依存関係をインストール中..."
    go mod tidy
    
    echo "⚙️  設定ファイルを作成中..."
    if [ ! -f .env ]; then
        cp .env.example .env
        echo "📝 .envファイルを作成しました"
        echo "⚠️  データベース設定を .env ファイルで確認・変更してください"
    fi
    
    echo ""
    echo "🎉 セットアップ完了！"
    echo ""
    echo "📋 次のステップ:"
    echo "  1. データベースを準備してください"
    echo "  2. .env ファイルでデータベース接続設定を確認してください"
    echo "  3. make run でアプリケーションを起動してください"
    echo ""
fi