package add_comment_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"mashu.example/internal/entity"
	entity_enums "mashu.example/internal/entity/enums"
	"mashu.example/internal/usecase/comment/add_comment"
	"mashu.example/internal/usecase/repository/mock"
)

func setup(t *testing.T) (*mock.MockPostRepo, *mock.MockUserRepo) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	return mock.NewMockPostRepo(mockCtrl), mock.NewMockUserRepo(mockCtrl)
}

func TestAddCommentUnderMyOwnPost(t *testing.T) {
	postRepo, userRepo := setup(t)

	ownerId := uuid.New()
	commentOwner := entity.NewUser(ownerId, "comment_owner", "comment owner display name", "comment_owner@email.com", true)

	postId := uuid.New()
	post := entity.NewPost(
		postId,
		"My First Post",
		"My first content",
		commentOwner,
		entity_enums.POST_PUBLIC,
	)

	postRepo.EXPECT().GetPostById(postId).Return(post, nil)
	userRepo.EXPECT().GetUserById(ownerId).Return(commentOwner, nil)
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

	assert.Nil(t, res.Err)
	assert.Equal(t, "Good!", post.Comments[0].Content)
	assert.Equal(t, post.ID, post.Comments[0].Post.ID)
	assert.Equal(t, commentOwner.ID, post.Comments[0].Owner.ID)
}

func TestAddMultipleCommentUnderPost(t *testing.T) {
	postRepo, userRepo := setup(t)

	commentOwner := entity.NewUser(uuid.New(), "comment_owner", "comment owner", "comment_owner@email.com", false)
	postOwner := entity.NewUser(uuid.New(), "post_owner", "post owner", "post_owner@email.com", true)
	post := entity.NewPost(uuid.New(), "Learning Domain Driven Design", "...", postOwner, entity_enums.POST_PUBLIC)

	// first comment
	userRepo.EXPECT().GetUserById(commentOwner.ID).Return(commentOwner, nil)
	postRepo.EXPECT().GetPostById(post.ID).Return(post, nil)
	postRepo.EXPECT().Save(gomock.AssignableToTypeOf(&entity.Post{})).Do(
		func(arg *entity.Post) { post = arg },
	)

	req := add_comment.NewAddCommentUseCaseReq(commentOwner.ID, post.ID, "good article!")
	res := add_comment.NewAddCommentUseCaseRes()
	uc := add_comment.NewAddCommentUseCase(userRepo, postRepo, req, res)

	uc.Execute()
	assert.Nil(t, res.Err)
	assert.Equal(t, 1, len(post.Comments))
	assert.Equal(t, "good article!", post.Comments[0].Content)
	assert.Equal(t, commentOwner.ID, post.Comments[0].Owner.ID)

	// second comment
	userRepo.EXPECT().GetUserById(postOwner.ID).Return(postOwner, nil)
	postRepo.EXPECT().GetPostById(post.ID).Return(post, nil)
	postRepo.EXPECT().Save(gomock.AssignableToTypeOf(&entity.Post{})).Do(
		func(arg *entity.Post) { post = arg },
	)

	req = add_comment.NewAddCommentUseCaseReq(postOwner.ID, post.ID, "thanks!")
	res = add_comment.NewAddCommentUseCaseRes()
	uc = add_comment.NewAddCommentUseCase(userRepo, postRepo, req, res)

	uc.Execute()
	assert.Nil(t, res.Err)
	assert.Equal(t, 2, len(post.Comments))
	assert.Equal(t, "thanks!", post.Comments[1].Content)
	assert.Equal(t, postOwner.ID, post.Comments[1].Owner.ID)
}

func TestAddCommentUnderPublicPost(t *testing.T) {
	postRepo, userRepo := setup(t)

	commentOwner := entity.NewUser(uuid.New(), "comment_owner", "comment owner", "comment_owner@email.com", false)
	postOwner := entity.NewUser(uuid.New(), "post_owner", "post owner", "post_owner@email.com", false)
	post := entity.NewPost(uuid.New(), "Learning Domain Driven Design", "...", postOwner, entity_enums.POST_PUBLIC)

	userRepo.EXPECT().GetUserById(commentOwner.ID).Return(commentOwner, nil)
	postRepo.EXPECT().GetPostById(post.ID).Return(post, nil)
	postRepo.EXPECT().Save(gomock.AssignableToTypeOf(&entity.Post{})).Do(
		func(arg *entity.Post) { post = arg },
	)

	req := add_comment.NewAddCommentUseCaseReq(commentOwner.ID, post.ID, "good article!")
	res := add_comment.NewAddCommentUseCaseRes()
	uc := add_comment.NewAddCommentUseCase(userRepo, postRepo, req, res)

	uc.Execute()
	assert.Nil(t, res.Err)
	assert.Equal(t, 1, len(post.Comments))
	assert.Equal(t, "good article!", post.Comments[0].Content)
	assert.Equal(t, commentOwner.ID, post.Comments[0].Owner.ID)
}

func TestAddCommentUnderFollowerOnlyPostWithoutFollow(t *testing.T) {
	postRepo, userRepo := setup(t)

	commentOwner := entity.NewUser(uuid.New(), "comment_owner", "comment owner", "comment_owner@email.com", false)
	postOwner := entity.NewUser(uuid.New(), "post_owner", "post owner", "post_owner@email.com", false)

	post := entity.NewPost(uuid.New(), "Learning Domain Driven Design", "...", postOwner, entity_enums.POST_FOLLOWER_ONLY)

	userRepo.EXPECT().GetUserById(commentOwner.ID).Return(commentOwner, nil)
	postRepo.EXPECT().GetPostById(post.ID).Return(post, nil)

	req := add_comment.NewAddCommentUseCaseReq(commentOwner.ID, post.ID, "good article!")
	res := add_comment.NewAddCommentUseCaseRes()
	uc := add_comment.NewAddCommentUseCase(userRepo, postRepo, req, res)

	uc.Execute()
	assert.Error(t, res.Err)
	assert.Equal(t, 0, len(post.Comments))
}

func TestAddCommentUnderFollowerOnlyPostWithFollow(t *testing.T) {
	postRepo, userRepo := setup(t)

	commentOwner := entity.NewUser(uuid.New(), "comment_owner", "comment owner", "comment_owner@email.com", false)
	postOwner := entity.NewUser(uuid.New(), "post_owner", "post owner", "post_owner@email.com", false)

	// comment owner follow post owner
	postOwner.Followers = append(postOwner.Followers, commentOwner.ID)

	post := entity.NewPost(uuid.New(), "Learning Domain Driven Design", "...", postOwner, entity_enums.POST_FOLLOWER_ONLY)

	userRepo.EXPECT().GetUserById(commentOwner.ID).Return(commentOwner, nil)
	postRepo.EXPECT().GetPostById(post.ID).Return(post, nil)
	postRepo.EXPECT().Save(gomock.AssignableToTypeOf(&entity.Post{})).Do(
		func(arg *entity.Post) { post = arg },
	)

	req := add_comment.NewAddCommentUseCaseReq(commentOwner.ID, post.ID, "good article!")
	res := add_comment.NewAddCommentUseCaseRes()
	uc := add_comment.NewAddCommentUseCase(userRepo, postRepo, req, res)

	uc.Execute()
	assert.Nil(t, res.Err)
	assert.Equal(t, 1, len(post.Comments))
	assert.Equal(t, "good article!", post.Comments[0].Content)
	assert.Equal(t, commentOwner.ID, post.Comments[0].Owner.ID)
}

func TestAddCommentUnderPrivatePost(t *testing.T) {
	postRepo, userRepo := setup(t)

	commentOwner := entity.NewUser(uuid.New(), "comment_owner", "comment owner", "comment_owner@email.com", false)
	postOwner := entity.NewUser(uuid.New(), "post_owner", "post owner", "post_owner@email.com", false)

	// comment owner follow post owner
	postOwner.Followers = append(postOwner.Followers, commentOwner.ID)

	post := entity.NewPost(uuid.New(), "Learning Domain Driven Design", "...", postOwner, entity_enums.POST_PRIVATE)

	userRepo.EXPECT().GetUserById(commentOwner.ID).Return(commentOwner, nil)
	postRepo.EXPECT().GetPostById(post.ID).Return(post, nil)

	req := add_comment.NewAddCommentUseCaseReq(commentOwner.ID, post.ID, "good article!")
	res := add_comment.NewAddCommentUseCaseRes()
	uc := add_comment.NewAddCommentUseCase(userRepo, postRepo, req, res)

	uc.Execute()
	assert.Error(t, res.Err)
	assert.Equal(t, 0, len(post.Comments))
}
