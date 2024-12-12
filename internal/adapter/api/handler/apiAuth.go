package handler

import (
	"fmt"
	"log"
	"net/http"
	"prior-chat-bot/configs"
	"prior-chat-bot/internal/adapter/api/model"
	"prior-chat-bot/internal/core/authentication"
	"prior-chat-bot/internal/core/port"
	"prior-chat-bot/internal/core/service"
	"strings"
)

func ExecuteHandlerSignIn(c port.MyServer, cfg *configs.Config) {
	fmt.Println("ExecuteHandlerHealthCheck")
	jwtTokenProvider := authentication.NewJwtTokenProvider(cfg.JWT.Secret, cfg.JWT.ExpirationAccessToken, cfg.JWT.ExpirationRefreshToken)

	authService := service.NewAuthService(c.GetRepo(), jwtTokenProvider)
	var request model.UserLoginRequest
	err := c.BindRequest(&request) // Bind request
	if err != nil {
		c.ToResponse(http.StatusBadRequest, "I0001", "Invalid request", nil)
		return
	}

	httpStatus, response := authService.SignIn(request)
	c.ToResponse(httpStatus, response.Code, response.Message, response.Data)
}

func ExecuteHandlerMe(c port.MyServer, cfg *configs.Config) {
	fmt.Println("ExecuteHandlerMe")
	jwtTokenProvider := authentication.NewJwtTokenProvider(cfg.JWT.Secret, cfg.JWT.ExpirationAccessToken, cfg.JWT.ExpirationRefreshToken)
	authService := service.NewAuthService(c.GetRepo(), jwtTokenProvider)
	authorizationHeader := c.GetRequest().Header.Get("Authorization")
	accessToken := strings.TrimPrefix(authorizationHeader, "Bearer ")
	log.Println("ExecuteHandlerMe -> jwt:", accessToken)
	httpStatus, response := authService.Me(accessToken)
	c.ToResponse(httpStatus, response.Code, response.Message, response.Data)
}

func ExecuteHandlerRegenerateToken(c port.MyServer, cfg *configs.Config) {
	fmt.Println("ExecuteHandlerRegenerateToken")
	jwtTokenProvider := authentication.NewJwtTokenProvider(cfg.JWT.Secret, cfg.JWT.ExpirationAccessToken, cfg.JWT.ExpirationRefreshToken)
	authService := service.NewAuthService(c.GetRepo(), jwtTokenProvider)
	refreshToken := c.GetRequest().Header.Get("Refresh-Token")
	httpStatus, response := authService.RegenerateToken(refreshToken)
	c.ToResponse(httpStatus, response.Code, response.Message, response.Data)

}

func ExecuteHandlerSignUp(c port.MyServer, cfg *configs.Config) {
	fmt.Println("ExecuteHandlerSignUp")
	jwtTokenProvider := authentication.NewJwtTokenProvider(cfg.JWT.Secret, cfg.JWT.ExpirationAccessToken, cfg.JWT.ExpirationRefreshToken)
	authService := service.NewAuthService(c.GetRepo(), jwtTokenProvider)
	var request model.UserSignUpRequest
	err := c.BindRequest(&request) // Bind request
	if err != nil {
		c.ToResponse(http.StatusBadRequest, "I0001", "Invalid request", nil)
		return
	}
	httpStatus, response := authService.SignUp(request)
	c.ToResponse(httpStatus, response.Code, response.Message, response.Data)

}
