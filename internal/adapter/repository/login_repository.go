package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"mashu.example/internal/usecase/repository"
)

type loginRepo struct {
	db           *redis.Client
	loginSessTTL int
}

func (lr *loginRepo) SaveUserSess(userId uuid.UUID) error {
	ctx := context.Background()

	if err := lr.db.Set(
		ctx,
		lr.getLoginSessKey(userId),
		userId,
		0,
	).Err(); err != nil {
		return err
	}

	return nil
}

func (lr *loginRepo) ClearUserSess(userId uuid.UUID) error {
	ctx := context.Background()

	if err := lr.db.Del(ctx, lr.getLoginSessKey(userId)).Err(); err != nil {
		return err
	}

	return nil
}

func (lr *loginRepo) getLoginSessKey(userId uuid.UUID) string {
	return fmt.Sprintf("sess:discord:login:%s", userId)
}

func NewLoginRepository(
	db *redis.Client,
	loginSessTTL int,
) repository.LoginRepo {
	return &loginRepo{db, loginSessTTL}
}
