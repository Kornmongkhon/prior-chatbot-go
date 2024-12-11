package authentication

import (
	"encoding/json"
	"errors"
	"log"
	"prior-chat-bot/internal/adapter/api/model"
)

type JwtTokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func Map(user model.UserLoginRequest, jwtTokenProvider *JwtTokenProvider) (*JwtTokenResponse, error) {
	tokenMapAccess := map[string]interface{}{
		"userId":       user.UserId,
		"email":        user.Email,
		"refreshToken": "N",
	}
	tokenMapRefresh := map[string]interface{}{
		"userId":       user.UserId,
		"email":        user.Email,
		"refreshToken": "Y",
	}

	// Convert to JSON
	accessTokenPayload, err := json.Marshal(tokenMapAccess)
	if err != nil {
		log.Println("Error marshalling accessTokenPayload:", err)
		return nil, err
	}

	refreshTokenPayload, err := json.Marshal(tokenMapRefresh)
	if err != nil {
		log.Println("Error marshalling refreshTokenPayload:", err)
		return nil, err
	}

	// Generate tokens
	accessToken, err := jwtTokenProvider.GenerateToken(string(accessTokenPayload))
	if err != nil {
		log.Println("Error generating access token:", err)
		return nil, err
	}

	refreshToken, err := jwtTokenProvider.GenerateRefreshToken(string(refreshTokenPayload))
	if err != nil {
		log.Println("Error generating refresh token:", err)
		return nil, err
	}

	// Return response
	return &JwtTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func ObjectToJsonString(object interface{}) (string, error) {
	jsonData, err := json.Marshal(object)
	if err != nil {
		return "", errors.New("error converting object to JSON string")
	}
	return string(jsonData), nil
}
