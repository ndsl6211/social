package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"mashu.example/internal/usecase/repository"
	"mashu.example/internal/usecase/user/login"
	"mashu.example/internal/usecase/user/logout"
	"mashu.example/internal/usecase/user/register"
)

type commandHandler func(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	opt map[string]*discordgo.ApplicationCommandInteractionDataOption,
)

type discordCommandHandler struct {
	userRepo  repository.UserRepo
	postRepo  repository.PostRepo
	groupRepo repository.GroupRepo
	loginRepo repository.LoginRepo
	dcRedis   *redis.Client

	botSess *discordgo.Session

	cmdHandlerMap map[string]commandHandler
}

func (h *discordCommandHandler) responseMessage(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	message string,
) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})
}

func (h *discordCommandHandler) register(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	opt map[string]*discordgo.ApplicationCommandInteractionDataOption,
) {
	req := register.NewRegisterUseCaseReq(
		opt["username"].StringValue(),
		opt["displayName"].StringValue(),
		opt["displayName"].StringValue(),
	)
	res := register.NewRegisterUseCaseRes()
	uc := register.NewRegisterUseCase(h.userRepo, req, res)
	uc.Execute()

	if res.Err != nil {
		logrus.Error("failed to run register usecase: ", res.Err)
		h.responseMessage(s, i, "註冊失敗, 好像有哪裡出錯ㄌ")
		return
	}

	h.responseMessage(s, i, "註冊成功")
}

func (h *discordCommandHandler) login(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	opt map[string]*discordgo.ApplicationCommandInteractionDataOption,
) {
	req := login.NewLoginUseCaseReq(opt["username"].StringValue())
	res := login.NewLoginUseCaseRes()
	uc := login.NewLoginUseCase(h.userRepo, req, res)
	uc.Execute()

	if res.Err != nil {
		logrus.Error("failed to run login usecase: ", res.Err)
		h.responseMessage(s, i, "登入失敗, 好像有哪裡出錯ㄌ")
		return
	}

	h.responseMessage(s, i, "登入成功")
}

func (h *discordCommandHandler) logout(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	opt map[string]*discordgo.ApplicationCommandInteractionDataOption,
) {
	req := logout.LogoutUseCaseReq{
		UserName: opt["username"].StringValue(),
	}
	res := logout.LogoutUseCaseRes{}
	uc := logout.NewLogoutUseCase(h.userRepo, h.loginRepo, &req, &res)
	uc.Execute()

	if res.Err != nil {
		h.responseMessage(s, i, res.Err.Error())
	}

	h.responseMessage(s, i, "登出成功, 若要繼續使用請再次登入")
}

func (h *discordCommandHandler) createPost(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	opt map[string]*discordgo.ApplicationCommandInteractionDataOption,
) {
	// ctx := context.Background()

	// req := create_post.CreatePostUseCaseReq{
	// 	Title:      opt["title"].StringValue(),
	// 	Content:    opt["content"].StringValue(),
	// 	Permission: entity_enums.PostPermission(opt["permission"].StringValue()),
	// }
}

//func (h *discordCommandHandler) checkIsLogin(dcUserId string, channelId string, s *discordgo.Session, sendReply bool) (uuid.UUID, string, bool) {
//userData, err := h.dcRedis.HGetAll(
//context.Background(),
//h.getRedisLoginSessKey(dcUserId),
//).Result()
//if err != nil {
//s.ChannelMessageSend(channelId, "好像有哪裡出錯了")
//s.ChannelEditComplex(channelId, &discordgo.ChannelEdit{
//Archived: true,
//Locked:   true,
//})
//return uuid.Nil, "", false
//}
//if len(userData) == 0 {
//logger.Info("the user is not logged in")
//if sendReply {
//s.ChannelMessageSend(channelId, "哥, 你還沒登入")
//}
//return uuid.Nil, "", false
//}

//return uuid.MustParse(userData["userId"]), userData["username"], true
//}
