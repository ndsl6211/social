package get_user

import (
	"github.com/google/uuid"
	"mashu.example/internal/usecase"
	"mashu.example/internal/usecase/repository"
)

type GetUserUseCaseReq struct {
	userId string
}

type GetUserUseCaseRes struct {
	ID          uuid.UUID
	UserName    string
	DisplayName string
	Email       string
	Public      bool

	Err error
}

type GetUserUseCase struct {
	userRepo repository.UserRepo

	Req *GetUserUseCaseReq
	Res *GetUserUseCaseRes
}

func (uc *GetUserUseCase) Execute() {
	userId := uuid.MustParse(uc.Req.userId)

	user, err := uc.userRepo.GetUserById(userId)
	if err != nil {
		uc.Res.Err = err
		return
	}

	uc.Res.ID = user.ID
	uc.Res.UserName = user.UserName
	uc.Res.DisplayName = user.DisplayName
	uc.Res.Email = user.Email
	uc.Res.Public = user.Public
	uc.Res.Err = nil
}

func NewGetUserUseCase(
	userRepo repository.UserRepo,
	req *GetUserUseCaseReq,
	res *GetUserUseCaseRes,
) usecase.UseCase {
	return &GetUserUseCase{userRepo, req, res}
}

func NewGetUserUseCaseReq(userId string) GetUserUseCaseReq {
	return GetUserUseCaseReq{userId}
}

func NewGetUserUseCaseRes() GetUserUseCaseRes {
	return GetUserUseCaseRes{}
}
