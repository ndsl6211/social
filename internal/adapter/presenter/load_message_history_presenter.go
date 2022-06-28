package presenter

import uc "mashu.example/internal/usecase/chat/load_message_history"

type MessageHistoryPresenter struct {
	res *uc.LoadMessageHistoryUseCaseRes
}

type MessageHistoryViewModel struct{}

func (mhp *MessageHistoryPresenter) BuildViewModel() MessageHistoryViewModel {
	return MessageHistoryViewModel{}
}

// constructor of message history presenter
func NewMessageHistoryPresenter(
	res *uc.LoadMessageHistoryUseCaseRes,
) Presenter[MessageHistoryViewModel] {
	return &MessageHistoryPresenter{res}
}
