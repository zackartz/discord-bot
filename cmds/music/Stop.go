package music

import (
	"github.com/zackartz/cmdlr2"
	"synergy/music"
)

var StopCommand = &cmdlr2.Command{
	Name:        "stop",
	Description: "Stops the queue and disconnects the bot.",
	Example:     "stop",
	Handler: func(ctx *cmdlr2.Ctx) {
		guild, err := ctx.Client.Guild(ctx.Event.Message.GuildID).Get()
		if err != nil {
			return
		}

		go music.StopStream(guild, ctx.Event.Message.Author.ID)
	},
}
