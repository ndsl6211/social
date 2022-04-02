package edit_post

import (
	"errors"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"mashu.example/internal/usecase"
	"mashu.example/internal/usecase/repository"
)

type EditPostUseCaseReq struct {
	postId     uuid.UUID
	ownerId    uuid.UUID
	newTitle   string
	newContent string
}

type EditPostUseCaseRes struct {
	Err     error
	message string
}

type EditPostUseCase struct {
	postRepo repository.PostRepo
	req      *EditPostUseCaseReq
	res      *EditPostUseCaseRes
}

func (uc *EditPostUseCase) Execute() {
	post, err := uc.postRepo.GetPostById(uc.req.postId)
	if err != nil {
		logrus.Infof("failed to get post (postId: %s)", uc.req.postId)
		uc.res.Err = err
		return
	}

	if post.Owner.ID != uc.req.ownerId {
		message := "only the post owner can edit the post"
		logrus.Error(message)
		uc.res.Err = errors.New(message)
		return
	}

	post.Title = uc.req.newTitle
	post.Content = uc.req.newContent

	uc.postRepo.Save(post)
}

func NewEditPostUseCase(
	postRepo repository.PostRepo,
	req *EditPostUseCaseReq,
	res *EditPostUseCaseRes,
) usecase.UseCase {
	return &EditPostUseCase{postRepo, req, res}
}

func NewEditPostUseCaseReq(
	postId uuid.UUID,
	ownerId uuid.UUID,
	newTitle string,
	newContent string,
) *EditPostUseCaseReq {
	return &EditPostUseCaseReq{postId, ownerId, newTitle, newContent}
}

func NewEditPostUseCaseRes() *EditPostUseCaseRes {
	return &EditPostUseCaseRes{}
}
