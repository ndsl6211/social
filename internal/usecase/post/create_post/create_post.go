package create_post

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
	ErrOwnerNotFound         = errors.New("owner not found")
	ErrGroupNotFound         = errors.New("group not found")
	ErrInvalidPostPermission = errors.New("group post should be public")
)

type CreatePostUseCaseReq struct {
	Title      string
	Content    string
	OwnerId    uuid.UUID
	GroupId    uuid.UUID
	Permission entity_enums.PostPermission
}

type CreatePostUseCaseRes struct {
	Err error
}

type CreatePostUseCase struct {
	userRepo  repository.UserRepo
	postRepo  repository.PostRepo
	groupRepo repository.GroupRepo

	req *CreatePostUseCaseReq
	res *CreatePostUseCaseRes
}

func (uc *CreatePostUseCase) Execute() {
	owner, err := uc.userRepo.GetUserById(uc.req.ownerId)
	if err != nil {
		uc.res.Err = ErrOwnerNotFound
		logrus.Error(uc.res.Err)
		return
	}

	var group *entity.Group = nil
	if uc.req.groupId != uuid.Nil {
		if group, err = uc.groupRepo.GetGroupById(uc.req.groupId); err != nil {
			uc.res.Err = ErrGroupNotFound
			logrus.Error(uc.res.Err)
			return
		}
	}

	post := entity.NewPost(uuid.New(), uc.req.title, uc.req.content, owner, group, uc.req.permission)
	if post == nil {
		uc.res.Err = ErrInvalidPostPermission
		logrus.Error(uc.res.Err)
		return
	}
	uc.postRepo.Save(post)

	uc.res.Err = nil
}

func NewCreatePostUseCase(
	userRepo repository.UserRepo,
	postRepo repository.PostRepo,
	groupRepo repository.GroupRepo,
	req *CreatePostUseCaseReq,
	res *CreatePostUseCaseRes,
) usecase.UseCase {
	return &CreatePostUseCase{userRepo, postRepo, groupRepo, req, res}
}

func NewCreatePostUseCaseReq(
	title string,
	content string,
	ownerId uuid.UUID,
	groupId uuid.UUID,
	permission entity_enums.PostPermission,
) *CreatePostUseCaseReq {
	return &CreatePostUseCaseReq{title, content, ownerId, groupId, permission}
}

func NewCreatePostUseCaseRes() *CreatePostUseCaseRes {
	return &CreatePostUseCaseRes{}
}
