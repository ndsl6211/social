package add_comment

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"mashu.example/internal/entity"
	"mashu.example/internal/entity/enums/post_permission"
	"mashu.example/internal/usecase"
	"mashu.example/internal/usecase/repository"
)

type AddCommentUseCaseReq struct {
	ownerId uuid.UUID
	postId  uuid.UUID
	content string
}

type AddCommentUseCaseRes struct {
	Err error
}

type AddCommentUseCase struct {
	userRepo repository.UserRepo
	postRepo repository.PostRepo
	req      *AddCommentUseCaseReq
	res      *AddCommentUseCaseRes
}

func (uc *AddCommentUseCase) Execute() {
	post, err := uc.postRepo.GetPostById(uc.req.postId)
	if err != nil {
		logrus.Errorf("failed to get post (postId: %s)", uc.req.postId)
		uc.res.Err = err
		return
	}

	if post.Permission == post_permission.PRIVATE {
		errMsg := fmt.Sprintf("can not add comment under private post (postId: %s)", uc.req.postId)
		logrus.Errorf(errMsg)
		uc.res.Err = errors.New(errMsg)
		return
	}

	commentOwner, err := uc.userRepo.GetUserById(uc.req.ownerId)
	if err != nil {
		logrus.Errorf("failed to get comment owner (userId: %s)", uc.req.ownerId)
		uc.res.Err = err
		return
	}

	if post.Permission == post_permission.FOLLOWER_ONLY {
		isFollower := false
		for _, followerID := range post.Owner.Followers {
			if followerID == commentOwner.ID {
				isFollower = true
				break
			}
		}

		if !isFollower {
			errMsg := "only the follower can comment below the follower-only post"
			logrus.Errorf(errMsg)
			uc.res.Err = errors.New(errMsg)
			return
		}
	}

	post.Comments = append(post.Comments, entity.NewComment(
		uuid.New(),
		commentOwner,
		post,
		uc.req.content,
	))

	uc.postRepo.Save(post)
}

func NewAddCommentUseCase(
	userRepo repository.UserRepo,
	postRepo repository.PostRepo,
	req *AddCommentUseCaseReq,
	res *AddCommentUseCaseRes,
) usecase.UseCase {
	return &AddCommentUseCase{userRepo, postRepo, req, res}
}

func NewAddCommentUseCaseReq(
	ownerId uuid.UUID,
	postId uuid.UUID,
	content string,
) *AddCommentUseCaseReq {
	return &AddCommentUseCaseReq{ownerId, postId, content}
}

func NewAddCommentUseCaseRes() *AddCommentUseCaseRes {
	return &AddCommentUseCaseRes{}
}
