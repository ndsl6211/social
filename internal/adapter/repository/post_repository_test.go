package adapter_repository_test

import (
	"errors"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	adapter_repository "mashu.example/internal/adapter/repository"
	"mashu.example/internal/entity"
	"mashu.example/internal/entity/enums/post_permission"
	"mashu.example/internal/usecase/repository"
	"mashu.example/pkg"
)

func setup() repository.PostRepo {
	db := pkg.NewMemoryGormClient()
	postRepo := adapter_repository.NewPostRepository(db)

	return postRepo
}

func TestGetPostById(t *testing.T) {
	postRepo := setup()

	postId := uuid.New()
	postRepo.Save(&entity.Post{
		ID:      postId,
		Title:   "Hi, Golang",
		Content: "Hi, Clean Architecture!\nHi, Domain Driven Design!",
		Owner: entity.NewUser(
			uuid.MustParse("10101010-0000-0000-0000-000000000000"),
			"owner",
			"owner display name",
			"owner@email.com",
			false,
		),
		Permission: post_permission.PUBLIC,
	})

	resultPost, err := postRepo.GetPostById(postId)
	if err != nil {
		t.Error("failed to get post")
	}

	assert.Equal(t, resultPost.ID, postId)
	assert.Equal(t, resultPost.Title, "Hi, Golang")
	assert.Equal(t, resultPost.Content, "Hi, Clean Architecture!\nHi, Domain Driven Design!")
	assert.Equal(t, resultPost.Permission, post_permission.PUBLIC)
	assert.Equal(t, resultPost.Owner.ID, uuid.MustParse("10101010-0000-0000-0000-000000000000"))
	assert.Equal(t, resultPost.Owner.UserName, "owner")
	assert.Equal(t, resultPost.Owner.DisplayName, "owner display name")
	assert.Equal(t, resultPost.Owner.Email, "owner@email.com")
	assert.Equal(t, resultPost.Owner.Public, false)
}

func TestGetNonExistPost(t *testing.T) {
	postRepo := setup()

	dummyPostId := uuid.MustParse("11111112-0000-0000-0000-000000000000")
	_, err := postRepo.GetPostById(dummyPostId)

	assert.Equal(t, errors.Is(err, gorm.ErrRecordNotFound), true)
}
