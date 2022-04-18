package send_message_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"mashu.example/internal/usecase/repository/mock"
)

func setup(t *testing.T) (*mock.MockUserRepo, *mock.MockChatRepo) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	return mock.NewMockUserRepo(mockCtrl), mock.NewMockChatRepo(mockCtrl)
}

func TestSendMessage(t *testing.T) {

	// userRepo, chatRepo := setup(t)

	// sender := entity.NewUser(uuid.New(), "sender", "Sender", "sender@email.com", true)
	// sender := entity.NewUser(uuid.New(), "sender", "Sender", "sender@email.com", true)

	// userRepo.EXPECT().GetUserById()
}
