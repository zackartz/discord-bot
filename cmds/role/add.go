package role

import (
	"errors"
	"fmt"
	"github.com/andersfylling/disgord"
	"github.com/zackartz/cmdlr2"
	"regexp"
	"strconv"
	"strings"
	"synergy/db"
	"synergy/models"
	"time"
)

var AddRoleCommand = &cmdlr2.Command{
	Name:        "add",
	Description: "Add a role to the commandlist",
	Usage:       "add @role :emoji:",
	Example:     "add @DPS :DPS:",
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

		var emoji *disgord.Emoji

		if strings.HasPrefix(ctx.Args.Get(1).Raw(), "<") {
			regex := regexp.MustCompile("\\d+")
			id := regex.Find([]byte(ctx.Args.Get(1).Raw()))
			eID, err := strconv.ParseInt(string(id), 10, 64)
			if err != nil {
				ctx.ResponseText(fmt.Sprintf("%v", err))
				return
			}

			emoji, err = ctx.Client.Guild(ctx.Event.Message.GuildID).Emoji(disgord.Snowflake(eID)).Get()
			if err != nil {
				ctx.ResponseText(fmt.Sprintf("%v", err))
				return
			}

			err = db.AddRoleToMessage(int64(role.ID), msg.ID, eID, emoji.Name, role.Name)
			if err != nil {
				ctx.ResponseText(fmt.Sprintf("%v", err))
				return
			}
		} else {
			emoji = &disgord.Emoji{Name: ctx.Args.Get(1).Raw()}
			err = db.AddRoleToMessage(int64(role.ID), msg.ID, 0, ctx.Args.Get(1).Raw(), role.Name)
			if err != nil {
				ctx.ResponseText(fmt.Sprintf("%v", err))
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

		err = ctx.Client.Channel(disgord.Snowflake(msg.ChannelID)).Message(disgord.Snowflake(msg.ID)).Reaction(emoji).Create()
		if err != nil {
			ctx.ResponseText(fmt.Sprintf("%v", err))
		}

		_ = ctx.Client.Channel(ctx.Event.Message.ChannelID).Message(ctx.Event.Message.ID).Delete()
	},
}

func getRoles(roles []*disgord.Role, rID int64) (*disgord.Role, error) {
	for _, r := range roles {
		if r.ID == disgord.Snowflake(rID) {
			return r, nil
		}
	}
	return nil, errors.New("couldn't find role")
}

func renderEmbed(msg *models.RoleMessage) *disgord.Embed {
	var fields []*disgord.EmbedField

	for _, x := range msg.Roles {
		var val string

		if x.EmojiID != 0 {
			val = fmt.Sprintf("<:%s:%d>", x.Emoji, x.EmojiID)
		} else {
			val = x.Emoji
		}

		fields = append(fields, &disgord.EmbedField{
			Name:   x.Name,
			Value:  val,
			Inline: true,
		})
	}

	return &disgord.Embed{
		Title:       "Roles",
		Type:        "rich",
		Description: "Use the following to pick roles!",
		Timestamp: disgord.Time{
			Time: time.Now(),
		},
		Color:  0xffff00,
		Fields: fields,
	}
}
