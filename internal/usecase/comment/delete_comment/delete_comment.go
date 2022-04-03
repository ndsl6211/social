package delete_comment

import (
	"errors"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
	"mashu.example/internal/entity"
	"mashu.example/internal/usecase"
	"mashu.example/internal/usecase/repository"
)

type DeleteCommentUseCaseReq struct {
	ownerId   uuid.UUID
	postId    uuid.UUID
	commentId uuid.UUID
}

type DeleteCommentUseCaseRes struct {
	Err error
}

type DeleteCommentUseCase struct {
	userRepo repository.UserRepo
	postRepo repository.PostRepo
	req      *DeleteCommentUseCaseReq
	res      *DeleteCommentUseCaseRes
}

func (uc *DeleteCommentUseCase) Execute() {
	post, err := uc.postRepo.GetPostById(uc.req.postId)
	if err != nil {
		logrus.Errorf("failed to get post (postId: %s)", uc.req.postId)
		uc.res.Err = err
		return
	}

	idx := slices.IndexFunc(post.Comments, func(comment *entity.Comment) bool {
		return comment.ID == uc.req.commentId
	})

	if post.Comments[idx].Owner.ID != uc.req.ownerId {
		errMsg := "only the comment owner can delete this comment"
		logrus.Errorf(errMsg)
		uc.res.Err = errors.New(errMsg)
		return
	}

	post.Comments = slices.Delete(post.Comments, idx, idx+1)

	uc.postRepo.Save(post)
}

func NewDeletePoseUseCase(
	userRepo repository.UserRepo,
	postRepo repository.PostRepo,
	req *DeleteCommentUseCaseReq,
	res *DeleteCommentUseCaseRes,
) usecase.UseCase {
	return &DeleteCommentUseCase{userRepo, postRepo, req, res}
}

func NewDeletePostUseCaseReq(
	ownerId uuid.UUID,
	postId uuid.UUID,
	commentId uuid.UUID,
) *DeleteCommentUseCaseReq {
	return &DeleteCommentUseCaseReq{ownerId, postId, commentId}
}

func NewDeletePostUseCaseRes() *DeleteCommentUseCaseRes {
	return &DeleteCommentUseCaseRes{}
}
