package discord

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"mashu.example/internal/usecase/repository"
)

type commandHandler func(string, []string, *discordgo.Session, *discordgo.MessageCreate)
type replyHandler func(key, cmd, channelId string, data map[string]string, reply string, s *discordgo.Session)

type botMessageHandler struct {
	userRepo  repository.UserRepo
	postRepo  repository.PostRepo
	groupRepo repository.GroupRepo
	dcRedis   *redis.Client

	botSess *discordgo.Session

	cmdHandlerMap   map[string]commandHandler
	replyHandlerMap map[string]replyHandler
}

func (h *botMessageHandler) HandleCommand(s *discordgo.Session, e *discordgo.MessageCreate) {
	msg := e.Content
	msgArr := strings.Split(msg, " ")
	cmd, params := msgArr[0][1:], msgArr[1:]
	logger.Infof("command: %s", cmd)
	logger.Infof("params: %v", params)

	handlerFunc, ok := h.cmdHandlerMap[cmd]
	if ok {
		handlerFunc(cmd, params, s, e)
	} else {
		s.ChannelMessageSendReply(e.ChannelID, "unknown command", e.Reference())
	}
}

func (h *botMessageHandler) register(cmd string, params []string, s *discordgo.Session, e *discordgo.MessageCreate) {
	logger.Info("register command received!")
	ctx := context.Background()

	ch, err := s.MessageThreadStart(
		e.ChannelID,
		e.ID,
		fmt.Sprintf("registration thread for %s%s", e.Author.Username, e.Author.Discriminator),
		60,
	)
	if err != nil {
		logger.Error("failed to create thread", err)
		return
	}

	h.dcRedis.HSet(
		ctx,
		h.getRedisCmdSessKey(e.Author.Username, e.Author.Discriminator, ch.ID, cmd),
		map[string]interface{}{"username": ""},
	)

	if _, err := s.ChannelMessageSend(ch.ID, "請輸入帳號"); err != nil {
		logger.Error("failed to send reply message", err)
		return
	}
}

func (h *botMessageHandler) login(cmd string, params []string, s *discordgo.Session, e *discordgo.MessageCreate) {
	logger.Info("login command received!")
	ctx := context.Background()
	dcUserId := fmt.Sprintf("%s%s", e.Author.Username, e.Author.Discriminator)

	ch, err := s.MessageThreadStart(
		e.ChannelID,
		e.ID,
		fmt.Sprintf("registration thread for %s", dcUserId),
		60,
	)
	if err != nil {
		logrus.Error("failed to create thread: ", err)
		return
	}

	if _, isLogin := h.checkIsLogin(dcUserId, s, ch); isLogin {
		s.ChannelMessageSend(ch.ID, "哥, 你已經登入ㄌ")
		return
	}

	if _, err := h.dcRedis.HSet(
		ctx,
		h.getRedisCmdSessKey(e.Author.Username, e.Author.Discriminator, ch.ID, cmd),
		map[string]interface{}{"username": ""},
	).Result(); err != nil {
		logrus.Error("failed to set command session:", err)
		return
	}

	if _, err := s.ChannelMessageSend(ch.ID, "請輸入帳號"); err != nil {
		logrus.Error("failed to send reply message", err)
		return
	}
}

func (h *botMessageHandler) logout(cmd string, params []string, s *discordgo.Session, e *discordgo.MessageCreate) {
	logger.Info("logout command received!")
	dcUserId := fmt.Sprintf("%s%s", e.Author.Username, e.Author.Discriminator)

	ch, err := s.MessageThreadStart(
		e.ChannelID,
		e.ID,
		fmt.Sprintf("registration thread for %s", dcUserId),
		60,
	)
	if err != nil {
		logrus.Error("failed to create thread: ", err)
		return
	}

	if _, isLogin := h.checkIsLogin(dcUserId, s, ch); !isLogin {
		return
	}

	if r, err := h.dcRedis.Del(
		context.Background(),
		h.getRedisLoginSessKey(fmt.Sprintf("%s%s", e.Author.Username, e.Author.Discriminator)),
	).Result(); err != nil {
		logrus.Error(err)
		s.ChannelMessageSend(ch.ID, "登出失敗, 好像有哪裡出錯ㄌ")
	} else {
		// if the user is not logged in
		if r == 0 {
			logger.Info("the user is not logged in")
			s.ChannelMessageSend(ch.ID, "哥, 你還沒登入")
		} else {
			s.ChannelMessageSend(ch.ID, "登出成功, 若要繼續使用請再次登入")
		}
	}

	s.ChannelEditComplex(ch.ID, &discordgo.ChannelEdit{
		Archived: true,
		Locked:   true,
	})
}

func (h *botMessageHandler) followUser(cmd string, params []string, s *discordgo.Session, e *discordgo.MessageCreate) {
	logger.Info("follow user command received!")
	// ctx := context.Background()

	// ch, err := s.MessageThreadStart(
	// 	e.ChannelID,
	// 	e.ID,
	// 	fmt.Sprintf("follow user [%s%s]", e.Author.Username, e.Author.Discriminator),
	// 	60,
	// )
	// if err != nil {
	// 	logrus.Error("failed to create thread: ", err)
	// 	return
	// }

	// if _, err := h.dcRedis.HSet(
	// 	ctx,
	// 	h.getRedisCmdSessKey(e.Author.Username, e.Author.Discriminator, ch.ID, cmd),
	// 	map[string]interface{}{""},
	// )
}

func (h *botMessageHandler) getRedisCmdSessKey(userName, userNum, channelId, cmd string) string {
	return fmt.Sprintf("sess:user:%s%s:%s:%s", userName, userNum, channelId, cmd)
}

func (h *botMessageHandler) getRedisCmdSessKeyForSearch(userName, userNum, channelId string) string {
	return fmt.Sprintf("sess:user:%s%s:%s:*", userName, userNum, channelId)
}

// get the redis key for storing login status
//
// dcUserId = username + 4-digit number
func (h *botMessageHandler) getRedisLoginSessKey(dcUserId string) string {
	return fmt.Sprintf("sess:discord:login:%s", dcUserId)
}

func (h *botMessageHandler) checkIsLogin(dcUserId string, s *discordgo.Session, ch *discordgo.Channel) (string, bool) {
	userId, err := h.dcRedis.Get(
		context.Background(),
		h.getRedisLoginSessKey(dcUserId),
	).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			logger.Info("the user is not logged in")
			s.ChannelMessageSend(ch.ID, "哥, 你還沒登入")
		} else {
			s.ChannelMessageSend(ch.ID, "好像有哪裡出錯了")
		}
		return "", false
	}

	s.ChannelEditComplex(ch.ID, &discordgo.ChannelEdit{
		Archived: true,
		Locked:   true,
	})
	return userId, true
}
