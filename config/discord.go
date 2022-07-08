package config

import (
	"os"

	"github.com/sirupsen/logrus"
)

var (
	DiscordBotToken  string
	DiscordBotPrefix string
)

func init() {
	logrus.Info()
	DiscordBotToken = os.Getenv("DISCORD_BOT_TOKEN")
	DiscordBotPrefix = "!"
}
