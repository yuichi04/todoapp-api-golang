package entity

import (
	"time"
)

// Todo はタスク管理システムの中核となるドメインエンティティです
// エンティティは業務データと業務ロジックを持つオブジェクトで、
// 一意性を持つID（識別子）によって他のオブジェクトと区別されます
type Todo struct {
	// ID は各Todoを一意に識別するための主キーです
	// Goでは大文字で始まるフィールドがpublicになります
	ID int `json:"id" gorm:"primaryKey;autoIncrement"`

	// Title はタスクのタイトル（件名）を表します
	// `gorm:"not null"` でデータベース上でNOT NULL制約を付けています
	// `json:"title"` でJSONシリアライゼーション時のキー名を指定しています
	Title string `json:"title" gorm:"not null"`

	// Description はタスクの詳細説明を格納します
	// 説明は任意項目なのでNOT NULL制約は付けていません
	Description string `json:"description"`

	// IsCompleted はタスクの完了状態を表すbool型フィールドです
	// `gorm:"default:false"` でデフォルト値をfalse（未完了）に設定
	IsCompleted bool `json:"is_completed" gorm:"default:false"`

	// CreatedAt はレコードの作成日時を自動的に記録します
	// GORMでは time.Time 型の CreatedAt フィールドを自動的に現在時刻で設定します
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`

	// UpdatedAt はレコードの更新日時を自動的に記録します
	// GORMでは time.Time 型の UpdatedAt フィールドを更新時に自動的に現在時刻で設定します
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName はGORMでデータベーステーブル名をカスタマイズするためのメソッドです
// デフォルトでは構造体名の複数形（todos）になりますが、明示的に指定することで
// テーブル名を制御できます（保守性の向上）
func (Todo) TableName() string {
	return "todos"
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