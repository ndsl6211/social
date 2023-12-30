package login

import (
	"mashu.example/internal/usecase/repository"
	"mashu.example/pkg/jwt"
)

type LoginUseCaseReq struct {
	userName string
}

type LoginUseCaseRes struct {
	AccessToken string
	Err         error
}

type LoginUseCase struct {
	userRepo repository.UserRepo
	req      *LoginUseCaseReq
	res      *LoginUseCaseRes
}

func (uc *LoginUseCase) Execute() {
	user, err := uc.userRepo.GetUserByUserName(uc.req.userName)
	if err != nil {
		uc.res.Err = err
		return
	}

	jwtClient := jwt.NewJWTClient(*jwt.NewAuthConfig())
	token, err := jwtClient.CreateToken(user.ID)
	if err != nil {
		uc.res.Err = err
		return
	}

	uc.res.AccessToken = token
}

func NewLoginUseCase(
	userRepo repository.UserRepo,
	req *LoginUseCaseReq,
	res *LoginUseCaseRes,
) *LoginUseCase {
	return &LoginUseCase{userRepo, req, res}
}

func NewLoginUseCaseReq(userName string) *LoginUseCaseReq {
	return &LoginUseCaseReq{userName}
}

func NewLoginUseCaseRes() *LoginUseCaseRes {
	return &LoginUseCaseRes{}
}
