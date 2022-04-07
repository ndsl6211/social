package register

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"mashu.example/internal/entity"
	"mashu.example/internal/usecase"
	"mashu.example/internal/usecase/repository"
)

type RegisterUseCaseReq struct {
	username    string
	displayName string
	email       string
}

type RegisterUseCaseRes struct {
	Err error
}

type RegisterUseCase struct {
	userRepo repository.UserRepo
	req      *RegisterUseCaseReq
	res      *RegisterUseCaseRes
}

func (uc *RegisterUseCase) Execute() {
	user := entity.NewUser(
		uuid.New(),
		uc.req.username,
		uc.req.displayName,
		uc.req.email,
		false,
	)

	if err := uc.userRepo.Save(user); err != nil {
		errMsg := fmt.Sprintf("user %s already exist", uc.req.username)
		logrus.Info(errMsg)
		uc.res.Err = errors.New(errMsg)
		return
	}
}

func NewRegisterUseCase(
	userRepo repository.UserRepo,
	req *RegisterUseCaseReq,
	res *RegisterUseCaseRes,
) usecase.UseCase {
	return &RegisterUseCase{userRepo, req, res}
}

func NewRegisterUseCaseReq(
	username string,
	displayName string,
	email string,
) *RegisterUseCaseReq {
	return &RegisterUseCaseReq{username, displayName, email}
}

func NewRegisterUseCaseRes() *RegisterUseCaseRes {
	return &RegisterUseCaseRes{}
}
