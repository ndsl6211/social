package register

import (
	"errors"

	"github.com/google/uuid"
	"mashu.example/internal/model"
	"mashu.example/internal/usecase/repository"
)

type RegisterUseCaseReq struct {
	username string
	email    string
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
	if _, err := uc.userRepo.GetUserByUserName(uc.req.username); err == nil {
		err = errors.New("user already exist")
		uc.res.Err = err
		return
	}

	user := model.NewUser(
		uuid.New(),
		uc.req.username,
		uc.req.email,
		false,
	)

	if err := uc.userRepo.Save(user); err != nil {
		uc.res.Err = err
		return
	}
}

func NewRegisterUseCase(
	userRepo repository.UserRepo,
	req *RegisterUseCaseReq,
	res *RegisterUseCaseRes,
) *RegisterUseCase {
	return &RegisterUseCase{userRepo, req, res}
}

func NewRegisterUseCaseReq(
	username string,
	email string,
) *RegisterUseCaseReq {
	return &RegisterUseCaseReq{username, email}
}

func NewRegisterUseCaseRes() *RegisterUseCaseRes {
	return &RegisterUseCaseRes{}
}
