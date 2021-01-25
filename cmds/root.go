package cmds

import (
	"github.com/zackartz/cmdlr2"
	"synergy/cmds/music"
)

var CommandList []*cmdlr2.Command

func init() {
	CommandList = append(CommandList, music.PlayCommand)
	CommandList = append(CommandList, music.StopCommand)
	CommandList = append(CommandList, music.SkipCommand)
	CommandList = append(CommandList, music.PauseCommand)
	CommandList = append(CommandList, music.UnPauseCommand)
}
