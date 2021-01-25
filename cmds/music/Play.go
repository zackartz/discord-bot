package music

import (
	"github.com/zackartz/cmdlr2"
	"strings"
	"synergy/music"
)

var PlayCommand = &cmdlr2.Command{
	Name:        "play",
	Description: "Adds a new song to the current queue.",
	Example:     "play [url/search terms]",
	Handler: func(ctx *cmdlr2.Ctx) {
		var url string

		if ctx.Args.Get(0).Raw() == "" {
			ctx.ResponseText("Try $play [song/url]")
		}

		if strings.HasPrefix(ctx.Args.Get(0).Raw(), "https://") {
			url = ctx.Args.Get(0).Raw()
		} else {
			url, _ = music.GetVideoFromName(ctx.Args.Raw())
		}

		go music.StreamToDiscord(ctx, url)
	},
}
