# マルチステージビルドを使用してイメージサイズを最小化
# ステージ1: ビルド環境
FROM golang:1.23-alpine AS builder

# ビルドに必要なパッケージをインストール
# git: go mod downloadでプライベートリポジトリを扱う場合に必要
# ca-certificates: HTTPSでの外部通信に必要な証明書
RUN apk add --no-cache git ca-certificates tzdata

# タイムゾーンの設定（アプリケーションでの時刻表示用）
ENV TZ=Asia/Tokyo

# 作業ディレクトリの設定
WORKDIR /app

# Go Modulesの依存関係ファイルをコピー
# これらのファイルを最初にコピーすることで、
# ソースコード変更時でも依存関係のダウンロードがキャッシュされる
COPY go.mod go.sum* ./

# 依存関係のダウンロード
# go mod downloadは依存関係をダウンロードのみ行い、ビルドは行わない
RUN go mod download

# ソースコードをコピー
COPY . .

# アプリケーションのビルド
# CGO_ENABLED=0: 静的リンクでバイナリを作成（Dockerイメージで実行しやすくなる）
# GOOS=linux: Linuxバイナリを明示的に指定
# -a: すべてのパッケージを再ビルド
# -installsuffix cgo: CGO無効時の識別子
# -ldflags: リンカーフラグ
#   -w: デバッグ情報を削除してバイナリサイズを削減
#   -s: シンボルテーブルを削除してバイナリサイズを削減
RUN CGO_ENABLED=0 GOOS=linux go build \
    -a -installsuffix cgo \
    -ldflags '-w -s' \
    -o todoapp \
    ./cmd/api

# ステージ2: 実行環境（最小構成）
FROM alpine:latest

# セキュリティアップデートとSSL証明書のインストール
# ca-certificates: HTTPS通信で必要
# tzdata: タイムゾーン情報
RUN apk --no-cache add ca-certificates tzdata

# タイムゾーンの設定
ENV TZ=Asia/Tokyo

# 非rootユーザーの作成（セキュリティ向上）
# アプリケーション専用のユーザーを作成してrootでの実行を避ける
RUN addgroup -g 1000 -S appgroup && \
    adduser -u 1000 -S appuser -G appgroup

# 作業ディレクトリの作成と所有者変更
WORKDIR /app
RUN chown -R appuser:appgroup /app

# ビルドステージからバイナリファイルをコピー
COPY --from=builder --chown=appuser:appgroup /app/todoapp .

# 非rootユーザーに切り替え
USER appuser

# アプリケーションが使用するポートを公開
# 8080はデフォルト設定、環境変数で変更可能
EXPOSE 8080

# ヘルスチェックの設定
# Dockerコンテナの健全性を監視するためのコマンド
# 30秒間隔で /health エンドポイントにアクセスしてチェック
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# アプリケーションの実行
# 実行ファイル名を直接指定
ENTRYPOINT ["./todoapp"]

# Dockerfileのベストプラクティス：
#
# 1. マルチステージビルド:
#    - ビルド環境と実行環境を分離
#    - 最終イメージサイズの最小化
#    - 不要な開発ツールの除外
#
# 2. レイヤーキャッシュの活用:
#    - 依存関係ファイルを先にコピー
#    - 変更頻度の低いレイヤーを下位に配置
#
# 3. セキュリティ:
#    - 非rootユーザーでの実行
#    - 最小限のベースイメージ使用
#    - 不要なパッケージの除外
#
# 4. 最適化:
#    - 静的バイナリの作成
#    - デバッグ情報の削除
#    - 適切な.dockerignoreの使用
#
# 5. 監視:
#    - ヘルスチェックの設定
#    - 適切なポート公開
#
# 6. 保守性:
#    - 明確なコメント
#    - 環境変数での設定可能化