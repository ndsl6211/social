package discord

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	entity_enums "mashu.example/internal/entity/enums"
	"mashu.example/internal/usecase/post/create_post"
	"mashu.example/internal/usecase/user/login"
	"mashu.example/internal/usecase/user/register"
)

func (h *botMessageHandler) HandleReply(s *discordgo.Session, e *discordgo.MessageCreate) {
	dcUserId := fmt.Sprintf("%s%s", e.Author.Username, e.Author.Discriminator)
	channelId := e.Message.ChannelID

	// try to find channel id based on user id
	ctx := context.Background()
	var cursor uint64
	keys, cursor, err := h.dcRedis.Scan(
		ctx,
		cursor,
		h.getRedisCmdSessKeyForSearch(dcUserId, channelId),
		0,
	).Result()
	if err != nil {
		logger.Error(err)
		return
	}

	// no active session
	if len(keys) == 0 {
		logger.Infof("there is no active session for dc user %s", dcUserId)
		return
	}

	activeSessKey := keys[0]
	keyArr := strings.Split(activeSessKey, ":")
	cmd := keyArr[len(keyArr)-1]
	prevReplyData, err := h.dcRedis.HGetAll(ctx, activeSessKey).Result()
	if err != nil {
		logrus.Error(err)
		return
	}

	replyHandlerFunc, ok := h.replyHandlerMap[cmd]
	if ok {
		replyHandlerFunc(activeSessKey, channelId, dcUserId, prevReplyData, e.Content, s)
	} else {
		logrus.Errorf("unknown command reply: %s", cmd)
		return
	}
}

func (h *botMessageHandler) handleRegisterReply(
	activeSessKey string,
	channelId string,
	dcUserId string,
	data map[string]string,
	reply string,
	s *discordgo.Session,
) {
	stage := len(data)
	ctx := context.Background()
	stageKeys := []string{"username", "displayName", "email"}
	stageMessages := []string{"", "請輸入使用者名稱", "請輸入email"}

	// save reply
	if _, err := h.dcRedis.HSet(ctx, activeSessKey, map[string]string{stageKeys[stage-1]: reply}).Result(); err != nil {
		logrus.Error(err)
		return
	}

	if stage < len(stageKeys) {
		// send next inquiry
		s.ChannelMessageSend(channelId, stageMessages[stage])

		// save placeholder for next inquiry
		if _, err := h.dcRedis.HSet(ctx, activeSessKey, map[string]string{stageKeys[stage]: ""}).Result(); err != nil {
			logrus.Error(err)
			return
		}
		return
	}

	completeData, err := h.dcRedis.HGetAll(ctx, activeSessKey).Result()
	if err != nil {
		logrus.Error(err)
		return
	}

	if _, err := h.dcRedis.Del(ctx, activeSessKey).Result(); err != nil {
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
	activeSessKey string,
	channelId string,
	dcUserId string,
	data map[string]string,
	reply string,
	s *discordgo.Session,
) {
	stage := len(data)
	ctx := context.Background()
	stageKeys := []string{"username"}
	stageMessages := []string{""}

	// save reply
	if _, err := h.dcRedis.HSet(ctx, activeSessKey, map[string]string{stageKeys[stage-1]: reply}).Result(); err != nil {
		logrus.Error(err)
		return
	}

	if stage < len(stageKeys) {
		// send next inquiry
		s.ChannelMessageSend(channelId, stageMessages[stage])

		// save placeholder for next inquiry
		if _, err := h.dcRedis.HSet(ctx, activeSessKey, map[string]string{stageKeys[stage]: ""}).Result(); err != nil {
			logrus.Error(err)
			return
		}
		return
	}

	completeData, err := h.dcRedis.HGetAll(ctx, activeSessKey).Result()
	if err != nil {
		logrus.Error(err)
		return
	}

	if _, err := h.dcRedis.Del(ctx, activeSessKey).Result(); err != nil {
		logrus.Error("failed to remove session from redis: ", err)
		return
	}

	req := login.NewLoginUseCaseReq(completeData["username"])
	res := login.NewLoginUseCaseRes()
	uc := login.NewLoginUseCase(h.userRepo, req, res)
	uc.Execute()

	if res.Err != nil {
		logrus.Error("failed to run login usecase: ", res.Err)
		s.ChannelMessageSend(channelId, "登入失敗, 好像有哪裡出錯ㄌ")
		s.ChannelEditComplex(channelId, &discordgo.ChannelEdit{
			Archived: true,
			Locked:   true,
		})
		return
	}

	user, _ := h.userRepo.GetUserByUserName(completeData["username"])

	if _, err := h.dcRedis.HSet(
		ctx,
		h.getRedisLoginSessKey(dcUserId),
		map[string]interface{}{"userId": user.ID.String(), "username": completeData["username"]},
	).Result(); err != nil {
		logrus.Error("failed to save login status: ", err)
		s.ChannelMessageSend(channelId, "登入失敗, 好像有哪裡出錯ㄌ")
		s.ChannelEditComplex(channelId, &discordgo.ChannelEdit{
			Archived: true,
			Locked:   true,
		})
		return
	}

	s.ChannelMessageSend(channelId, "登入成功")
	s.ChannelEditComplex(channelId, &discordgo.ChannelEdit{
		Archived: true,
		Locked:   true,
	})
}

func (h *botMessageHandler) handleCreatePostReply(
	activeSessKey string,
	channelId string,
	dcUserId string,
	data map[string]string,
	reply string,
	s *discordgo.Session,
) {
	userId, _, ok := h.checkIsLogin(dcUserId, channelId, s, true)
	if !ok {
		return
	}

	stage := len(data)
	ctx := context.Background()
	stageKeys := []string{"title", "content", "permission"}
	stageMessages := []string{"", "請輸入貼文內容", "請輸入該貼文的權限(0=公開, 1=僅限追蹤者, 2=私人)"}

	// save reply
	if _, err := h.dcRedis.HSet(ctx, activeSessKey, map[string]string{stageKeys[stage-1]: reply}).Result(); err != nil {
		logrus.Error(err)
		return
	}

	if stage < len(stageKeys) {
		// send next inquiry
		s.ChannelMessageSend(channelId, stageMessages[stage])

		// save placeholder for next inquiry
		if _, err := h.dcRedis.HSet(ctx, activeSessKey, map[string]string{stageKeys[stage]: ""}).Result(); err != nil {
			logrus.Error(err)
			return
		}
		return
	}

	completeData, err := h.dcRedis.HGetAll(ctx, activeSessKey).Result()
	if err != nil {
		logrus.Error(err)
		return
	}

	if _, err := h.dcRedis.Del(ctx, activeSessKey).Result(); err != nil {
		logrus.Error("failed to remove session from redis: ", err)
		return
	}

	permission, err := strconv.Atoi(completeData["permission"])
	if err != nil {
		logrus.Error("invalid permission:", permission)
		s.ChannelMessageSend(channelId, "無效的選擇")
		return
	}
	req := create_post.NewCreatePostUseCaseReq(
		completeData["title"],
		completeData["content"],
		userId,
		uuid.Nil,
		entity_enums.PostPermission(permission),
	)
	res := create_post.NewCreatePostUseCaseRes()
	uc := create_post.NewCreatePostUseCase(h.userRepo, h.postRepo, h.groupRepo, req, res)
	uc.Execute()

	if res.Err != nil {
		logrus.Error("failed to run create post usecase: ", res.Err)
		s.ChannelMessageSend(channelId, "創建貼文失敗, 好像有哪裡出錯ㄌ")
		s.ChannelEditComplex(channelId, &discordgo.ChannelEdit{
			Archived: true,
			Locked:   true,
		})
		return
	}

	s.ChannelMessageSend(channelId, "這是你的貼文")
	s.ChannelMessageSendEmbed(channelId, &discordgo.MessageEmbed{
		Title:       completeData["title"],
		Description: completeData["content"],
	})
	s.ChannelEditComplex(channelId, &discordgo.ChannelEdit{
		Archived: true,
		Locked:   true,
	})
}

func (h *botMessageHandler) handleFollowUserReply(
	activeSessKey string,
	channelId string,
	dcUserId string,
	data map[string]string,
	reply string,
	s *discordgo.Session,
) {
	// stage := len(data)
	// ctx := context.Background()
	// stageKeys := []string{"followeeId"}
	// stageMessages := []string{""}

	// // save reply
	// if _, err := h.dcRedis.HSet(ctx, activeSessKey, map[string]string{stageKeys[stage-1]: reply}).Result(); err != nil {
	// 	logrus.Error(err)
	// 	return
	// }

	// if stage < len(stageKeys) {
	// 	// send next inquiry
	// 	s.ChannelMessageSend(channelId, stageMessages[stage])

	// 	// save placeholder for next inquiry
	// 	if _, err := h.dcRedis.HSet(ctx, activeSessKey, map[string]string{stageKeys[stage]: ""}).Result(); err != nil {
	// 		logrus.Error(err)
	// 		return
	// 	}
	// 	return
	// }

	// completeData, err := h.dcRedis.HGetAll(ctx, activeSessKey).Result()
	// if err != nil {
	// 	logrus.Error(err)
	// 	return
	// }

	// if _, err := h.dcRedis.Del(ctx, activeSessKey).Result(); err != nil {
	// 	logrus.Error("failed to remove session from redis: ", err)
	// 	return
	// }

	// if err != nil {

	// }
	// res := follow_user.NewFollowUserUseCaseReq()
}
