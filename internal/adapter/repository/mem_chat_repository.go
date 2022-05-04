package repository

import (
	"github.com/google/uuid"
	"golang.org/x/exp/slices"
	chat "mashu.example/internal/entity/chat"
	"mashu.example/internal/usecase/repository"
)

type memChatRepo struct {
	directMessages []chat.DirectMessage
}

func (mct *memChatRepo) GetDirectMessage(dmId uuid.UUID) (*chat.DirectMessage, error) {
	idx := slices.IndexFunc(mct.directMessages, func(dm chat.DirectMessage) bool {
		return dm.ID == dmId
	})

	if idx == -1 {
		return nil, &repository.ErrDMNotFound{}
	}

	return &mct.directMessages[idx], nil
}

func (mct *memChatRepo) GetDMByUserId(userA uuid.UUID, userB uuid.UUID) (*chat.DirectMessage, error) {
	idx := slices.IndexFunc(mct.directMessages, func(dm chat.DirectMessage) bool {
		return (dm.Creator.ID == userA && dm.Receiver.ID == userB) || (dm.Creator.ID == userB && dm.Receiver.ID == userA)
	})

	if idx == -1 {
		return nil, &repository.ErrDMNotFound{}
	}

	return &mct.directMessages[idx], nil
}

func (mct *memChatRepo) SaveDirectMessage(dm *chat.DirectMessage) error {
	idx := slices.IndexFunc(mct.directMessages, func(d chat.DirectMessage) bool {
		return d.ID == dm.ID
	})
	if idx != -1 {
		slices.Delete(mct.directMessages, idx, idx+1)
	}
	mct.directMessages = append(mct.directMessages, *dm)

	return nil
}

func NewMemChatRepository() repository.ChatRepo {
	return &memChatRepo{}
}
