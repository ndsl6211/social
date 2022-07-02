package presenter

import uc "mashu.example/internal/usecase/chat/load_message_history"

type MessageHistoryPresenter struct {
	res *uc.LoadMessageHistoryUseCaseRes
}

type MessageHistoryViewModel struct {
	Messages map[string][]uc.MessageDTO
}

func (mhp *MessageHistoryPresenter) BuildViewModel() MessageHistoryViewModel {
	mhvm := MessageHistoryViewModel{}
	mhvm.Messages = mhp.res.MessageMap

	return mhvm
}

// constructor of message history presenter
func NewMessageHistoryPresenter(
	res *uc.LoadMessageHistoryUseCaseRes,
) Presenter[MessageHistoryViewModel] {
	return &MessageHistoryPresenter{res}
}
