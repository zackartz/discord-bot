package music

import (
	"github.com/andersfylling/disgord"
	"github.com/kkdai/youtube/v2"
	"github.com/zackartz/cmdlr2"
)

const (
	Channels  int = 2
	FrameRate int = 48000
	FrameSize int = 960
	MaxBytes  int = (FrameSize * 2) * 2
)

type ServerQueue struct {
	MusicQueue []*MusicInfo
	vc         disgord.VoiceConnection
}

type MusicInfo struct {
	SongInfo      *youtube.Video
	Url           string
	TextChannelID disgord.Snowflake
	GuildID       disgord.Snowflake
	ChannelID     disgord.Snowflake
	Pause         chan bool
	Stop          bool
	Skip          chan bool
	Ctx           *cmdlr2.Ctx
}

var (
	ServersVC  = map[disgord.Snowflake]*ServerQueue{}
	QueueBlock = map[string]bool{}
	QueueChan  = make(chan *MusicInfo, 101)
)
