package bind_discord_user

import (
	"errors"

	"mashu.example/internal/usecase/repository"
)

type BindDCUserUseCaseReq struct {
	discordUserId string
	userName      string
}

type BindDCUserUseCaseRes struct {
	Err error
}

type BindDCUserUseCase struct {
	userRepo repository.UserRepo
	req      *BindDCUserUseCaseReq
	res      *BindDCUserUseCaseRes
}

func (uc *BindDCUserUseCase) Execute() {
	user, err := uc.userRepo.GetUserByUserName(uc.req.userName)
	if err != nil {
		err = errors.New("user doesn't exist")
		uc.res.Err = err
	}

	user.DiscordUserId = uc.req.discordUserId

	if err := uc.userRepo.Save(user); err != nil {
		uc.res.Err = err
		return
	}
}

func NewBindDCUserUseCase(
	userRepo repository.UserRepo,
	req *BindDCUserUseCaseReq,
	res *BindDCUserUseCaseRes,
) *BindDCUserUseCase {
	return &BindDCUserUseCase{userRepo, req, res}
}

func NewBindDCUserUseCaseReq(discordUserId, userName string) *BindDCUserUseCaseReq {
	return &BindDCUserUseCaseReq{discordUserId, userName}
}

func NewBindDCUserUseCaseRes() *BindDCUserUseCaseRes {
	return &BindDCUserUseCaseRes{}
}
