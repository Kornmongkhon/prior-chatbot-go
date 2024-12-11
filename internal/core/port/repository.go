package port

import "prior-chat-bot/internal/adapter/api/model"

type MyRepo interface {
	FindUserByEmail(email string) (model.UserLoginRequest, error)
}
