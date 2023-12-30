package discord

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/redis/go-redis/v9"

	"mashu.example/config"
	"mashu.example/internal/usecase/repository"
	"mashu.example/pkg"
)

var (
	logger    = pkg.NewScopedLogger("DISCORD")
	cmdPrefix = config.DiscordBotPrefix
)

type DiscordBot struct {
	handler discordCommandHandler
	botSess *discordgo.Session
}

func NewDiscordBot(
	userRepo repository.UserRepo,
	postRepo repository.PostRepo,
	groupRepo repository.GroupRepo,
	dcRedis *redis.Client,
) (*DiscordBot, error) {

	botSess, err := discordgo.New("Bot " + config.DiscordBotToken)
	if err != nil {
		logger.Error("failed to new discord bot")
		return nil, err
	}

	discordBot := &DiscordBot{
		handler: discordCommandHandler{
			userRepo:      userRepo,
			postRepo:      postRepo,
			groupRepo:     groupRepo,
			dcRedis:       dcRedis,
			botSess:       botSess,
			cmdHandlerMap: map[string]commandHandler{},
		},
		botSess: botSess,
	}

	botSess.AddHandler(discordBot.messageHandler)
	botSess.AddHandler(discordBot.InteractionCreateHandler)

	return discordBot, nil
}

func (b *DiscordBot) registerDiscordBotCommandHandler() {
	b.handler.cmdHandlerMap["register"] = b.handler.register
	b.handler.cmdHandlerMap["login"] = b.handler.login
	b.handler.cmdHandlerMap["logout"] = b.handler.logout
	b.handler.cmdHandlerMap["createPost"] = b.handler.createPost
}

func (b *DiscordBot) Start() {
	logger.Info("discord bot started")

	if err := b.botSess.Open(); err != nil {
		logger.Error("failed to open websocket connection")
		return
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	b.botSess.Close()
}

func (b *DiscordBot) messageHandler(s *discordgo.Session, e *discordgo.MessageCreate) {
	logger.
		WithField("sender_name", e.Author.Username).
		WithField("sender_ID", e.Author.ID).
		WithField("sender_number", e.Author.Discriminator).
		WithField("channel_ID", e.Message.ChannelID).
		Info(e.Content)

	// ignore all message sent by bot
	if e.Author.Bot {
		return
	}

	if len(e.Content) == 0 {
		return
	}
}

func (b *DiscordBot) InteractionCreateHandler(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
) {
	if h, ok := b.handler.cmdHandlerMap[i.ApplicationCommandData().Name]; ok {
		optionMap := make(
			map[string]*discordgo.ApplicationCommandInteractionDataOption,
			len(i.ApplicationCommandData().Options),
		)
		for _, opt := range i.ApplicationCommandData().Options {
			optionMap[opt.Name] = opt
		}

		h(s, i, optionMap)
	}
}
