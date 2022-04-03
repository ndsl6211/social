package add_comment_test

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"mashu.example/internal/entity"
	"mashu.example/internal/usecase/comment/add_comment"
	"mashu.example/internal/usecase/repository/mock"
)

func setup(t *testing.T) (*mock.MockPostRepo, *mock.MockUserRepo) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	return mock.NewMockPostRepo(mockCtrl), mock.NewMockUserRepo(mockCtrl)
}

func TestAddComment(t *testing.T) {
	postRepo, userRepo := setup(t)

	ownerId := uuid.New()
	owner := entity.NewUser(ownerId, "owner", "owner display name", "owner@email.com", true)

	postId := uuid.New()
	post := entity.NewPost(
		postId,
		"My First Post",
		"My first content",
		entity.NewUser(ownerId, "post_owner", "owner display name", "owner@email.com", true),
		true,
	)

	postRepo.EXPECT().GetPostById(postId).Return(post, nil)
	userRepo.EXPECT().GetUserById(ownerId).Return(owner, nil)
	postRepo.EXPECT().Save(gomock.AssignableToTypeOf(&entity.Post{})).Do(
		func(arg *entity.Post) { post = arg },
	)

	req := add_comment.NewAddCommentUseCaseReq(ownerId, postId, "Good!")
	res := add_comment.NewAddCommentUseCaseRes()
	uc := add_comment.NewAddCommentUseCase(userRepo, postRepo, req, res)

	uc.Execute()

	if res.Err != nil {
		t.Errorf("failed to execute usecase")
	}

	fmt.Println(post)

	assert.Nil(t, res.Err)
	assert.Equal(t, "Good!", post.Comments[0].Content)
	assert.Equal(t, post.ID, post.Comments[0].Post.ID)
	assert.Equal(t, owner.ID, post.Comments[0].Owner.ID)
}
