package create_post

import (
	"github.com/google/uuid"
	"mashu.example/internal/entity"
	"mashu.example/internal/entity/enums/post_permission"
	"mashu.example/internal/usecase"
	"mashu.example/internal/usecase/repository"
)

type CreatePostUseCaseReq struct {
	Title      string
	Content    string
	OwnerId    uuid.UUID
	Permission post_permission.PostPermission
}

type CreatePostUseCaseRes struct {
	Err error
}

type CreatePostUseCase struct {
	userRepo repository.UserRepo
	postRepo repository.PostRepo

	req *CreatePostUseCaseReq
	res *CreatePostUseCaseRes
}

func (uc *CreatePostUseCase) Execute() {
	owner, err := uc.userRepo.GetUserById(uc.req.OwnerId)
	if err != nil {
		uc.res.Err = err
		return
	}

	post := entity.NewPost(uuid.New(), uc.req.Title, uc.req.Content, owner, uc.req.Permission)
	uc.postRepo.Save(post)

	uc.res.Err = nil
}

func NewCreatePostUseCase(
	userRepo repository.UserRepo,
	postRepo repository.PostRepo,
	req *CreatePostUseCaseReq,
	res *CreatePostUseCaseRes,
) usecase.UseCase {
	return &CreatePostUseCase{userRepo, postRepo, req, res}
}

func NewCreatePostUseCaseReq(
	title string,
	content string,
	ownerId uuid.UUID,
	permission post_permission.PostPermission,
) *CreatePostUseCaseReq {
	return &CreatePostUseCaseReq{title, content, ownerId, permission}
}

func NewCreatePostUseCaseRes() *CreatePostUseCaseRes {
	return &CreatePostUseCaseRes{}
}