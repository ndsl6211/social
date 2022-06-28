package tests

import (
	"testing"

	"github.com/golang/mock/gomock"
	"mashu.example/internal/usecase/repository/mock"
)

func SetupTestRepositories(t *testing.T) (*mock.MockUserRepo, *mock.MockPostRepo, *mock.MockGroupRepo, *mock.MockChatRepo) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	return mock.NewMockUserRepo(mockCtrl), mock.NewMockPostRepo(mockCtrl), mock.NewMockGroupRepo(mockCtrl), mock.NewMockChatRepo(mockCtrl)
}
