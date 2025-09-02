package entity

import (
	"time"
)

// Todo はタスク管理システムの中核となるドメインエンティティです
// エンティティは業務データと業務ロジックを持つオブジェクトで、
// 一意性を持つID（識別子）によって他のオブジェクトと区別されます
//
// 標準パッケージでのエンティティ設計の学習ポイント：
// 1. 構造体タグはJSONのみ（ORMタグは使用しない）
// 2. データベース制約は実装層で管理
// 3. ビジネスロジックをメソッドとして定義
// 4. 時刻フィールドの適切な管理
type Todo struct {
	// ID は各Todoを一意に識別するための主キーです
	// Goでは大文字で始まるフィールドがpublicになります
	// データベースでは AUTO_INCREMENT の主キーとして扱われます
	ID int `json:"id"`

	// Title はタスクのタイトル（件名）を表します
	// `json:"title"` でJSONシリアライゼーション時のキー名を指定
	// データベースでのNOT NULL制約は実装層で管理されます
	Title string `json:"title"`

	// Description はタスクの詳細説明を格納します
	// 説明は任意項目として扱われます
	Description string `json:"description"`

	// IsCompleted はタスクの完了状態を表すbool型フィールドです
	// デフォルト値（false = 未完了）の設定は実装層で行います
	IsCompleted bool `json:"is_completed"`

	// CreatedAt はレコードの作成日時を記録します
	// 標準パッケージでは明示的に現在時刻を設定する必要があります
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt はレコードの更新日時を記録します
	// 更新時には明示的に現在時刻を設定する必要があります
	UpdatedAt time.Time `json:"updated_at"`
}

// IsValid はTodoエンティティのビジネスルールを検証するメソッドです
// ドメイン層でのバリデーションロジックを担当します
// 戻り値がtrueなら有効、falseなら無効なデータです
func (t *Todo) IsValid() bool {
	// タイトルが空文字でないかチェック
	// strings.TrimSpace() で前後の空白を除去してから長さをチェックしています
	return len(t.Title) > 0 && len(t.Title) <= 100
}

// MarkAsCompleted はタスクを完了状態にするビジネスロジックです
// エンティティ内でのステート変更ロジックをカプセル化しています
func (t *Todo) MarkAsCompleted() {
	t.IsCompleted = true
}

// MarkAsIncomplete はタスクを未完了状態に戻すビジネスロジックです
func (t *Todo) MarkAsIncomplete() {
	t.IsCompleted = false
}
