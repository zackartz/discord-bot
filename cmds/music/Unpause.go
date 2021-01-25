package music

import (
	"github.com/zackartz/cmdlr2"
	"synergy/music"
)

var UnPauseCommand = &cmdlr2.Command{
	Name:        "unpause",
	Description: "Unpauses a paused queue.",
	Example:     "unpause",
	Handler: func(ctx *cmdlr2.Ctx) {
		guild, err := ctx.Client.Guild(ctx.Event.Message.GuildID).Get()
		if err != nil {
			return
		}

		go music.UnPauseStream(guild, ctx.Event.Message.Author.ID)
	},
}
