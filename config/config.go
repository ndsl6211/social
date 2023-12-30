package config

import (
	"strings"

	"github.com/spf13/viper"
)

type SocialServerConfig struct {
	Discord struct {
		BotToken string
	}
}

func SetupConfig() SocialServerConfig {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	config := SocialServerConfig{
		Discord: struct {
			BotToken string
		}{
			BotToken: viper.GetString("discord.bot.token"),
		},
	}

	return config
}
