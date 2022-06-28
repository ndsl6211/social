package comment_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"mashu.example/internal/entity"
	entity_enums "mashu.example/internal/entity/enums"
	usecase "mashu.example/internal/usecase/comment/delete_comment"
	"mashu.example/internal/usecase/tests"
)

func TestDeleteComment(t *testing.T) {
	userRepo, postRepo, _, _ := tests.SetupTestRepositories(t)

	ownerId := uuid.New()
	owner := entity.NewUser(ownerId, "owner", "owner display name", "owner@email.com", true)

	postId := uuid.New()
	post := entity.NewPost(
		postId,
		"My First Post",
		"My first content",
		entity.NewUser(ownerId, "post_owner", "owner display name", "owner@email.com", true),
		entity_enums.POST_PUBLIC,
	)
	commentId := uuid.New()
	post.Comments = append(post.Comments, entity.NewComment(commentId, owner, post, "Good!"))
	assert.Equal(t, 1, len(post.Comments))

	postRepo.EXPECT().GetPostById(postId).Return(post, nil)
	postRepo.EXPECT().Save(gomock.AssignableToTypeOf(&entity.Post{})).Do(
		func(arg *entity.Post) { post = arg },
	)

	req := usecase.NewDeletePostUseCaseReq(ownerId, postId, commentId)
	res := usecase.NewDeletePostUseCaseRes()
	uc := usecase.NewDeletePoseUseCase(userRepo, postRepo, req, res)

	uc.Execute()

	if res.Err != nil {
		t.Errorf("failed to execute usecase")
	}

	assert.Nil(t, res.Err)
	assert.Equal(t, 0, len(post.Comments))
}

func TestDeleteNotMyOwnComment(t *testing.T) {
	userRepo, postRepo, _, _ := tests.SetupTestRepositories(t)

	ownerId := uuid.New()
	owner := entity.NewUser(ownerId, "owner", "owner display name", "owner@email.com", true)

	postId := uuid.New()
	post := entity.NewPost(
		postId,
		"My First Post",
		"My first content",
		entity.NewUser(ownerId, "post_owner", "owner display name", "owner@email.com", true),
		entity_enums.POST_PUBLIC,
	)
	commentId := uuid.New()
	post.Comments = append(post.Comments, entity.NewComment(commentId, owner, post, "Good!"))
	assert.Equal(t, 1, len(post.Comments))

	postRepo.EXPECT().GetPostById(postId).Return(post, nil)
	postRepo.EXPECT().Save(gomock.AssignableToTypeOf(&entity.Post{})).Do(
		func(arg *entity.Post) { post = arg },
	)

	req := usecase.NewDeletePostUseCaseReq(uuid.New(), postId, commentId)
	res := usecase.NewDeletePostUseCaseRes()
	uc := usecase.NewDeletePoseUseCase(userRepo, postRepo, req, res)

	uc.Execute()

	if res.Err == nil {
		t.Errorf("failed to execute usecase")
	}

	assert.Equal(t, "only the comment owner can delete this comment", res.Err.Error())
}
