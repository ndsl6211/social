package create_post_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"mashu.example/internal/entity"
	entity_enums "mashu.example/internal/entity/enums"
	"mashu.example/internal/usecase/post/create_post"
	"mashu.example/internal/usecase/tests"
)

func TestCreatePost(t *testing.T) {
	userRepo, postRepo, groupRepo, _ := tests.SetupTestRepositories(t)

	owner := entity.NewUser(uuid.New(), "owner", "Owner", "owner@email.com", false)

	userRepo.EXPECT().GetUserById(owner.ID).Return(owner, nil)

	var resultPost *entity.Post
	postRepo.EXPECT().Save(gomock.AssignableToTypeOf(&entity.Post{})).Do(
		func(arg *entity.Post) { resultPost = arg },
	)

	req := create_post.NewCreatePostUseCaseReq(
		"Hi, Golang",
		"Hello world!\nHello Clean Architecture!\nHello Domain Driven Design!",
		owner.ID,
		uuid.Nil,
		entity_enums.POST_PUBLIC,
	)
	res := create_post.NewCreatePostUseCaseRes()
	uc := create_post.NewCreatePostUseCase(userRepo, postRepo, groupRepo, req, res)

	uc.Execute()

	assert.Nil(t, res.Err)
	assert.Equal(t, resultPost.Title, "Hi, Golang")
	assert.Equal(t, resultPost.Content, "Hello world!\nHello Clean Architecture!\nHello Domain Driven Design!")
	assert.Equal(t, resultPost.Owner, owner)
	assert.Equal(t, resultPost.Permission, entity_enums.POST_PUBLIC)
	assert.Nil(t, resultPost.Group())
}

func TestCreatePostButOwnerDoesNotExist(t *testing.T) {
	userRepo, postRepo, groupRepo, _ := tests.SetupTestRepositories(t)

	owner := entity.NewUser(uuid.New(), "owner", "Owner", "owner@email.com", false)

	userRepo.EXPECT().GetUserById(owner.ID).Return(nil, gorm.ErrRecordNotFound)

	req := create_post.NewCreatePostUseCaseReq(
		"Hi, Golang",
		"Hello world!\nHello Clean Architecture!\nHello Domain Driven Design!",
		owner.ID,
		uuid.Nil,
		entity_enums.POST_PUBLIC,
	)
	res := create_post.NewCreatePostUseCaseRes()
	uc := create_post.NewCreatePostUseCase(userRepo, postRepo, groupRepo, req, res)

	uc.Execute()

	assert.ErrorIs(t, res.Err, create_post.ErrOwnerNotFound)
}

func TestCreatePostInGroup(t *testing.T) {
	userRepo, postRepo, groupRepo, _ := tests.SetupTestRepositories(t)

	owner := entity.NewUser(uuid.New(), "owner", "Owner", "owner@email.com", false)
	group := entity.NewGroup(uuid.New(), "group", owner, entity_enums.GROUP_PUBLIC)

	userRepo.EXPECT().GetUserById(owner.ID).Return(owner, nil)
	groupRepo.EXPECT().GetGroupById(group.ID).Return(group, nil)

	var resultPost *entity.Post
	postRepo.EXPECT().Save(gomock.AssignableToTypeOf(&entity.Post{})).Do(
		func(arg *entity.Post) { resultPost = arg },
	)

	req := create_post.NewCreatePostUseCaseReq(
		"Hi, Golang",
		"Hello world!\nHello Clean Architecture!\nHello Domain Driven Design!",
		owner.ID,
		group.ID,
		entity_enums.POST_PUBLIC,
	)
	res := create_post.NewCreatePostUseCaseRes()
	uc := create_post.NewCreatePostUseCase(userRepo, postRepo, groupRepo, req, res)

	uc.Execute()

	assert.Nil(t, res.Err)
	assert.Equal(t, resultPost.Title, "Hi, Golang")
	assert.Equal(t, resultPost.Content, "Hello world!\nHello Clean Architecture!\nHello Domain Driven Design!")
	assert.Equal(t, resultPost.Owner, owner)
	assert.Equal(t, resultPost.Permission, entity_enums.POST_PUBLIC)
	assert.Equal(t, resultPost.Group().ID, group.ID)
	assert.Equal(t, resultPost.Group().Owner.ID, group.Owner.ID)
}

func TestCreatePostInNonExistentGroup(t *testing.T) {
	userRepo, postRepo, groupRepo, _ := tests.SetupTestRepositories(t)

	owner := entity.NewUser(uuid.New(), "owner", "Owner", "owner@email.com", false)
	group := entity.NewGroup(uuid.New(), "group", owner, entity_enums.GROUP_PUBLIC)

	userRepo.EXPECT().GetUserById(owner.ID).Return(owner, nil)
	groupRepo.EXPECT().GetGroupById(group.ID).Return(nil, gorm.ErrRecordNotFound)

	req := create_post.NewCreatePostUseCaseReq(
		"Hi, Golang",
		"Hello world!\nHello Clean Architecture!\nHello Domain Driven Design!",
		owner.ID,
		group.ID,
		entity_enums.POST_PUBLIC,
	)
	res := create_post.NewCreatePostUseCaseRes()
	uc := create_post.NewCreatePostUseCase(userRepo, postRepo, groupRepo, req, res)

	uc.Execute()

	assert.ErrorIs(t, res.Err, create_post.ErrGroupNotFound)
}

func TestCreatePostInGroupWithInvalidPermission(t *testing.T) {
	userRepo, postRepo, groupRepo, _ := tests.SetupTestRepositories(t)

	owner := entity.NewUser(uuid.New(), "owner", "Owner", "owner@email.com", false)
	group := entity.NewGroup(uuid.New(), "group", owner, entity_enums.GROUP_PUBLIC)

	userRepo.EXPECT().GetUserById(owner.ID).Return(owner, nil)
	groupRepo.EXPECT().GetGroupById(group.ID).Return(group, nil)

	req := create_post.NewCreatePostUseCaseReq(
		"Hi, Golang",
		"Hello world!\nHello Clean Architecture!\nHello Domain Driven Design!",
		owner.ID,
		group.ID,
		entity_enums.POST_PRIVATE,
	)
	res := create_post.NewCreatePostUseCaseRes()
	uc := create_post.NewCreatePostUseCase(userRepo, postRepo, groupRepo, req, res)

	uc.Execute()

	assert.ErrorIs(t, res.Err, create_post.ErrInvalidPostPermission)
}
