package repository

import (
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	chat "mashu.example/internal/entity/chat"
	"mashu.example/internal/usecase/repository"
)

type redisChatRepo struct {
	db *redis.Client
}

func (cr *redisChatRepo) GetDirectMessage(dmId uuid.UUID) (*chat.DirectMessage, error) {
	return nil, nil
}

func (cr *redisChatRepo) GetDMByUserId(
	userA uuid.UUID,
	userB uuid.UUID,
) (*chat.DirectMessage, error) {
	return nil, nil
}

func (cr *redisChatRepo) SaveDirectMessage(dm *chat.DirectMessage) error {
	return nil
}

func NewRedisChatRepository(db *redis.Client) repository.ChatRepo {
	return &redisChatRepo{db}
}
