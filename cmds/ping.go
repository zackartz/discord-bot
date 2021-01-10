package cmds

import "github.com/zackartz/cmdlr2"

var PingCommand = &cmdlr2.Command{
	Name:        "ping",
	Description: "Pings",
	Example:     "ping",
	Handler: func(ctx *cmdlr2.Ctx) {
		ctx.ResponseText("pong")
	},
}

func init() {
	CommandList = append(CommandList, PingCommand)
}
