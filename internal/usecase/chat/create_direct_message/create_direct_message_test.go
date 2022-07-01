package create_direct_message_test

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"mashu.example/internal/entity"
	chat "mashu.example/internal/entity/chat"
	usecase "mashu.example/internal/usecase/chat/create_direct_message"
	"mashu.example/internal/usecase/repository"
	"mashu.example/internal/usecase/tests"
)

func TestCreateDirectMessageToPublicUser(t *testing.T) {
	userRepo, _, _, chatRepo := tests.SetupTestRepositories(t)

	sender := entity.NewUser(uuid.New(), "sender", "Sender", "sender@email.com", true)
	receiver := entity.NewUser(uuid.New(), "receiver", "Receiver", "receiver@email.com", true)

	userRepo.EXPECT().GetUserById(sender.ID).Return(sender, nil)
	userRepo.EXPECT().GetUserById(receiver.ID).Return(receiver, nil)
	chatRepo.EXPECT().GetDMByUserId(sender.ID, receiver.ID).Return(nil, &repository.ErrDMNotFound{})

	var createdDm *chat.DirectMessage
	chatRepo.
		EXPECT().
		SaveDirectMessage(gomock.AssignableToTypeOf(&chat.DirectMessage{})).
		Do(func(arg *chat.DirectMessage) { createdDm = arg })

	req := usecase.NewCreateDirectMessageUseCaseReq(sender.ID, receiver.ID)
	res := usecase.NewCreateDirectMessageUseCaseRes()
	uc := usecase.NewCreateDirectMessageUseCase(chatRepo, userRepo, req, res)

	uc.Execute()

	assert.Nil(t, res.Err)
	assert.Equal(t, sender.ID, createdDm.Creator.ID)
	assert.Equal(t, receiver.ID, createdDm.Receiver.ID)
	assert.Empty(t, createdDm.Messages)
}

func TestCreateDirectMessageToPrivateUserWithoutFollowing(t *testing.T) {
	userRepo, _, _, chatRepo := tests.SetupTestRepositories(t)

	sender := entity.NewUser(uuid.New(), "sender", "Sender", "sender@email.com", true)
	receiver := entity.NewUser(uuid.New(), "receiver", "Receiver", "receiver@email.com", false)

	userRepo.EXPECT().GetUserById(sender.ID).Return(sender, nil)
	userRepo.EXPECT().GetUserById(receiver.ID).Return(receiver, nil)
	chatRepo.EXPECT().GetDMByUserId(sender.ID, receiver.ID).Return(nil, &repository.ErrDMNotFound{})

	req := usecase.NewCreateDirectMessageUseCaseReq(sender.ID, receiver.ID)
	res := usecase.NewCreateDirectMessageUseCaseRes()
	uc := usecase.NewCreateDirectMessageUseCase(chatRepo, userRepo, req, res)

	uc.Execute()
	fmt.Println(res.Err.Error())

	assert.ErrorIs(t, res.Err, usecase.ErrSenderDoNotFollowPrivateReceiver)
	assert.Zero(t, res.DirectMessageId)
}

func TestCreateDirectMessageToPrivateUserAfterFollow(t *testing.T) {
	userRepo, _, _, chatRepo := tests.SetupTestRepositories(t)

	sender := entity.NewUser(uuid.New(), "sender", "Sender", "sender@email.com", true)
	receiver := entity.NewUser(uuid.New(), "receiver", "Receiver", "receiver@email.com", false)
	sender.Followings = append(sender.Followings, receiver.ID)
	receiver.Followers = append(receiver.Followers, sender.ID)

	userRepo.EXPECT().GetUserById(sender.ID).Return(sender, nil)
	userRepo.EXPECT().GetUserById(receiver.ID).Return(receiver, nil)
	chatRepo.EXPECT().GetDMByUserId(sender.ID, receiver.ID).Return(nil, &repository.ErrDMNotFound{})

	var createdDm *chat.DirectMessage
	chatRepo.
		EXPECT().
		SaveDirectMessage(gomock.AssignableToTypeOf(&chat.DirectMessage{})).
		Do(func(arg *chat.DirectMessage) { createdDm = arg })

	req := usecase.NewCreateDirectMessageUseCaseReq(sender.ID, receiver.ID)
	res := usecase.NewCreateDirectMessageUseCaseRes()
	uc := usecase.NewCreateDirectMessageUseCase(chatRepo, userRepo, req, res)

	uc.Execute()

	assert.Nil(t, res.Err)
	assert.Equal(t, sender.ID, createdDm.Creator.ID)
	assert.Equal(t, receiver.ID, createdDm.Receiver.ID)
}

func TestCreateDuplicateDirectMessage(t *testing.T) {
	userRepo, _, _, chatRepo := tests.SetupTestRepositories(t)

	user1 := entity.NewUser(uuid.New(), "user1", "User 1", "user1@email.com", true)
	user2 := entity.NewUser(uuid.New(), "user2", "User 2", "user2@email.com", true)

	userRepo.EXPECT().GetUserById(user1.ID).Return(user1, nil)
	userRepo.EXPECT().GetUserById(user2.ID).Return(user2, nil)

	dm := chat.NewDirectMessage(uuid.New(), user1, user2)
	chatRepo.EXPECT().GetDMByUserId(user1.ID, user2.ID).Return(dm, nil)

	req := usecase.NewCreateDirectMessageUseCaseReq(user1.ID, user2.ID)
	res := usecase.NewCreateDirectMessageUseCaseRes()
	uc := usecase.NewCreateDirectMessageUseCase(chatRepo, userRepo, req, res)

	uc.Execute()

	assert.ErrorIs(t, res.Err, usecase.ErrChatRoomAlreadyExist)
	assert.Zero(t, res.DirectMessageId)
}

func TestCreateDirectMessageButUserNotExist(t *testing.T) {
	userRepo, _, _, chatRepo := tests.SetupTestRepositories(t)

	user1 := entity.NewUser(uuid.New(), "user1", "User 1", "user1@email.com", true)
	user2Id := uuid.New()

	userRepo.EXPECT().GetUserById(user1.ID).Return(user1, nil)
	userRepo.EXPECT().GetUserById(user2Id).Return(nil, gorm.ErrRecordNotFound)

	req := usecase.NewCreateDirectMessageUseCaseReq(user1.ID, user2Id)
	res := usecase.NewCreateDirectMessageUseCaseRes()
	uc := usecase.NewCreateDirectMessageUseCase(chatRepo, userRepo, req, res)

	uc.Execute()

	assert.NotNil(t, res.Err)
	assert.Equal(t, res.Err.Error(), (&repository.ErrUserNotFound{UserId: user2Id}).Error())
}
