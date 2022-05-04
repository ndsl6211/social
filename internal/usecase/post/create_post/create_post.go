package create_post

import (
	"github.com/google/uuid"
	"mashu.example/internal/entity"
	entity_enums "mashu.example/internal/entity/enums"
	"mashu.example/internal/usecase"
	"mashu.example/internal/usecase/repository"
)

type CreatePostUseCaseReq struct {
	title      string
	content    string
	ownerId    uuid.UUID
	permission entity_enums.PostPermission
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
	owner, err := uc.userRepo.GetUserById(uc.req.ownerId)
	if err != nil {
		uc.res.Err = err
		return
	}

	post := entity.NewPost(uuid.New(), uc.req.title, uc.req.content, owner, uc.req.permission)
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
	permission entity_enums.PostPermission,
) *CreatePostUseCaseReq {
	return &CreatePostUseCaseReq{title, content, ownerId, permission}
}

func NewCreatePostUseCaseRes() *CreatePostUseCaseRes {
	return &CreatePostUseCaseRes{}
}
