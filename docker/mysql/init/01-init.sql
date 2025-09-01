-- MySQL初期化スクリプト
-- Docker Compose起動時にデータベースの初期設定を行う

-- 文字セットとコレーションを明示的に設定
-- 日本語対応のためにutf8mb4を使用
ALTER DATABASE todoapp CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- テーブルの作成は GORMのAutoMigrateに任せるため、
-- ここでは基本的なデータベース設定のみを行う

-- ユーザー権限の確認（デバッグ用）
-- SELECT User, Host FROM mysql.user WHERE User = 'todouser';
-- SHOW GRANTS FOR 'todouser'@'%';

-- 設定確認用クエリ（必要に応じて有効化）
-- SHOW VARIABLES LIKE 'character_set%';
-- SHOW VARIABLES LIKE 'collation%';

-- サンプルデータの投入（開発環境用）
-- 実際のアプリケーションテストで使用できる初期データ
-- INSERT INTO todos (title, description, is_completed, created_at, updated_at) VALUES
-- ('サンプルタスク1', '最初のタスクです', false, NOW(), NOW()),
-- ('サンプルタスク2', '完了済みのタスクです', true, NOW(), NOW()),
-- ('サンプルタスク3', 'Docker環境での動作確認用', false, NOW(), NOW());

-- パフォーマンス設定の確認
-- SHOW VARIABLES LIKE 'innodb_buffer_pool_size';
-- SHOW VARIABLES LIKE 'max_connections';