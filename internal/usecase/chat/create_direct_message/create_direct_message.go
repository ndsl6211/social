package create_direct_message

import (
	"errors"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
	entity "mashu.example/internal/entity/chat"
	"mashu.example/internal/usecase"
	"mashu.example/internal/usecase/repository"
)

var (
	ErrChatRoomAlreadyExist             = errors.New("chatroom already exist")
	ErrSenderDoNotFollowPrivateReceiver = errors.New("sender don't follow the private receiver")
)

type CreateDirectMessageUseCaseReq struct {
	SenderId   uuid.UUID
	ReceiverId uuid.UUID
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
	sender, err := uc.userRepo.GetUserById(uc.req.SenderId)
	if err != nil {
		uc.res.Err = &repository.ErrUserNotFound{UserId: uc.req.SenderId}
		logrus.Error(uc.res.Err)
		return
	}

	receiver, err := uc.userRepo.GetUserById(uc.req.ReceiverId)
	if err != nil {
		uc.res.Err = &repository.ErrUserNotFound{UserId: uc.req.ReceiverId}
		logrus.Error(uc.res.Err)
		return
	}

	dm, err := uc.chatRepo.GetDMByUserId(sender.ID, receiver.ID)
	if dm != nil {
		uc.res.Err = ErrChatRoomAlreadyExist
		logrus.Error(ErrChatRoomAlreadyExist.Error())
		return
	}
	if _, ok := err.(*repository.ErrDMNotFound); !ok {
		uc.res.Err = err
		logrus.Warn(err)
		return
	}

	if !receiver.Public && !slices.Contains(receiver.Followers, sender.ID) {
		uc.res.Err = ErrSenderDoNotFollowPrivateReceiver
		logrus.Error(ErrSenderDoNotFollowPrivateReceiver.Error())
		return
	}

	id := uuid.New()
	dm = entity.NewDirectMessage(id, sender, receiver)
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
