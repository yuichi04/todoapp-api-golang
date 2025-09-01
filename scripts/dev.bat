@echo off
REM 開発環境セットアップスクリプト（Windows用）
REM Air（ホットリロード）を使用した開発サーバーの起動
REM
REM 使用方法:
REM scripts\dev.bat
REM
REM このスクリプトが行うこと:
REM 1. Airのインストール確認と自動インストール
REM 2. 必要なディレクトリの作成
REM 3. 開発サーバーの起動（ホットリロード有効）

echo === Todo API 開発環境セットアップ（Windows）===

REM 現在のディレクトリを確認
if not exist "go.mod" (
    echo エラー: go.modファイルが見つかりません。
    echo プロジェクトのルートディレクトリで実行してください。
    pause
    exit /b 1
)

REM 1. Goのインストール確認
echo Goのインストールをチェック中...
go version >nul 2>&1
if errorlevel 1 (
    echo エラー: Goがインストールされていません。
    echo 先にGoをインストールしてください: https://golang.org/dl/
    pause
    exit /b 1
)

REM 2. Airのインストール確認とインストール
echo Air（ホットリロードツール）をチェック中...
air -v >nul 2>&1
if errorlevel 1 (
    echo Airがインストールされていません。インストールを開始...
    go install github.com/air-verse/air@latest
    
    REM インストール確認
    air -v >nul 2>&1
    if errorlevel 1 (
        echo 警告: airコマンドが見つかりません。
        echo 環境変数PATHにGoのbinディレクトリが含まれているか確認してください。
        echo 一般的なパス: %%USERPROFILE%%\go\bin
        pause
        exit /b 1
    )
) else (
    echo Air はすでにインストールされています
    air -v
)

REM 3. 必要なディレクトリの作成
echo 開発用ディレクトリを作成中...

REM 一時ディレクトリの作成（バイナリファイル用）
if not exist "tmp" mkdir tmp

REM ログディレクトリの作成（オプション）
if not exist "logs" mkdir logs

REM 4. 設定ファイルの確認
if not exist ".air.toml" (
    echo エラー: .air.tomlファイルが見つかりません。
    echo 設定ファイルが必要です。
    pause
    exit /b 1
)

echo 設定ファイル(.air.toml)を確認: OK

REM 5. 環境変数の設定（開発環境用）
set APP_ENV=development
set SERVER_HOST=localhost
set SERVER_PORT=8080
set DB_DRIVER=mysql
set DB_HOST=localhost
set DB_PORT=3306
set DB_NAME=todoapp
set DB_USER=todouser
set DB_PASSWORD=todopass

echo.
echo === 開発環境設定 ===
echo APP_ENV: %APP_ENV%
echo サーバー: http://%SERVER_HOST%:%SERVER_PORT%
echo データベース: %DB_DRIVER%://%DB_USER%@%DB_HOST%:%DB_PORT%/%DB_NAME%
echo.

REM 6. データベースの確認メッセージ
echo === 重要: データベースについて ===
echo 開発を始める前に、MySQLデータベースが起動していることを確認してください。
echo.
echo Dockerを使用する場合:
echo   docker-compose up -d mysql
echo.
echo ローカルMySQLを使用する場合:
echo   上記の接続情報でデータベースが利用可能であることを確認してください。
echo.

REM 7. 使用方法の説明
echo === Air（ホットリロード）の使用方法 ===
echo • ファイルを編集・保存すると自動的にアプリケーションが再起動されます
echo • 終了するには Ctrl+C を押してください
echo • ログは tmp\build-errors.log に出力されます
echo.

REM 8. 開発サーバーの起動
echo 開発サーバーを起動します...
echo 開発サーバーが起動したら http://localhost:8080/health でヘルスチェックを確認できます
echo.

REM Airの実行
air -c .air.toml

REM スクリプト終了時の処理
if errorlevel 1 (
    echo.
    echo 開発サーバーの起動中にエラーが発生しました。
    echo ログファイル tmp\build-errors.log を確認してください。
    pause
)