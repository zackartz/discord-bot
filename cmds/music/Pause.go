package music

import (
	"github.com/zackartz/cmdlr2"
	"synergy/music"
)

var PauseCommand = &cmdlr2.Command{
	Name:        "pause",
	Description: "Pauses the bot's playback.",
	Example:     "pause",
	Handler: func(ctx *cmdlr2.Ctx) {
		guild, err := ctx.Client.Guild(ctx.Event.Message.GuildID).Get()
		if err != nil {
			return
		}

		go music.PauseStream(guild, ctx.Event.Message.Author.ID)
	},
}
