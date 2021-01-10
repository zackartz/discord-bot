package role

import (
	"fmt"
	"github.com/andersfylling/disgord"
	"github.com/zackartz/cmdlr2"
	"synergy/db"
)

var CreateMessageRoleCommand = &cmdlr2.Command{
	Name:        "create",
	Description: "Create a new role selection message.",
	Usage:       "create",
	Example:     "create",
	Handler: func(ctx *cmdlr2.Ctx) {
		embed := &disgord.Embed{
			Title:       "Roles",
			Description: "Use the following to pick roles!",
			Fields:      []*disgord.EmbedField{},
		}

		if ctx.Event.Message.Author.ID == 133314498214756352 || ctx.Event.Message.Author.ID == 271787171889807360 {
			m, _ := ctx.Client.Channel(ctx.Event.Message.ChannelID).CreateMessage(&disgord.CreateMessageParams{Embed: embed})

			err := db.CreateRoleMessage(int64(m.ID), int64(ctx.Event.Message.ChannelID), int64(ctx.Event.Message.GuildID))
			if err != nil {
				_ = ctx.Client.Channel(ctx.Event.Message.ChannelID).DeleteMessages(&disgord.DeleteMessagesParams{Messages: []disgord.Snowflake{m.ID}})
				ctx.ResponseText(fmt.Sprintf("%v", err))
			}
		}
	},
}
