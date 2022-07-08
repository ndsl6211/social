package discord

import (
	"context"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"mashu.example/internal/usecase/user/login"
	"mashu.example/internal/usecase/user/register"
)

func (h *botMessageHandler) HandleReply(s *discordgo.Session, e *discordgo.MessageCreate) {
	userName := e.Author.Username
	userNum := e.Author.Discriminator
	channelId := e.Message.ChannelID

	// try to find channel id based on user id
	ctx := context.Background()
	var cursor uint64
	keys, cursor, err := h.dcRedis.Scan(
		ctx,
		cursor,
		h.getRedisCmdSessKeyForSearch(userName, userNum, channelId),
		0,
	).Result()
	if err != nil {
		logger.Error(err)
		return
	}

	// no active session
	if len(keys) == 0 {
		logger.Infof("there is no active session for user %s%s", userName, userNum)
		return
	}

	key := keys[0]
	keyArr := strings.Split(key, ":")
	cmd := keyArr[len(keyArr)-1]
	info, err := h.dcRedis.HGetAll(ctx, key).Result()
	if err != nil {
		logrus.Error(err)
		return
	}

	replyHandlerFunc, ok := h.replyHandlerMap[cmd]
	if ok {
		replyHandlerFunc(key, cmd, channelId, info, e.Content, s)
	} else {
		logrus.Errorf("unknown command reply: %s", cmd)
		return
	}
}

func (h *botMessageHandler) handleRegisterReply(
	key string,
	cmd string,
	channelId string,
	data map[string]string,
	reply string,
	s *discordgo.Session,
) {
	stage := len(data)
	ctx := context.Background()
	stageKeys := []string{"username", "displayName", "email"}
	stageMessages := []string{"", "請輸入使用者名稱", "請輸入email"}

	// save reply
	if _, err := h.dcRedis.HSet(ctx, key, map[string]string{stageKeys[stage-1]: reply}).Result(); err != nil {
		logrus.Error(err)
		return
	}

	if stage < len(stageKeys) {
		// send next inquiry
		s.ChannelMessageSend(channelId, stageMessages[stage])

		// save placeholder for next inquiry
		if _, err := h.dcRedis.HSet(ctx, key, map[string]string{stageKeys[stage]: ""}).Result(); err != nil {
			logrus.Error(err)
			return
		}
		return
	}

	completeData, err := h.dcRedis.HGetAll(ctx, key).Result()
	if err != nil {
		logrus.Error(err)
		return
	}

	if _, err := h.dcRedis.Del(ctx, key).Result(); err != nil {
		logrus.Error("failed to remove session from redis: ", err)
		return
	}

	req := register.NewRegisterUseCaseReq(
		completeData["username"],
		completeData["displayName"],
		completeData["email"],
	)
	res := register.NewRegisterUseCaseRes()
	uc := register.NewRegisterUseCase(h.userRepo, req, res)
	uc.Execute()
	if res.Err != nil {
		logrus.Error("failed to run register usecase: ", res.Err)
		s.ChannelMessageSend(channelId, "註冊失敗, 好像有哪裡出錯ㄌ")
		return
	}

	s.ChannelMessageSend(channelId, "註冊完成")
	s.ChannelEditComplex(channelId, &discordgo.ChannelEdit{
		Archived: true,
		Locked:   true,
	})
}

func (h *botMessageHandler) handleLoginReply(
	key string,
	cmd string,
	channelId string,
	data map[string]string,
	reply string,
	s *discordgo.Session,
) {
	stage := len(data)
	ctx := context.Background()
	stageKeys := []string{"username"}
	stageMessages := []string{""}

	// save reply
	if _, err := h.dcRedis.HSet(ctx, key, map[string]string{stageKeys[stage-1]: reply}).Result(); err != nil {
		logrus.Error(err)
		return
	}

	if stage < len(stageKeys) {
		// send next inquiry
		s.ChannelMessageSend(channelId, stageMessages[stage])

		// save placeholder for next inquiry
		if _, err := h.dcRedis.HSet(ctx, key, map[string]string{stageKeys[stage]: ""}).Result(); err != nil {
			logrus.Error(err)
			return
		}
		return
	}

	completeData, err := h.dcRedis.HGetAll(ctx, key).Result()
	if err != nil {
		logrus.Error(err)
		return
	}

	if _, err := h.dcRedis.Del(ctx, key).Result(); err != nil {
		logrus.Error("failed to remove session from redis: ", err)
		return
	}

	req := login.NewLoginUseCaseReq(completeData["username"])
	res := login.NewLoginUseCaseRes()
	uc := login.NewLoginUseCase(h.userRepo, req, res)
	uc.Execute()

	if res.Err != nil {
		logrus.Error("failed to run register usecase: ", res.Err)
		s.ChannelMessageSend(channelId, "登入失敗, 好像有哪裡出錯ㄌ")
		return
	}

	splitKey := strings.Split(key, ":")
	if _, err := h.dcRedis.Set(
		ctx,
		h.getRedisLoginSessKey(splitKey[2]),
		reply,
		0,
	).Result(); err != nil {
		logrus.Error("failed to save login status: ", err)
		s.ChannelMessageSend(channelId, "登入失敗, 好像有哪裡出錯ㄌ")
		return
	}

	s.ChannelMessageSend(channelId, "登入成功")
	s.ChannelEditComplex(channelId, &discordgo.ChannelEdit{
		Archived: true,
		Locked:   true,
	})
}
