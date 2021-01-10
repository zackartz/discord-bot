package events

import (
	"fmt"
	"github.com/andersfylling/disgord"
	"synergy/db"
)

func EmojiAdd(s disgord.Session, h *disgord.MessageReactionAdd) {
	r, err := db.GetRoleByEmoji(h.PartialEmoji.Name)
	if err != nil {
		return
	}

	u, err := s.CurrentUser().Get()
	if err != nil {
		return
	}

	if h.UserID == u.ID {
		return
	}

	rm, err := db.GetRoleMessageByMessageID(r.MessageID)
	if err != nil {
		return
	}

	err = s.Guild(disgord.Snowflake(rm.GuildID)).Member(h.UserID).AddRole(disgord.Snowflake(r.ID))
	if err != nil {
		s.Channel(h.ChannelID).CreateMessage(&disgord.CreateMessageParams{
			Content: fmt.Sprintf("%v", err),
		})
	}
}
