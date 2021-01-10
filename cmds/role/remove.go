package role

import (
	"fmt"
	"github.com/andersfylling/disgord"
	"github.com/zackartz/cmdlr2"
	"strconv"
	"synergy/db"
	"synergy/models"
)

var RemoveRoleCommand = &cmdlr2.Command{
	Name:        "remove",
	Description: "Remove a role to the commandlist",
	Usage:       "remove @role",
	Example:     "remove @DPS",
	Handler: func(ctx *cmdlr2.Ctx) {
		rID, err := strconv.ParseInt(ctx.Args.Get(0).AsRoleMentionID(), 10, 64)
		if err != nil {
			ctx.ResponseText(fmt.Sprintf("%v", err))
			return
		}

		var role *disgord.Role
		roles, err := ctx.Client.Guild(ctx.Event.Message.GuildID).GetRoles()
		if err != nil {
			ctx.ResponseText(fmt.Sprintf("%v", err))
			return
		}

		role, err = getRoles(roles, rID)
		if err != nil {
			ctx.ResponseText(fmt.Sprintf("%v", err))
			return
		}

		msg, err := db.GetRoleMessageByChannelID(int64(ctx.Event.Message.ChannelID))
		if err != nil {
			ctx.ResponseText(fmt.Sprintf("%v", err))
			return
		}

		var r *models.Role
		for _, x := range msg.Roles {
			if x.ID == int64(role.ID) {
				r = x
				break
			}
		}

		if r == nil {
			return
		}

		var emoji *disgord.Emoji

		if r.EmojiID != 0 {
			emoji, err = ctx.Client.Guild(ctx.Event.Message.GuildID).Emoji(disgord.Snowflake(r.EmojiID)).Get()
			if err != nil {
				return
			}

			err = db.RemoveRoleFromMessge(int64(role.ID))
			if err != nil {
				return
			}
		} else {
			emoji = &disgord.Emoji{Name: r.Emoji}
			err = db.RemoveRoleFromMessge(int64(role.ID))
			if err != nil {
				return
			}
		}

		msg, err = db.GetRoleMessageByChannelID(int64(ctx.Event.Message.ChannelID))
		if err != nil {
			ctx.ResponseText(fmt.Sprintf("%v", err))
			return
		}

		embed := renderEmbed(msg)

		_, err = ctx.Client.Channel(ctx.Event.Message.ChannelID).Message(disgord.Snowflake(msg.ID)).Update().SetEmbed(embed).Execute()

		err = ctx.Client.Channel(disgord.Snowflake(msg.ChannelID)).Message(disgord.Snowflake(msg.ID)).Reaction(emoji).DeleteOwn()
		if err != nil {
			ctx.ResponseText(fmt.Sprintf("%v", err))
		}

		_ = ctx.Client.Channel(ctx.Event.Message.ChannelID).Message(ctx.Event.Message.ID).Delete()
	},
}
