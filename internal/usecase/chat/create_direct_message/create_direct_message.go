package create_direct_message

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
	chat "mashu.example/internal/entity/chat"
	"mashu.example/internal/usecase"
	"mashu.example/internal/usecase/repository"
)

var (
	ErrChatRoomAlreadyExist             = errors.New("chatroom already exist")
	ErrSenderDoNotFollowPrivateReceiver = errors.New("sender don't follow the private receiver")
)

type CreateDirectMessageUseCaseReq struct {
	senderId   uuid.UUID
	receiverId uuid.UUID
}

type CreateDirectMessageUseCaseRes struct {
	DirectMessageId uuid.UUID
	Err             error
}

type CreateDirectMessageUseCase struct {
	chatRepo repository.ChatRepo
	userRepo repository.UserRepo
	req      *CreateDirectMessageUseCaseReq
	res      *CreateDirectMessageUseCaseRes
}

func (uc *CreateDirectMessageUseCase) Execute() {
	sender, err := uc.userRepo.GetUserById(uc.req.senderId)
	if err != nil {
		errMsg := fmt.Sprintf("user %s not exist", uc.req.senderId)
		logrus.Error(errMsg)
		uc.res.Err = errors.New(errMsg)
		return
	}

	receiver, err := uc.userRepo.GetUserById(uc.req.receiverId)
	if err != nil {
		errMsg := fmt.Sprintf("user %s not exist", uc.req.receiverId)
		logrus.Error(errMsg)
		uc.res.Err = errors.New(errMsg)
		return
	}

	dm, err := uc.chatRepo.GetDMByUserId(sender.ID, receiver.ID)
	if dm != nil {
		uc.res.Err = ErrChatRoomAlreadyExist
		logrus.Info(ErrChatRoomAlreadyExist.Error())
		return
	}
	if _, ok := err.(*repository.ErrDMNotFound); !ok {
		uc.res.Err = err
		logrus.Error(err)
		return
	}

	if !receiver.Public && !slices.Contains(receiver.Followers, sender.ID) {
		uc.res.Err = ErrSenderDoNotFollowPrivateReceiver
		logrus.Error(ErrSenderDoNotFollowPrivateReceiver.Error())
		return
	}

	id := uuid.New()
	dm = chat.NewDirectMessage(id, sender, receiver)
	if err := uc.chatRepo.SaveDirectMessage(dm); err != nil {
		uc.res.Err = err
		return
	}
	uc.res.DirectMessageId = id
}

func NewCreateDirectMessageUseCase(
	chatRepo repository.ChatRepo,
	userRepo repository.UserRepo,
	req *CreateDirectMessageUseCaseReq,
	res *CreateDirectMessageUseCaseRes,
) usecase.UseCase {
	return &CreateDirectMessageUseCase{chatRepo, userRepo, req, res}
}

func NewCreateDirectMessageUseCaseReq(
	senderId uuid.UUID,
	receiverId uuid.UUID,
) *CreateDirectMessageUseCaseReq {
	return &CreateDirectMessageUseCaseReq{senderId, receiverId}
}

func NewCreateDirectMessageUseCaseRes() *CreateDirectMessageUseCaseRes {
	return &CreateDirectMessageUseCaseRes{}
}
