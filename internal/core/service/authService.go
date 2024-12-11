package service

import (
	"fmt"
	"log"
	"net/http"
	"prior-chat-bot/internal/adapter/api/model"
	"prior-chat-bot/internal/core/authentication"
	"prior-chat-bot/internal/core/domain"
	"prior-chat-bot/internal/core/port"
	"strings"
)

type AuthService struct {
	repo             port.MyRepo
	jwtTokenProvider *authentication.JwtTokenProvider // pointer type
}

func NewAuthService(repo port.MyRepo, jwtTokenProvider *authentication.JwtTokenProvider) *AuthService {
	return &AuthService{
		repo:             repo,
		jwtTokenProvider: jwtTokenProvider, // passing pointer here
	}
}

func (s *AuthService) SignIn(request model.UserLoginRequest) (int, domain.ResponseT[interface{}]) {
	log.Println("AuthService -> signIn")
	errors := validateUserNameAndPassword(request.Email, request.Password)
	if len(errors) > 0 {
		errorMessage := fmt.Sprintf("Invalid Request require %s", strings.Join(errors, ", "))
		return http.StatusBadRequest, domain.ResponseT[interface{}]{
			Code:    "I0001",
			Message: errorMessage,
		}
	}
	user, err := s.repo.FindUserByEmail(request.Email)
	if err != nil {
		log.Println("Error finding user:", err)
		return http.StatusBadRequest, domain.ResponseT[interface{}]{
			Code:    "I0001",
			Message: "User not found",
		}
	}
	// ตรวจสอบ password และอื่น ๆ
	if !authentication.MatchPassword(user.Password, request.Password) {
		return http.StatusUnauthorized, domain.ResponseT[interface{}]{
			Code:    "I0002",
			Message: "Incorrect username or password",
		}
	}
	jwtResponse, err := authentication.Map(user, s.jwtTokenProvider)
	return http.StatusOK, domain.ResponseT[interface{}]{
		Code:    "S0000",
		Message: "Success",
		Data:    jwtResponse,
	}
}

func (s *AuthService) Me(accessToken string) (int, domain.ResponseT[interface{}]) {
	log.Println("AuthService -> me")
	claims, err := s.jwtTokenProvider.DecodeTokenClaims(accessToken)
	if err != nil {
		return 0, domain.ResponseT[interface{}]{}
	}
	log.Println("AuthService -> claims:", claims)
	userId, ok := claims["userId"].(string) // สมมติว่า userId เป็น string
	if !ok {
		log.Println("AuthService -> userId not found in claims")
		return http.StatusUnauthorized, domain.ResponseT[interface{}]{
			Code:    "I0007",
			Message: "User ID not found in token.",
		}
	}
	log.Println("AuthService -> userId:", userId)
	return http.StatusOK, domain.ResponseT[interface{}]{
		Code:    "S0000",
		Message: "Success",
	}
}

func validateUserNameAndPassword(email, password string) []string {
	errors := []string{}
	if email == "" {
		errors = append(errors, "email")
	}
	if password == "" {
		errors = append(errors, "password")
	}
	return errors
}
