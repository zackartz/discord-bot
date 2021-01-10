package models

import "time"

type RoleMessage struct {
	ID        int64     `json:"id"`
	ChannelID int64     `json:"channel_id"`
	GuildID   int64     `json:"guild_id"`
	Roles     []*Role   `gorm:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}
