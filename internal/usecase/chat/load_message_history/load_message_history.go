package load_message_history

import (
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	entity "mashu.example/internal/entity/chat"
	"mashu.example/internal/usecase"
	"mashu.example/internal/usecase/repository"
)

// DTO for entity Message
type MessageDTO struct {
	ID        uuid.UUID
	OwnerId   uuid.UUID
	Content   string
	Timestamp time.Time
}

func NewMessageDTO(m entity.Message) MessageDTO {
	return MessageDTO{
		ID:        m.ID,
		OwnerId:   m.OwnerId,
		Content:   m.Content,
		Timestamp: m.Timestamp,
	}
}

type LoadMessageHistoryUseCaseReq struct {
	userId uuid.UUID
}

type LoadMessageHistoryUseCaseRes struct {
	MessageMap map[string][]MessageDTO // map the dm id to a list of message
	Err        error
}

type LoadMessageHistoryUseCase struct {
	userRepo repository.UserRepo
	chatRepo repository.ChatRepo
	req      *LoadMessageHistoryUseCaseReq
	res      *LoadMessageHistoryUseCaseRes
}

func (uc *LoadMessageHistoryUseCase) Execute() {
	user, err := uc.userRepo.GetUserById(uc.req.userId)
	if err != nil {
		uc.res.Err = &repository.ErrUserNotFound{UserId: uc.req.userId}
		logrus.Error(uc.res.Err)
		return
	}

	dms, err := uc.chatRepo.GetDMsByPartUserId(user.ID)
	if err != nil {
		uc.res.Err = err
		logrus.Error(err)
		return
	}

	for _, dm := range dms {
		msgDTOs := []MessageDTO{}
		for _, msg := range dm.Messages {
			msgDTOs = append(msgDTOs, NewMessageDTO(*msg))
		}
		uc.res.MessageMap[dm.ID.String()] = msgDTOs
	}

	uc.res.Err = nil
}

func NewLoadMessageHistoryUseCase(
	userRepo repository.UserRepo,
	chatRepo repository.ChatRepo,
	req *LoadMessageHistoryUseCaseReq,
	res *LoadMessageHistoryUseCaseRes,
) usecase.UseCase {
	return &LoadMessageHistoryUseCase{userRepo, chatRepo, req, res}
}

func NewLoadMessageHistoryUseCaseReq(
	userId uuid.UUID,
) *LoadMessageHistoryUseCaseReq {
	return &LoadMessageHistoryUseCaseReq{userId}
}

func NewLoadMessageHistoryUseCaseRes() *LoadMessageHistoryUseCaseRes {
	return &LoadMessageHistoryUseCaseRes{
		MessageMap: map[string][]MessageDTO{},
	}
}
