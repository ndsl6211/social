package login

import (
	"github.com/sirupsen/logrus"
	"mashu.example/internal/usecase"
	"mashu.example/internal/usecase/repository"
	"mashu.example/pkg/jwt"
)

type LoginUseCaseReq struct {
	username string
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
	user, err := uc.userRepo.GetUserByUserName(uc.req.username)
	if err != nil {
		logrus.Errorf("user %s doesn't exist", uc.req.username)
		uc.res.Err = err
		return
	}

	jwtClient := jwt.NewJWTClient(*jwt.NewAuthConfig())
	token, err := jwtClient.CreateToken(user.ID)
	if err != nil {
		logrus.Errorf("failed to generate token")
		uc.res.Err = err
		return
	}

	uc.res.AccessToken = token
}

func NewLoginUseCase(
	userRepo repository.UserRepo,
	req *LoginUseCaseReq,
	res *LoginUseCaseRes,
) usecase.UseCase {
	return &LoginUseCase{userRepo, req, res}
}

func NewLoginUseCaseReq(username string) *LoginUseCaseReq {
	return &LoginUseCaseReq{username}
}

func NewLoginUseCaseRes() *LoginUseCaseRes {
	return &LoginUseCaseRes{}
}
