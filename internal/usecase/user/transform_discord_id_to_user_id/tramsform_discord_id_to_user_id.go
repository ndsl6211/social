package transform_discord_id_to_user_id

import "mashu.example/internal/usecase/repository"

type TransformDiscordIdToUserIdUseCaseReq struct {
	discordUserId string
}

type TransformDiscordIdToUserIdUseCaseRes struct {
	UserId string
	Err    error
}

type TransformDiscordIdToUserIdUseCase struct {
	userRepo repository.UserRepo
	req      *TransformDiscordIdToUserIdUseCaseReq
	res      *TransformDiscordIdToUserIdUseCaseRes
}

func (uc *TransformDiscordIdToUserIdUseCase) Execute() {
	user, err := uc.userRepo.GetUserByDiscordUserId(uc.req.discordUserId)
	if err != nil {
		uc.res.UserId = ""
		uc.res.Err = err
		return
	}

	uc.res.UserId = user.ID.String()
}
