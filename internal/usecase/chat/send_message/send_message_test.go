package send_message_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"mashu.example/internal/entity"
	chat "mashu.example/internal/entity/chat"
	usecase "mashu.example/internal/usecase/chat/send_message"
	"mashu.example/internal/usecase/tests"
)

func TestSendMessage(t *testing.T) {
	userRepo, _, _, chatRepo := tests.SetupTestRepositories(t)

	sender := entity.NewUser(uuid.New(), "sender", "Sender", "sender@email.com", true)
	receiver := entity.NewUser(uuid.New(), "receiver", "Receiver", "receiver@email.com", true)

	userRepo.EXPECT().GetUserById(sender.ID).Return(sender, nil)
	userRepo.EXPECT().GetUserById(receiver.ID).Return(receiver, nil)

	dmId := uuid.New()
	dm := chat.NewDirectMessage(dmId, sender, receiver)

	var updatedDM *chat.DirectMessage
	chatRepo.EXPECT().GetDMByUserId(sender.ID, receiver.ID).Return(dm, nil)
	chatRepo.
		EXPECT().
		SaveDirectMessage(gomock.AssignableToTypeOf(&chat.DirectMessage{})).
		Do(func(arg *chat.DirectMessage) { updatedDM = arg })

	now := time.Now()
	req := usecase.NewSendMessageUseCaseReq(
		sender.ID,
		receiver.ID,
		"Hi! How are you?",
		now,
	)
	res := usecase.NewSendMessageUseCaseRes()
	uc := usecase.NewSendMessageUseCase(userRepo, chatRepo, req, res)

	uc.Execute()

	assert.Nil(t, res.Err)
	assert.Len(t, updatedDM.Messages, 1)
	fmt.Println(updatedDM.Messages[0].Content)
	assert.Equal(t, "Hi! How are you?", updatedDM.Messages[0].Content)
	assert.Equal(t, now, updatedDM.Messages[0].Timestamp)
}
