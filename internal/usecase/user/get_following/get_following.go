package get_following

import (
	"github.com/google/uuid"
	"mashu.example/internal/usecase"
	"mashu.example/internal/usecase/dto"
	"mashu.example/internal/usecase/repository"
)

type GetFollowingUseCaseReq struct {
	userId uuid.UUID
}

type GetFollowingUseCaseRes struct {
	Users []*dto.FollowingInfo
	Err   error
}

type GetFollowingUseCase struct {
	userRepo repository.UserRepo
	req      *GetFollowingUseCaseReq
	res      *GetFollowingUseCaseRes
}

func (uc *GetFollowingUseCase) Execute() {
	user, err := uc.userRepo.GetUserById(uc.req.userId)
	if err != nil {
		uc.res.Err = err
		return
	}

	followingInfos := []*dto.FollowingInfo{}
	for _, followingId := range user.Followings {
		followingUser, err := uc.userRepo.GetUserById(followingId)
		if err != nil {
			uc.res.Err = err
			return
		}

		followingInfos = append(followingInfos, &dto.FollowingInfo{
			ID:       followingUser.ID,
			UserName: followingUser.UserName,
			Email:    followingUser.Email,
			Public:   followingUser.Public,
		})
	}

	uc.res.Users = followingInfos
	uc.res.Err = nil
}

func NewGetFollowingUseCase(
	userRepo repository.UserRepo,
	req *GetFollowingUseCaseReq,
	res *GetFollowingUseCaseRes,
) usecase.UseCase {
	return &GetFollowingUseCase{userRepo, req, res}
}

func NewGetFollowingUseCaseReq(userId uuid.UUID) *GetFollowingUseCaseReq {
	return &GetFollowingUseCaseReq{userId}
}

func NewGetFollowingUseCaseRes() *GetFollowingUseCaseRes {
	return &GetFollowingUseCaseRes{}
}
