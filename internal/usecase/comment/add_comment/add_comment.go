package add_comment

import (
	"errors"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"mashu.example/internal/entity"
	entity_enums "mashu.example/internal/entity/enums"
	"mashu.example/internal/usecase"
	"mashu.example/internal/usecase/repository"
)

var (
	ErrAddCommentUnderPrivatePost     = errors.New("can not add comment under private post")
	ErrAddCommentUnderFollowrOnlyPost = errors.New("only the follower can comment under the follower-only post")
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
		uc.res.Err = &repository.ErrPostNotFound{PostId: uc.req.postId}
		logrus.Error(uc.res.Err)
		return
	}

	if post.Permission == entity_enums.POST_PRIVATE {
		uc.res.Err = ErrAddCommentUnderPrivatePost
		logrus.Error(uc.res.Err)
		return
	}

	commentOwner, err := uc.userRepo.GetUserById(uc.req.ownerId)
	if err != nil {
		uc.res.Err = &repository.ErrUserNotFound{UserId: uc.req.ownerId}
		logrus.Error(uc.res.Err)
		return
	}

	if post.Permission == entity_enums.POST_FOLLOWER_ONLY {
		isFollower := false
		for _, followerID := range post.Owner.Followers {
			if followerID == commentOwner.ID {
				isFollower = true
				break
			}
		}

		if !isFollower {
			uc.res.Err = ErrAddCommentUnderFollowrOnlyPost
			logrus.Error(uc.res.Err)
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
