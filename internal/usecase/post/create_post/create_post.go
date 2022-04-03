package create_post

import (
	"github.com/google/uuid"
	"mashu.example/internal/entity"
	"mashu.example/internal/usecase"
	"mashu.example/internal/usecase/repository"
)

type CreatePostUseCaseReq struct {
	Title   string
	Content string
	OwnerId uuid.UUID
	Public  bool
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

	post := entity.NewPost(uuid.New(), uc.req.Title, uc.req.Content, owner, uc.req.Public)
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
	public bool,
) *CreatePostUseCaseReq {
	return &CreatePostUseCaseReq{title, content, ownerId, public}
}

func NewCreatePostUseCaseRes() *CreatePostUseCaseRes {
	return &CreatePostUseCaseRes{}
}
