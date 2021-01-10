package main

import (
	"github.com/andersfylling/disgord"
	"github.com/sirupsen/logrus"
	"github.com/zackartz/cmdlr2"
	"os"
	"synergy/cmds"
	"synergy/db"
	"synergy/events"
)

var log = &logrus.Logger{
	Out:       os.Stderr,
	Formatter: new(logrus.TextFormatter),
	Hooks:     make(logrus.LevelHooks),
	Level:     logrus.DebugLevel,
}

func main() {
	db.Init()

	// Set up a new Disgord client
	client := disgord.New(disgord.Config{
		BotToken: os.Getenv("DISCORD_TOKEN"),
		Logger:   log,
	})

	router := cmdlr2.Create(&cmdlr2.Router{
		Prefixes:         []string{"$"},
		Client:           client,
		BotsAllowed:      false,
		IgnorePrefixCase: true,
	})

	router.RegisterCMDList(cmds.CommandList)

	router.RegisterDefaultHelpCommand(client)

	router.Initialize(client)

	client.Gateway().MessageReactionAdd(events.EmojiAdd)
	client.Gateway().MessageReactionRemove(events.EmojiRemove)

	defer client.Gateway().StayConnectedUntilInterrupted()
}
