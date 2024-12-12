package port

import "prior-chat-bot/internal/adapter/api/model"

type MyRepo interface {
	FindUserByEmail(email string) (model.UserLoginRequest, error)
	FindUserById(userId float64) (model.UserLoginModel, error)
	SignUp(request model.UserSignUpRequest, password string) error
}
