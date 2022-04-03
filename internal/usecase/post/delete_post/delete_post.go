package delete_post

import (
	"errors"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"mashu.example/internal/usecase"
	"mashu.example/internal/usecase/repository"
)

type DeletePostUseCaseReq struct {
	postId  uuid.UUID
	ownerId uuid.UUID
}

type DeletePostUseCaseRes struct {
	Err error
}

type DeletePostUseCase struct {
	postRepo repository.PostRepo
	req      *DeletePostUseCaseReq
	res      *DeletePostUseCaseRes
}

func (uc *DeletePostUseCase) Execute() {
	post, err := uc.postRepo.GetPostById(uc.req.postId)
	if err != nil {
		logrus.Errorf("failed to get post (postId: %s", uc.req.postId)
		uc.res.Err = err
		return
	}

	if post.Owner.ID != uc.req.ownerId {
		message := "only the post owner can delete the post"
		logrus.Error(message)
		uc.res.Err = errors.New(message)
		return
	}

	if err := uc.postRepo.Delete(post.ID); err != nil {
		logrus.Errorf("failed to delete post (postId: %s", uc.req.postId)
		uc.res.Err = err
		return
	}
}

func NewDeletePoseUseCase(
	postRepo repository.PostRepo,
	req *DeletePostUseCaseReq,
	res *DeletePostUseCaseRes,
) usecase.UseCase {
	return &DeletePostUseCase{postRepo, req, res}
}

func NewDeletePostUseCaseReq(postId uuid.UUID, userId uuid.UUID) *DeletePostUseCaseReq {
	return &DeletePostUseCaseReq{postId, userId}
}

func NewDeletePostUseCaseRes() *DeletePostUseCaseRes {
	return &DeletePostUseCaseRes{}
}
