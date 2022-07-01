package load_message_history_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"mashu.example/internal/entity"
	chat "mashu.example/internal/entity/chat"
	usecase "mashu.example/internal/usecase/chat/load_message_history"
	"mashu.example/internal/usecase/tests"
)

func TestLoadMessageHistory(t *testing.T) {
	userRepo, _, _, chatRepo := tests.SetupTestRepositories(t)

	user1 := entity.NewUser(uuid.New(), "user1", "User1", "user1@email.com", true)
	user2 := entity.NewUser(uuid.New(), "user2", "User2", "user2@email.com", true)
	user3 := entity.NewUser(uuid.New(), "user3", "User3", "user3@email.com", true)
	user4 := entity.NewUser(uuid.New(), "user4", "User4", "user4@email.com", true)
	dms := []*chat.DirectMessage{
		chat.NewDirectMessage(uuid.New(), user1, user2),
		chat.NewDirectMessage(uuid.New(), user1, user3),
		chat.NewDirectMessage(uuid.New(), user4, user1),
	}
	dms[0].Messages = append(dms[0].Messages, chat.NewMessage(uuid.New(), user1.ID, "hi, user2"))
	dms[0].Messages = append(dms[0].Messages, chat.NewMessage(uuid.New(), user2.ID, "hello, user1"))
	dms[1].Messages = append(dms[1].Messages, chat.NewMessage(uuid.New(), user1.ID, "excuse me?"))
	dms[2].Messages = append(dms[2].Messages, chat.NewMessage(uuid.New(), user4.ID, "are you available now?"))
	dms[2].Messages = append(dms[2].Messages, chat.NewMessage(uuid.New(), user1.ID, "of course, wuzzup?"))

	userRepo.EXPECT().GetUserById(user1.ID).Return(user1, nil)
	chatRepo.EXPECT().GetDMsByPartUserId(user1.ID).Return(dms, nil)

	req := usecase.NewLoadMessageHistoryUseCaseReq(user1.ID)
	res := usecase.NewLoadMessageHistoryUseCaseRes()
	uc := usecase.NewLoadMessageHistoryUseCase(userRepo, chatRepo, req, res)

	uc.Execute()

	assert.Nil(t, res.Err)
	assert.Len(t, res.MessageMap, 3)
	assert.Len(t, res.MessageMap[dms[0].ID.String()], 2)
	assert.Len(t, res.MessageMap[dms[1].ID.String()], 1)
	assert.Len(t, res.MessageMap[dms[2].ID.String()], 2)
}
