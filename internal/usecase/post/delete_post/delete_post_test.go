package delete_post_test

import (
	"errors"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"mashu.example/internal/entity"
	"mashu.example/internal/usecase/post/delete_post"
	"mashu.example/internal/usecase/repository/mock"
)

func setup(t *testing.T) *mock.MockPostRepo {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	return mock.NewMockPostRepo(mockCtrl)
}

func TestDeletePost(t *testing.T) {
	postRepo := setup(t)

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
		true,
	)

	postRepo.EXPECT().GetPostById(postId).Return(post, nil)
	postRepo.EXPECT().Delete(postId).Return(nil)

	req := delete_post.NewDeletePostUseCaseReq(postId, post.Owner.ID)
	res := delete_post.NewDeletePostUseCaseRes()
	uc := delete_post.NewDeletePoseUseCase(postRepo, req, res)

	uc.Execute()

	if res.Err != nil {
		t.Errorf("failed to execute usecase")
	}
}

func TestDeleteNonExistPost(t *testing.T) {
	postRepo := setup(t)

	postId := uuid.MustParse("10101010-1010-1010-1010-101010101010")
	userId := uuid.MustParse("01010101-0101-0101-0101-010101010101")
	postRepo.EXPECT().GetPostById(postId).Return(nil, gorm.ErrRecordNotFound)

	req := delete_post.NewDeletePostUseCaseReq(postId, userId)
	res := delete_post.NewDeletePostUseCaseRes()
	uc := delete_post.NewDeletePoseUseCase(postRepo, req, res)

	uc.Execute()

	if res.Err == nil {
		t.Errorf("failed to execute usecase")
	}

	assert.Equal(t, errors.Is(res.Err, gorm.ErrRecordNotFound), true)
}

func TestDeleteNotMyPost(t *testing.T) {
	postRepo := setup(t)

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
		true,
	)

	postRepo.EXPECT().GetPostById(postId).Return(post, nil)
	notOwnerId := uuid.MustParse("02020202-0202-0202-0202-020202020202")

	req := delete_post.NewDeletePostUseCaseReq(postId, notOwnerId)
	res := delete_post.NewDeletePostUseCaseRes()
	uc := delete_post.NewDeletePoseUseCase(postRepo, req, res)

	uc.Execute()

	if res.Err == nil {
		t.Errorf("failed to execute usecase")
	}

	assert.Equal(t, res.Err.Error(), "only the post owner can delete the post")
}
