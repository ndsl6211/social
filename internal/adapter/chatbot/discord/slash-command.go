package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "register",
			Description: "register",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "username",
					Description: "使用者帳號",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "email",
					Description: "電子信箱",
					Required:    true,
				},
			},
		},
		{
			Name:        "login",
			Description: "login",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "username",
					Description: "使用者帳號",
					Required:    true,
				},
			},
		},
		{
			Name:        "logout",
			Description: "logout",
		},
		{
			Name:        "create-post",
			Description: "創建貼文",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "title",
					Description: "標題",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "content",
					Description: "內文",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "permission",
					Description: "貼文權限",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "public",
							Value: "PUBLIC",
						},
						{
							Name:  "follower only",
							Value: "FOLLOWER_ONLY",
						},
						{
							Name:  "private",
							Value: "PRIVATE",
						},
					},
				},
			},
		},
	}

	registeredCommands = make([]*discordgo.ApplicationCommand, 0)
)

func RegisterSlashCommand(s *discordgo.Session) {
	logrus.Info("registering slash command...")

	for i, c := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, "", c)
		if err != nil {
			logrus.Panicf("failed to create command %v: %v", c.Name, err)
		}
		registeredCommands[i] = cmd
	}
}

func DeleteSlashCommand(s *discordgo.Session) {
	logrus.Info("removing slash commands...")

	for _, c := range registeredCommands {
		err := s.ApplicationCommandDelete(s.State.User.ID, "", c.ID)
		if err != nil {
			logrus.Panicf("failed to delete command %v: %v", c.Name, err)
		}
	}
}
