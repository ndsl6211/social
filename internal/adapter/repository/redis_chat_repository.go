package repository

import (
	"github.com/redis/go-redis/v9"
	"github.com/google/uuid"
	chat "mashu.example/internal/entity/chat"
	"mashu.example/internal/usecase/repository"
)

type redisChatRepo struct {
	db *redis.Client
}

func (rcr *redisChatRepo) GetDirectMessage(dmId uuid.UUID) (*chat.DirectMessage, error) {
	return nil, nil
}

func (rcr *redisChatRepo) GetDMByUserId(
	userA uuid.UUID,
	userB uuid.UUID,
) (*chat.DirectMessage, error) {
	return nil, nil
}

func (rcr *redisChatRepo) GetDMsByPartUserId(userId uuid.UUID) ([]*chat.DirectMessage, error) {
	return nil, nil
}

func (rcr *redisChatRepo) SaveDirectMessage(dm *chat.DirectMessage) error {
	return nil
}

func NewRedisChatRepository(db *redis.Client) repository.ChatRepo {
	return &redisChatRepo{db}
}
