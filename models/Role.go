package models

import "time"

type Role struct {
	ID        int64
	Name      string
	Emoji     string
	EmojiID   int64
	MessageID int64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
