package handle_follow_request

import (
	"errors"

	"github.com/google/uuid"
	"golang.org/x/exp/slices"
	"mashu.example/internal/entity"
	"mashu.example/internal/usecase"
	"mashu.example/internal/usecase/repository"
)

type HandleFollowRequestAction string

const (
	ACCEPT_FOLLOW_REQUEST HandleFollowRequestAction = "ACCEPT"
	REJECT_FOLLOW_REQUEST HandleFollowRequestAction = "REJECT"
)

type HandleFollowRequestUsecaseReq struct {
	userId     string
	followerId string
	action     HandleFollowRequestAction
}

type HandleFollowRequestUsecaseRes struct {
	Err error
}

type HandleFollowRequestUsecase struct {
	userRepo repository.UserRepo
	Req      *HandleFollowRequestUsecaseReq
	Res      *HandleFollowRequestUsecaseRes
}

func (uc *HandleFollowRequestUsecase) Execute() {
	userId := uuid.MustParse(uc.Req.userId)
	followerId := uuid.MustParse(uc.Req.followerId)

	user, err := uc.userRepo.GetUserById(userId)
	if err != nil {
		uc.Res.Err = err
		return
	}
	follower, err := uc.userRepo.GetUserById(followerId)
	if err != nil {
		uc.Res.Err = err
		return
	}

	var targetFollowReqIdx int
	var targetFollowReq *entity.FollowRequest = nil
	for idx, followReq := range user.FollowRequests {
		if followReq.From == followerId && followReq.To == userId {
			targetFollowReq = followReq
			targetFollowReqIdx = idx
			break
		}
	}
	if targetFollowReq == nil {
		uc.Res.Err = errors.New("follow request not found in followee")
		return
	}
	user.FollowRequests = slices.Delete(user.FollowRequests, targetFollowReqIdx, targetFollowReqIdx+1)

	targetFollowReq = nil
	for idx, followReq := range follower.FollowRequests {
		if followReq.From == followerId && followReq.To == userId {
			targetFollowReq = followReq
			targetFollowReqIdx = idx
			break
		}
	}
	if targetFollowReq == nil {
		uc.Res.Err = errors.New("follow request not found in follower")
		return
	}
	follower.FollowRequests = slices.Delete(follower.FollowRequests, targetFollowReqIdx, targetFollowReqIdx+1)

	if uc.Req.action == ACCEPT_FOLLOW_REQUEST {
		user.AddFollower(followerId)
		follower.AddFollowing(userId)
	}

	uc.userRepo.Save(user)
	uc.userRepo.Save(follower)

	uc.Res.Err = nil
}

func NewHandleFollowRequestUsecase(
	userRepo repository.UserRepo,
	req *HandleFollowRequestUsecaseReq,
	res *HandleFollowRequestUsecaseRes,
) usecase.UseCase {
	return &HandleFollowRequestUsecase{userRepo, req, res}
}

func NewHandleFollowRequestUsecaseReq(
	userId string,
	followerId string,
	action HandleFollowRequestAction,
) HandleFollowRequestUsecaseReq {
	return HandleFollowRequestUsecaseReq{userId, followerId, action}
}

func NewHandleFollowRequestUsecaseRes() HandleFollowRequestUsecaseRes {
	return HandleFollowRequestUsecaseRes{}
}