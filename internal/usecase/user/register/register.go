package register

import (
	"github.com/google/uuid"
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
		uc.res.Err = err
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

func NewRegisterUseCaseReq() *RegisterUseCaseReq {
	return &RegisterUseCaseReq{}
}

func NewRegisterUseCaseRes() *RegisterUseCaseRes {
	return &RegisterUseCaseRes{}
}
