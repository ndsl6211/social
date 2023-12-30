package logout

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"mashu.example/internal/usecase"
	"mashu.example/internal/usecase/repository"
)

type LogoutUseCaseReq struct {
	UserName string
}

type LogoutUseCaseRes struct {
	Err error
}

type LogoutUseCase struct {
	userRepo  repository.UserRepo
	loginRepo repository.LoginRepo
	req       *LogoutUseCaseReq
	res       *LogoutUseCaseRes
}

func (uc *LogoutUseCase) Execute() {
	user, err := uc.userRepo.GetUserByUserName(uc.req.UserName)
	if err != nil {
		err := errors.New(fmt.Sprintf("failed to get user by username %s", uc.req.UserName))
		logrus.Error(err)
		uc.res.Err = err
	}

	if err := uc.loginRepo.ClearUserSess(user.ID); err != nil {
		err := errors.New(fmt.Sprintf("failed to clear user %s's login session", user.UserName))
		logrus.Error(err)
		uc.res.Err = err
	}
}

func NewLogoutUseCase(
	userRepo repository.UserRepo,
	loginRepo repository.LoginRepo,
	req *LogoutUseCaseReq,
	res *LogoutUseCaseRes,
) usecase.UseCase {
	return &LogoutUseCase{userRepo, loginRepo, req, res}
}
