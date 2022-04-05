package edit_post_test

import (
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"mashu.example/internal/entity"
	"mashu.example/internal/entity/enums/post_permission"
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

	postId := uuid.MustParse("10101010-1010-1010-1010-101010101010")
	post := entity.NewPost(
		postId,
		"My First Post",
		"My first content",
		&entity.User{
			ID:          uuid.MustParse("01010101-0101-0101-0101-010101010101"),
			UserName:    "post_owner",
			DisplayName: "Post Owner",
			Email:       "owner@email.com",
			Public:      true,
		},
		post_permission.PUBLIC,
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

	if res.Err != nil {
		t.Errorf("failed to execute usecase")
	}

	assert.Equal(t, updatedPost.ID, post.ID)
	assert.Equal(t, updatedPost.Title, newTitle)
	assert.Equal(t, updatedPost.Content, newContent)
	assert.Equal(t, updatedPost.Owner.ID, post.Owner.ID)
}

func TestEditNotMyOwnPost(t *testing.T) {
	postRepo, _ := setup(t)

	postId := uuid.MustParse("10101010-1010-1010-1010-101010101010")
	post := entity.NewPost(
		postId,
		"My First Post",
		"My first content",
		&entity.User{
			ID:          uuid.MustParse("01010101-0101-0101-0101-010101010101"),
			UserName:    "post_owner",
			DisplayName: "Post Owner",
			Email:       "owner@email.com",
			Public:      true,
		},
		post_permission.PUBLIC,
	)

	postRepo.EXPECT().GetPostById(postId).Return(post, nil)
	nonOwnerId := uuid.MustParse("02020202-0202-0202-0202-020202020202")

	newTitle := "My First Post (revised)"
	newContent := "My first content (revised)"
	req := edit_post.NewEditPostUseCaseReq(postId, nonOwnerId, newTitle, newContent)
	res := edit_post.NewEditPostUseCaseRes()
	uc := edit_post.NewEditPostUseCase(postRepo, req, res)

	uc.Execute()

	if res.Err == nil {
		t.Errorf("only the post owner can edit the post")
		return
	}

	assert.Equal(t, res.Err.Error(), "only the post owner can edit the post")
}
