package send_message

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	chat "mashu.example/internal/entity/chat"
	"mashu.example/internal/usecase"
	"mashu.example/internal/usecase/repository"
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
		errMsg := "sender not exist"
		logrus.Error(errMsg)
		uc.res.Err = errors.New(errMsg)
		return
	}
	receiver, err := uc.userRepo.GetUserById(uc.req.receiverId)
	if err != nil {
		errMsg := "receiver not exist"
		logrus.Error(errMsg)
		uc.res.Err = errors.New(errMsg)
		return
	}

	dm, err := uc.chatRepo.GetDMByUserId(sender.ID, receiver.ID)
	if err != nil {
		errMsg := "chatroom does not exist"
		logrus.Error(errMsg)
		uc.res.Err = errors.New(errMsg)
		return
	}

	dm.Messages = append(dm.Messages, chat.NewMessage(
		uuid.New(),
		sender,
		uc.req.message,
	))

	if err := uc.chatRepo.SaveDirectMessage(dm); err != nil {
		uc.res.Err = err
		return
	}
}

func NewSendMessageUseCase() usecase.UseCase {
	return &SendMessageUseCase{}
}

func NewSendMessageUseCaseReq() *SendMessageUseCaseReq {
	return &SendMessageUseCaseReq{}
}

func NewSendMessageUseCaseRes() *SendMessageUseCaseRes {
	return &SendMessageUseCaseRes{}
}
