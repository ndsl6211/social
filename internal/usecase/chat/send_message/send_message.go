package chat

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	chat "mashu.example/internal/entity/chat"
	"mashu.example/internal/usecase"
	"mashu.example/internal/usecase/repository"
)

var (
	ErrChatRoomNotExist  = errors.New("chatroom not exist")
	ErrSaveMessageFailed = errors.New("failed to save message")
)

type SendMessageUseCaseReq struct {
	senderId   uuid.UUID
	receiverId uuid.UUID
	message    string
	timestamp  time.Time
}

type SendMessageUseCaseRes struct {
	Err error
}

type SendMessageUseCase struct {
	userRepo repository.UserRepo
	chatRepo repository.ChatRepo
	req      *SendMessageUseCaseReq
	res      *SendMessageUseCaseRes
}

func (uc *SendMessageUseCase) Execute() {
	sender, err := uc.userRepo.GetUserById(uc.req.senderId)
	if err != nil {
		uc.res.Err = &repository.ErrUserNotFound{UserId: uc.req.senderId}
		logrus.Error(uc.res.Err)
		return
	}
	receiver, err := uc.userRepo.GetUserById(uc.req.receiverId)
	if err != nil {
		uc.res.Err = &repository.ErrUserNotFound{UserId: uc.req.senderId}
		logrus.Error(uc.res.Err)
		return
	}

	dm, err := uc.chatRepo.GetDMByUserId(sender.ID, receiver.ID)
	if err != nil {
		uc.res.Err = ErrChatRoomNotExist
		logrus.Error(ErrChatRoomNotExist.Error())
		return
	}

	dm.Messages = append(dm.Messages, chat.NewMessageWithTime(
		uuid.New(),
		sender.ID,
		uc.req.message,
		uc.req.timestamp,
	))

	if err := uc.chatRepo.SaveDirectMessage(dm); err != nil {
		uc.res.Err = ErrSaveMessageFailed
		logrus.Error(uc.res.Err)
		return
	}
}

func NewSendMessageUseCase(
	userRepo repository.UserRepo,
	chatRepo repository.ChatRepo,
	req *SendMessageUseCaseReq,
	res *SendMessageUseCaseRes,
) usecase.UseCase {
	return &SendMessageUseCase{userRepo, chatRepo, req, res}
}

func NewSendMessageUseCaseReq(
	senderId uuid.UUID,
	receiverId uuid.UUID,
	message string,
	timestamp time.Time,
) *SendMessageUseCaseReq {
	return &SendMessageUseCaseReq{senderId, receiverId, message, timestamp}
}

func NewSendMessageUseCaseRes() *SendMessageUseCaseRes {
	return &SendMessageUseCaseRes{}
}
