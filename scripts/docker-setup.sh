#!/bin/bash

# Docker環境セットアップスクリプト
# Todo APIアプリケーションをDocker環境で起動するためのスクリプト

set -e  # エラー時にスクリプトを終了

# 色付きメッセージ用の定数
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# ヘルパー関数
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# スクリプトの使用方法を表示
show_usage() {
    echo "使用方法: $0 [コマンド]"
    echo ""
    echo "利用可能なコマンド:"
    echo "  setup     - Docker環境の初期セットアップ"
    echo "  start     - アプリケーションの起動"
    echo "  stop      - アプリケーションの停止"
    echo "  restart   - アプリケーションの再起動"
    echo "  logs      - ログの表示"
    echo "  clean     - データベースとイメージのクリーンアップ"
    echo "  reset     - 完全リセット（データ削除）"
    echo "  help      - このヘルプメッセージを表示"
}

# Dockerの動作確認
check_docker() {
    print_info "Dockerの動作確認..."
    
    if ! command -v docker &> /dev/null; then
        print_error "Dockerがインストールされていません"
        exit 1
    fi
    
    if ! docker info &> /dev/null; then
        print_error "Dockerが起動していません"
        exit 1
    fi
    
    if ! command -v docker-compose &> /dev/null; then
        print_error "Docker Composeがインストールされていません"
        exit 1
    fi
    
    print_success "Docker環境が正常に動作しています"
}

# 初期セットアップ
setup() {
    print_info "Docker環境の初期セットアップを開始..."
    
    check_docker
    
    # .env.dockerファイルがない場合はコピー
    if [ ! -f .env.docker ]; then
        print_warning ".env.dockerファイルが見つかりません"
        if [ -f .env.example ]; then
            cp .env.example .env.docker
            print_info ".env.exampleから.env.dockerを作成しました"
        fi
    fi
    
    # MySQLの初期化ディレクトリ作成
    mkdir -p docker/mysql/init
    
    print_success "初期セットアップが完了しました"
}

# アプリケーション起動
start() {
    print_info "アプリケーションを起動しています..."
    
    check_docker
    
    # イメージのビルドとコンテナ起動
    docker-compose --env-file .env.docker up -d --build
    
    print_success "アプリケーションが起動しました"
    print_info "API: http://localhost:8080"
    print_info "phpMyAdmin: http://localhost:8081"
    print_info "ログ確認: $0 logs"
}

# アプリケーション停止
stop() {
    print_info "アプリケーションを停止しています..."
    
    docker-compose --env-file .env.docker down
    
    print_success "アプリケーションが停止しました"
}

# アプリケーション再起動
restart() {
    print_info "アプリケーションを再起動しています..."
    
    stop
    sleep 2
    start
}

# ログ表示
show_logs() {
    print_info "ログを表示します (Ctrl+C で終了)..."
    
    docker-compose --env-file .env.docker logs -f
}

# クリーンアップ
clean() {
    print_warning "データベースデータとイメージを削除します"
    read -p "続行しますか？ (y/N): " -n 1 -r
    echo
    
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        print_info "クリーンアップを実行しています..."
        
        # コンテナ停止と削除
        docker-compose --env-file .env.docker down -v
        
        # イメージ削除
        docker rmi todoapp-api-golang-todoapp 2>/dev/null || true
        
        # 未使用のボリュームとネットワークを削除
        docker volume prune -f
        docker network prune -f
        
        print_success "クリーンアップが完了しました"
    else
        print_info "クリーンアップがキャンセルされました"
    fi
}

# 完全リセット
reset() {
    print_warning "すべてのデータとイメージを削除します（復元不可能）"
    read -p "本当に続行しますか？ (y/N): " -n 1 -r
    echo
    
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        print_info "完全リセットを実行しています..."
        
        # すべてのコンテナを停止
        docker-compose --env-file .env.docker down -v --rmi all
        
        # 関連するボリューム削除
        docker volume rm todoapp-api-golang_mysql_data 2>/dev/null || true
        
        print_success "完全リセットが完了しました"
    else
        print_info "リセットがキャンセルされました"
    fi
}

# メイン処理
main() {
    case "${1:-help}" in
        setup)
            setup
            ;;
        start)
            start
            ;;
        stop)
            stop
            ;;
        restart)
            restart
            ;;
        logs)
            show_logs
            ;;
        clean)
            clean
            ;;
        reset)
            reset
            ;;
        help|*)
            show_usage
            ;;
    esac
}

# スクリプト実行
main "$@"