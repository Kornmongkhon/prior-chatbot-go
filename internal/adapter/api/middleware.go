package api

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"prior-chat-bot/internal/adapter/repository"
	"prior-chat-bot/internal/core/authentication"
	"prior-chat-bot/internal/core/domain"
	"strings"
)

var jwtProvider *authentication.JwtTokenProvider

func InterceptorFilter(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		repo := c.Get("repo")
		if repo == nil {
			return c.JSON(http.StatusInternalServerError, domain.ResponseT[interface{}]{
				Code:    "E9999",
				Message: "Database connection not initialized",
			})
		}
		authRepo, ok := repo.(*repository.AuthRepository)
		if !ok {
			log.Println("InterceptorFilter -> Error casting repo to *repository.AuthRepository")
			response := domain.ResponseT[interface{}]{
				Code:    "E9999",
				Message: "The system has a problem. Please contact the system administrator.",
			}
			return c.JSON(http.StatusInternalServerError, response)
		}
		//log.Println(c.Request().RequestURI)
		if isAlwaysAllowEndPoint(c.Request().RequestURI) {
			return next(c)
		}
		accessToken := c.Request().Header.Get("Authorization")
		accessToken = strings.TrimPrefix(accessToken, "Bearer ")
		log.Println("InterceptorFilter -> jwt:", accessToken)
		if accessToken == "" {
			response := domain.ResponseT[interface{}]{
				Code:    "E9998",
				Message: "Do not have permission to access this content. Please sign in.",
			}
			return c.JSON(http.StatusUnauthorized, response)
		} else if !jwtProvider.ValidateToken(accessToken) {
			log.Println("InterceptorFilter -> Invalid token")
			response := domain.ResponseT[interface{}]{
				Code:    "I0006",
				Message: "Access Token expired. Please sign in again.",
			}
			return c.JSON(http.StatusUnauthorized, response)
		}
		claims, err := jwtProvider.DecodeTokenClaims(accessToken)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, domain.ResponseT[interface{}]{
				Code:    "I0006",
				Message: "Access Token expired. Please sign in again.",
			})
		}

		jwtSubject := claims["sub"].(string)
		log.Println("InterceptorFilter -> jwtSubject:", jwtSubject)
		var userJsonObject map[string]interface{}
		err = json.Unmarshal([]byte(jwtSubject), &userJsonObject)
		if err != nil {
			log.Println("Error unmarshalling subject: %v", err)
		}
		email := userJsonObject["email"].(string)
		refreshToken, _ := userJsonObject["refreshToken"].(string)
		if refreshToken == "" || refreshToken == "Y" {
			log.Println("InterceptorFilter -> Missing header refresh token")
			response := domain.ResponseT[interface{}]{
				Code:    "E9999",
				Message: "The system has a problem. Please contact the system administrator.",
			}
			return c.JSON(http.StatusUnauthorized, response)
		}

		log.Printf("InterceptorFilter -> email: %s", email)
		if authRepo.DB == nil {
			log.Println("InterceptorFilter -> DB connection is nil")
			response := domain.ResponseT[interface{}]{
				Code:    "E9999",
				Message: "Internal Server Error: Database connection is nil",
			}
			return c.JSON(http.StatusInternalServerError, response)
		}
		_, err = authRepo.FindUserByEmail(email)
		if err != nil {
			log.Println("InterceptorFilter -> User not found")
			response := domain.ResponseT[interface{}]{
				Code:    "E9999",
				Message: "The system has a problem. Please contact the system administrator.",
			}
			return c.JSON(http.StatusUnauthorized, response)
		}
		return next(c)

	}
}

func isAlwaysAllowEndPoint(requestUrl string) bool {
	allowedEndpoints := []string{
		"/prior_chatbot_api/api/v1/health/check",
		"/prior_chatbot_api/api/v1/auth/sign-in",
		"/prior_chatbot_api/api/v1/auth/regenerate-tokens",
	}
	for _, endpoint := range allowedEndpoints {
		if requestUrl == endpoint {
			return true
		}
	}

	return false

}
