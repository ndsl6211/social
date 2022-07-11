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

// type of handler function for the command sent from user
//
// - cmd: the command sent from user
// - params: the parameters of the command
// - channelId: the channel id of the created thread
// - dcUserId: the id of user in discord (username + 4-digits number)
// - s: the active chatting session of discord
// - e: the message event of discord
type commandHandler func(cmd string, params []string, channelId string, dcUserId string, s *discordgo.Session, e *discordgo.MessageCreate)

// type of handler function for the reply from user in thread
//
// - key: the active session key get from redis
// - channelId: the channel id of the active discord thread
// - dcUserId: the id of user in discord (username + 4-digits number)
// - data: data get from redis (all the replies from user)
// - reply: the latest reply from user
// - s: the active chatting session of discord
type replyHandler func(activeSessKey, channelId, dcUserId string, data map[string]string, reply string, s *discordgo.Session)

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
	dcUserId := fmt.Sprintf("%s%s", e.Author.Username, e.Author.Discriminator)

	logger.Infof("command: %s", cmd)
	logger.Infof("params: %v", params)

	handlerFunc, ok := h.cmdHandlerMap[cmd]
	if ok {
		ch, err := s.MessageThreadStart(
			e.ChannelID,
			e.ID,
			fmt.Sprintf("Thread for %s", dcUserId),
			60,
		)
		if err != nil {
			logger.Error("failed to create thread", err)
			return
		}
		handlerFunc(cmd, params, ch.ID, dcUserId, s, e)
	} else {
		s.ChannelMessageSendReply(e.ChannelID, "unknown command", e.Reference())
	}
}

func (h *botMessageHandler) register(cmd string, params []string, channelId string, dcUserId string, s *discordgo.Session, e *discordgo.MessageCreate) {
	ctx := context.Background()

	h.dcRedis.HSet(
		ctx,
		h.getRedisCmdSessKey(dcUserId, channelId, cmd),
		map[string]interface{}{"username": ""},
	)

	if _, err := s.ChannelMessageSend(channelId, "請輸入帳號"); err != nil {
		logger.Error("failed to send reply message", err)
		return
	}
}

func (h *botMessageHandler) login(cmd string, params []string, channelId string, dcUserId string, s *discordgo.Session, e *discordgo.MessageCreate) {
	ctx := context.Background()

	if _, isLogin := h.checkIsLogin(dcUserId, channelId, s, false); isLogin {
		s.ChannelMessageSend(channelId, "哥, 你已經登入ㄌ")
		s.ChannelEditComplex(channelId, &discordgo.ChannelEdit{
			Archived: true,
			Locked:   true,
		})
		return
	}

	if _, err := h.dcRedis.HSet(
		ctx,
		h.getRedisCmdSessKey(dcUserId, channelId, cmd),
		map[string]interface{}{"username": ""},
	).Result(); err != nil {
		logrus.Error("failed to set command session:", err)
		return
	}

	if _, err := s.ChannelMessageSend(channelId, "請輸入帳號"); err != nil {
		logrus.Error("failed to send reply message", err)
		return
	}
}

func (h *botMessageHandler) logout(cmd string, params []string, channelId string, dcUserId string, s *discordgo.Session, e *discordgo.MessageCreate) {
	if _, isLogin := h.checkIsLogin(dcUserId, channelId, s, true); !isLogin {
		return
	}

	if _, err := h.dcRedis.Del(
		context.Background(),
		h.getRedisLoginSessKey(dcUserId),
	).Result(); err != nil {
		logrus.Error(err)
		s.ChannelMessageSend(channelId, "登出失敗, 好像有哪裡出錯ㄌ")
	} else {
		s.ChannelMessageSend(channelId, "登出成功, 若要繼續使用請再次登入")
	}

	s.ChannelEditComplex(channelId, &discordgo.ChannelEdit{
		Archived: true,
		Locked:   true,
	})
}

func (h *botMessageHandler) followUser(cmd string, params []string, channelId string, dcUserId string, s *discordgo.Session, e *discordgo.MessageCreate) {
	// ctx := context.Background()

	// if _, err := h.dcRedis.HSet(
	// 	ctx,
	// 	h.getRedisCmdSessKey(dcUserId, channelId, cmd),
	// 	map[string]interface{}{"followeeId": ""},
	// ).Result(); err != nil {
	// 	logrus.Error("failed to set command session:", err)
	// 	return
	// }

	// if _, err := s.ChannelMessageSend(channelId, "請輸入你要追蹤的使用者"); err != nil {
	// 	logrus.Error("failed to send reply message", err)
	// 	return
	// }
}

func (h *botMessageHandler) getRedisCmdSessKey(dcUserId, channelId, cmd string) string {
	return fmt.Sprintf("sess:user:%s:%s:%s", dcUserId, channelId, cmd)
}

func (h *botMessageHandler) getRedisCmdSessKeyForSearch(dcUserId, channelId string) string {
	return fmt.Sprintf("sess:user:%s:%s:*", dcUserId, channelId)
}

// get the redis key for storing login status
//
// dcUserId = username + 4-digit number
func (h *botMessageHandler) getRedisLoginSessKey(dcUserId string) string {
	return fmt.Sprintf("sess:discord:login:%s", dcUserId)
}

func (h *botMessageHandler) checkIsLogin(dcUserId string, channelId string, s *discordgo.Session, sendReply bool) (string, bool) {
	userId, err := h.dcRedis.Get(
		context.Background(),
		h.getRedisLoginSessKey(dcUserId),
	).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			logger.Info("the user is not logged in")
			if sendReply {
				s.ChannelMessageSend(channelId, "哥, 你還沒登入")
			}
		} else {
			s.ChannelMessageSend(channelId, "好像有哪裡出錯了")
		}
		s.ChannelEditComplex(channelId, &discordgo.ChannelEdit{
			Archived: true,
			Locked:   true,
		})
		return "", false
	}

	return userId, true
}
