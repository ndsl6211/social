package create_post_test

import (
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"mashu.example/internal/entity"
	"mashu.example/internal/usecase/create_post"
	"mashu.example/internal/usecase/repository/mock"
)

func setup(t *testing.T) (*mock.MockPostRepo, *mock.MockUserRepo) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	return mock.NewMockPostRepo(mockCtrl), mock.NewMockUserRepo(mockCtrl)
}

func TestCreatePost(t *testing.T) {
	postRepo, userRepo := setup(t)

	ownerId := uuid.MustParse("00000000-0000-0000-0000-000000000000")
	owner := entity.NewUser(
		ownerId,
		"owner",
		"owner display name",
		"owner@email.com",
		false,
	)

	userRepo.EXPECT().GetUserById(ownerId).Return(owner, nil)

	var resultPost *entity.Post
	postRepo.EXPECT().Save(gomock.AssignableToTypeOf(&entity.Post{})).Do(
		func(arg *entity.Post) { resultPost = arg },
	)

	req := create_post.NewCreatePostUseCaseReq(
		"Hi, Golang",
		"Hello world!\nHello Clean Architecture!\nHello Domain Driven Design!",
		ownerId,
		true,
	)
	res := create_post.NewCreatePostUseCaseRes()
	uc := create_post.NewCreatePostUseCase(userRepo, postRepo, &req, &res)

	uc.Execute()

	if res.Err != nil {
		t.Errorf("failed to execute usecase")
	}

	assert.Equal(t, resultPost.Title, "Hi, Golang")
	assert.Equal(t, resultPost.Content, "Hello world!\nHello Clean Architecture!\nHello Domain Driven Design!")
	assert.Equal(t, resultPost.Owner, owner)
	assert.Equal(t, resultPost.Public, true)
}
