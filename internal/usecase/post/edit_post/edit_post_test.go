package edit_post_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"mashu.example/internal/entity"
	entity_enums "mashu.example/internal/entity/enums"
	"mashu.example/internal/usecase/post/edit_post"
	"mashu.example/internal/usecase/repository/mock"
)

func setup(t *testing.T) (*mock.MockPostRepo, *mock.MockUserRepo) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	return mock.NewMockPostRepo(mockCtrl), mock.NewMockUserRepo(mockCtrl)
}

func TestEditPost(t *testing.T) {
	postRepo, _ := setup(t)

	postId := uuid.New()
	post := entity.NewPost(
		postId,
		"My First Post",
		"My first content",
		&entity.User{
			ID:          uuid.New(),
			UserName:    "post_owner",
			DisplayName: "Post Owner",
			Email:       "owner@email.com",
			Public:      true,
		},
		entity_enums.POST_PUBLIC,
	)

	var updatedPost *entity.Post
	postRepo.EXPECT().GetPostById(postId).Return(post, nil)
	postRepo.EXPECT().Save(gomock.AssignableToTypeOf(&entity.Post{})).Do(
		func(arg *entity.Post) { updatedPost = arg },
	)

	newTitle := "My First Post (revised)"
	newContent := "My first content (revised)"
	req := edit_post.NewEditPostUseCaseReq(postId, post.Owner.ID, newTitle, newContent)
	res := edit_post.NewEditPostUseCaseRes()
	uc := edit_post.NewEditPostUseCase(postRepo, req, res)

	uc.Execute()

	assert.Nil(t, res.Err)
	assert.Equal(t, updatedPost.ID, post.ID)
	assert.Equal(t, updatedPost.Title, newTitle)
	assert.Equal(t, updatedPost.Content, newContent)
	assert.Equal(t, updatedPost.Owner.ID, post.Owner.ID)
	assert.Greater(t, updatedPost.UpdatedAt, updatedPost.CreatedAt)
}

func TestEditNotMyOwnPost(t *testing.T) {
	postRepo, _ := setup(t)

	postId := uuid.New()
	post := entity.NewPost(
		postId,
		"My First Post",
		"My first content",
		&entity.User{
			ID:          uuid.New(),
			UserName:    "post_owner",
			DisplayName: "Post Owner",
			Email:       "owner@email.com",
			Public:      true,
		},
		entity_enums.POST_PUBLIC,
	)

	postRepo.EXPECT().GetPostById(postId).Return(post, nil)
	nonOwnerId := uuid.New()

	newTitle := "My First Post (revised)"
	newContent := "My first content (revised)"
	req := edit_post.NewEditPostUseCaseReq(postId, nonOwnerId, newTitle, newContent)
	res := edit_post.NewEditPostUseCaseRes()
	uc := edit_post.NewEditPostUseCase(postRepo, req, res)

	uc.Execute()

	assert.ErrorIs(t, res.Err, edit_post.ErrNotOwnerOfPost)
}
