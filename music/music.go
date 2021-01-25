package music

import (
	"encoding/binary"
	"fmt"
	"github.com/andersfylling/disgord"
	"github.com/andersfylling/disgord/json"
	"github.com/jonas747/dca"
	"github.com/kkdai/youtube/v2"
	"github.com/zackartz/cmdlr2"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

func CreateNameToGet(s string) string {
	sm := strings.Split(s, " ")
	name := ""
	for _, add := range sm {
		name += add + "%20"
	}
	name = strings.TrimSuffix(name, "%20")
	return name
}

func SkipMusic(guild *disgord.Guild, authorID disgord.Snowflake) {
	for _, vx := range guild.VoiceStates {
		if vx.UserID == authorID {
			ServersVC[guild.ID].MusicQueue[0].Skip = make(chan bool)
			ServersVC[guild.ID].MusicQueue[0].Skip <- true
		}
	}
}

func StopStream(guild *disgord.Guild, authorID disgord.Snowflake) {
	for _, vx := range guild.VoiceStates {
		if vx.UserID == authorID {
			ServersVC[guild.ID].MusicQueue[0].Stop = true
		}
	}
}

func PauseStream(guild *disgord.Guild, authorID disgord.Snowflake) {
	ServersVC[guild.ID].MusicQueue[0].Pause <- true
}

func UnPauseStream(guild *disgord.Guild, authorID disgord.Snowflake) {
	ServersVC[guild.ID].MusicQueue[0].Pause <- false
}

func GetVideoFromName(name string) (string, bool) {
	nameToGet := CreateNameToGet(name)

	get := fmt.Sprintf("https://www.googleapis.com/youtube/v3/search?part=snippet&maxResults=1&q=%s&type=video&key="+os.Getenv("YT_KEY"), nameToGet)
	fmt.Println(get)
	r, err := http.Get(get)
	if err != nil {
		return get, false
	}

	resp_iotil, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return name, false
	}
	var body map[string]interface{}
	err = json.Unmarshal(resp_iotil, &body)
	if len(body["items"].([]interface{})) == 0 {
		return name, false
	}
	items := body["items"].([]interface{})[0]
	items_info := items.(map[string]interface{})
	id2 := items_info["id"].(map[string]interface{})
	video_id := id2["videoId"]
	return video_id.(string), true
}

func UrlOrNot(s string) (string, bool) {
	if len(s) < 15 {
		id, status := GetVideoFromName(s)
		return "https://www.youtube.com/watch?v=" + id, status
	} else if s[0:14] == "www.youtube.com" {
		return "https://" + s, true
	} else if len(s) >= 24 {
		if s[0:23] == "https://www.youtube.com/" {
			return s, true
		} else {
			id, status := GetVideoFromName(s)
			return "https://www.youtube.com/watch?v=" + id, status
		}
	} else {
		id, status := GetVideoFromName(s)
		return "https://www.youtube.com/watch?v=" + id, status
	}
}

func StreamToDiscord(ctx *cmdlr2.Ctx, url string) {
	channel, err := ctx.Client.Channel(ctx.Event.Message.ChannelID).Get()
	if err != nil {
		return
	}
	guild, err := ctx.Client.Guild(ctx.Event.Message.GuildID).Get()
	if err != nil {
		return
	}

	for _, v := range guild.VoiceStates {
		if v.UserID == ctx.Event.Message.Author.ID {
			mi := new(MusicInfo)
			url, vid := GetUrl(url)
			mi.SongInfo = vid
			mi.Url = url
			mi.GuildID = ctx.Event.Message.GuildID
			mi.ChannelID = v.ChannelID
			mi.Ctx = ctx
			mi.TextChannelID = channel.ID
			mi.Pause = make(chan bool, 1)
			QueueChan <- mi
		}
	}

}

func QueueWay() {
	for {
		mi := <-QueueChan
		if mi == nil {
			continue
		}
		if ServersVC[mi.GuildID] == nil {
			vc := GetVcConnection(mi.Ctx, mi.GuildID, mi.ChannelID)
			ServersVC[mi.GuildID] = new(ServerQueue)
			ServersVC[mi.GuildID].vc = vc
			ServersVC[mi.GuildID].MusicQueue = append(ServersVC[mi.GuildID].MusicQueue, mi)
			go playOnSever(mi.GuildID)
		} else {
			_, _ = ServersVC[mi.GuildID].MusicQueue[0].Ctx.Client.Channel(ServersVC[mi.GuildID].MusicQueue[0].TextChannelID).CreateMessage(&disgord.CreateMessageParams{
				Embed: &disgord.Embed{
					Title:       mi.SongInfo.Title,
					Description: fmt.Sprintf("Added **%s** to the queue!", mi.SongInfo.Title),
					Color:       0xFF0000,
					Author: &disgord.EmbedAuthor{
						Name: mi.SongInfo.Author,
					},
				},
			})
			ServersVC[mi.GuildID].MusicQueue = append(ServersVC[mi.GuildID].MusicQueue, mi)
		}
	}
}

func playOnSever(guildID disgord.Snowflake) {
	for {
		if len(ServersVC[guildID].MusicQueue) == 0 {
			break
		}
		ServersVC[guildID].MusicQueue[0].Ctx.Client.Channel(ServersVC[guildID].MusicQueue[0].TextChannelID).CreateMessage(&disgord.CreateMessageParams{
			Embed: &disgord.Embed{
				Title:       ServersVC[guildID].MusicQueue[0].SongInfo.Title,
				Description: fmt.Sprintf("Now playing **%s**!", ServersVC[guildID].MusicQueue[0].SongInfo.Title),
				Color:       0xFF0000,
				Author: &disgord.EmbedAuthor{
					Name: ServersVC[guildID].MusicQueue[0].SongInfo.Author,
				},
			},
		})
		stopStatus := Play(ServersVC[guildID].MusicQueue[0].GuildID, ServersVC[guildID].MusicQueue[0].ChannelID, ServersVC[guildID].MusicQueue[0].Url, ServersVC[guildID].vc)
		if stopStatus == true {
			os.Remove(ServersVC[guildID].MusicQueue[0].Url)
			break
		}
		os.Remove(ServersVC[guildID].MusicQueue[0].Url)
		ServersVC[guildID].MusicQueue = ServersVC[guildID].MusicQueue[1:]
	}
	if ServersVC[guildID].vc == nil {
		return
	}
	ServersVC[guildID].vc.Close()
	delete(ServersVC, guildID)
}

func GetUrl(videoURl string) (string, *youtube.Video) {
	client := youtube.Client{
		Debug:      false,
		HTTPClient: http.DefaultClient,
	}
	vid, err := client.GetVideo(videoURl)
	if err != nil {
		return "", nil
	}
	format := vid.Formats.FindByItag(251)
	url, _ := client.GetStreamURL(vid, format)

	os.Mkdir("tmp", 0755)
	filename := fmt.Sprintf("./tmp/" + RandStringRunes(10))
	fmt.Println("Downloading ", vid.Title, " to ", filename)

	go func() {
		resp, err := http.Get(url)
		if err != nil {
			return
		}
		defer resp.Body.Close()

		f, err := os.Create(filename)
		if err != nil {
			return
		}
		defer f.Close()

		_, err = io.Copy(f, resp.Body)
		if err != nil {
			return
		}
	}()

	return filename, vid
}

func Play(guildID, channelID disgord.Snowflake, url string, vc disgord.VoiceConnection) bool {
	options := dca.StdEncodeOptions
	options.RawOutput = true
	options.Bitrate = 96
	options.Application = "lowdelay"
	options.Volume = 50
	encode, err := dca.EncodeFile(url, options)
	if err != nil {
		fmt.Printf("%v 3\n", err)
	}
	if vc == nil {
		return false
	}
	vc.StartSpeaking()
	defer encode.Cleanup()
	for {
		select {
		case <-ServersVC[guildID].MusicQueue[0].Pause:
			<-ServersVC[guildID].MusicQueue[0].Pause
		case <-ServersVC[guildID].MusicQueue[0].Skip:
			return false
		default:
			if ServersVC[guildID].MusicQueue[0].Stop == true {
				encode.Cleanup()
				return true
			}
			var sz_frame int16
			err := binary.Read(encode, binary.LittleEndian, &sz_frame)
			if err != nil {
				fmt.Println(err)
				return false
			}
			Inbuf := make([]byte, sz_frame)
			_ = binary.Read(encode, binary.LittleEndian, &Inbuf)
			err = vc.SendOpusFrame(Inbuf)
			if err != nil {
				fmt.Printf("%v 2\n", err)
			}
		}
	}
	return false
}

func GetVcConnection(ctx *cmdlr2.Ctx, guildID, channelID disgord.Snowflake) disgord.VoiceConnection {
	voice, _ := ctx.Client.Guild(guildID).VoiceChannel(channelID).Connect(false, true)
	return voice
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
