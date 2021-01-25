package music

import (
	"github.com/zackartz/cmdlr2"
	"synergy/music"
)

var SkipCommand = &cmdlr2.Command{
	Name:        "skip",
	Description: "Skips a song in the queue.",
	Example:     "skip",
	Handler: func(ctx *cmdlr2.Ctx) {
		guild, err := ctx.Client.Guild(ctx.Event.Message.GuildID).Get()
		if err != nil {
			return
		}

		go music.SkipMusic(guild, ctx.Event.Message.Author.ID)
	},
}
