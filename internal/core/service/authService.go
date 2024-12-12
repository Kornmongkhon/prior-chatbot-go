package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"prior-chat-bot/internal/adapter/api/model"
	"prior-chat-bot/internal/core/authentication"
	"prior-chat-bot/internal/core/domain"
	"prior-chat-bot/internal/core/port"
	"regexp"
	"strings"
	"time"
)

type AuthService struct {
	repo             port.MyRepo
	jwtTokenProvider *authentication.JwtTokenProvider // pointer type
}

func NewAuthService(repo port.MyRepo, jwtTokenProvider *authentication.JwtTokenProvider) *AuthService {
	return &AuthService{
		repo:             repo,
		jwtTokenProvider: jwtTokenProvider, // passing pointer
	}
}

func (s *AuthService) SignIn(request model.UserLoginRequest) (int, domain.ResponseT[interface{}]) {
	log.Println("AuthService -> signIn")
	errors := validateEmailAndPassword(request.Email, request.Password)
	if len(errors) > 0 {
		errorMessage := fmt.Sprintf("Invalid Request require %s", strings.Join(errors, ", "))
		return http.StatusBadRequest, domain.ResponseT[interface{}]{
			Code:    "I0001",
			Message: errorMessage,
		}
	}
	httpStatus, emailValidationResponse := validateEmailPattern(request.Email)
	if httpStatus != http.StatusOK {
		return httpStatus, emailValidationResponse
	}
	user, err := s.repo.FindUserByEmail(request.Email)
	if err != nil {
		log.Println("Error finding user:", err)
		return http.StatusBadRequest, domain.ResponseT[interface{}]{
			Code:    "I0001",
			Message: "User not found",
		}
	}
	// validate match password
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
	log.Println("AuthService -> claims")

	jwtSubject, ok := claims["sub"].(string)
	if !ok {
		log.Println("AuthService -> jwtSubject not found")
		return http.StatusUnauthorized, domain.ResponseT[interface{}]{
			Code:    "E9998",
			Message: "Do not have permission to access this content. Please sign in.",
		}
	}
	log.Println("AuthService -> jwtSubject:", jwtSubject)
	var jwtSubjectData map[string]interface{}
	err = json.Unmarshal([]byte(jwtSubject), &jwtSubjectData)
	if err != nil {
		log.Println("Error unmarshalling jwtSubject:", err)
		return http.StatusUnauthorized, domain.ResponseT[interface{}]{
			Code:    "E9998",
			Message: "Do not have permission to access this content. Please sign in.",
		}
	}
	//log.Println("AuthService -> jwtSubjectData:", jwtSubjectData)
	userId, err := s.repo.FindUserById(jwtSubjectData["userId"].(float64))
	if err != nil {
		log.Println("Error finding user:", err)
		return http.StatusNotFound, domain.ResponseT[interface{}]{
			Code:    "I0001",
			Message: "User not found",
		}
	}
	return http.StatusOK, domain.ResponseT[interface{}]{
		Code:    "S0000",
		Message: "Success",
		Data:    userId,
	}
}

func (s *AuthService) RegenerateToken(refreshToken string) (int, domain.ResponseT[interface{}]) {
	log.Println("AuthService -> regenerateToken")
	if !s.jwtTokenProvider.ValidateToken(refreshToken) {
		return http.StatusUnauthorized, domain.ResponseT[interface{}]{
			Code:    "I0007",
			Message: "Refresh Token expired. Please sign in again.",
		}
	}
	claims, err := s.jwtTokenProvider.DecodeTokenClaims(refreshToken)
	if err != nil {
		log.Println("Error decoding refreshToken:", err)
		return http.StatusUnauthorized, domain.ResponseT[interface{}]{
			Code:    "E9999",
			Message: "The system has a problem. Please contact the system administrator.",
		}
	}
	jwtSubject, ok := claims["sub"].(string)
	if !ok {
		log.Println("AuthService -> jwtSubject not found")
		return http.StatusUnauthorized, domain.ResponseT[interface{}]{
			Code:    "E9998",
			Message: "Do not have permission to access this content. Please sign in.",
		}
	}
	var jwtSubjectData map[string]interface{}
	err = json.Unmarshal([]byte(jwtSubject), &jwtSubjectData)
	if err != nil {
		log.Println("Error unmarshalling jwtSubject:", err)
		return http.StatusUnauthorized, domain.ResponseT[interface{}]{
			Code:    "E9998",
			Message: "Do not have permission to access this content. Please sign in.",
		}
	}
	email, err := s.repo.FindUserByEmail(jwtSubjectData["email"].(string))
	if err != nil {
		log.Println("Error finding email:", err)
		return http.StatusNotFound, domain.ResponseT[interface{}]{
			Code:    "I0001",
			Message: "Email not found",
		}
	}
	jwtResponse, err := authentication.Map(email, s.jwtTokenProvider)
	if err != nil {
		log.Println("Error mapping jwtResponse:", err)
		return http.StatusInternalServerError, domain.ResponseT[interface{}]{
			Code:    "E9999",
			Message: "The system has a problem. Please contact the system administrator.",
		}
	}
	return http.StatusOK, domain.ResponseT[interface{}]{
		Code:    "S0000",
		Message: "Success",
		Data:    jwtResponse,
	}
}

func (s *AuthService) SignUp(request model.UserSignUpRequest) (int, domain.ResponseT[interface{}]) {
	errors := validateSignUpRequest(request)
	if len(errors) > 0 {
		errorMessage := fmt.Sprintf("Invalid Request require %s", strings.Join(errors, ", "))
		return http.StatusBadRequest, domain.ResponseT[interface{}]{
			Code:    "I0001",
			Message: errorMessage,
		}
	}
	httpStatus, emailValidationResponse := validateEmailPattern(request.Email)
	if httpStatus != http.StatusOK {
		return httpStatus, emailValidationResponse
	}
	httpStatus, passwordValidationResponse := validatePasswordMatch(request.Password, request.ConfirmPassword)
	if httpStatus != http.StatusOK {
		return httpStatus, passwordValidationResponse
	}
	hashedPassword, err := authentication.HashPassword(request.Password)
	if err != nil {
		log.Println("Error hashing password:", err)
		return http.StatusInternalServerError, domain.ResponseT[interface{}]{
			Code:    "E9999",
			Message: "The system has a problem. Please contact the system administrator.",
		}
	}
	s.repo.SignUp(request, hashedPassword)
	return http.StatusOK, domain.ResponseT[interface{}]{
		Code:    "S0000",
		Message: "Success",
	}
}

func validatePasswordMatch(password string, password2 string) (int, domain.ResponseT[interface{}]) {
	if password != password2 {
		return http.StatusBadRequest, domain.ResponseT[interface{}]{
			Code:    "I0001",
			Message: "Passwords do not match. Please try again.",
		}
	}
	return http.StatusOK, domain.ResponseT[interface{}]{}
}

func validateEmailAndPassword(email, password string) []string {
	errors := []string{}
	if email == "" {
		errors = append(errors, "email")
	}
	if password == "" {
		errors = append(errors, "password")
	}
	return errors
}

func validateSignUpRequest(request model.UserSignUpRequest) []string {
	errors := []string{}
	if request.Email == "" {
		errors = append(errors, "email")
	}
	if request.Password == "" {
		errors = append(errors, "password")
	}
	if request.ConfirmPassword == "" {
		errors = append(errors, "confirmPassword")
	}
	if request.Mobile == "" {
		errors = append(errors, "mobile")
	}
	if request.Dob == "" {
		errors = append(errors, "dob")
	} else {
		valid, message := validateDob(request.Dob)
		if !valid {
			errors = append(errors, message)
		}
	}
	if request.Sex == "" {
		errors = append(errors, "sex")
	}
	return errors
}

func validateDob(dob string) (bool, string) {
	patternDate := "2006/01/02"

	parsedDob, err := time.Parse(patternDate, dob)
	if err != nil {
		return false, "Invalid date format. Please use yyyy/MM/dd."
	}
	if parsedDob.After(time.Now()) {
		return false, "Date of birth cannot be in the future."
	}

	// Check if the user is older than 18 years (optional validation)
	//ageLimit := 18
	//currentYear := time.Now().Year()
	//yearOfBirth := parsedDob.Year()
	//
	// Check age
	//if currentYear-yearOfBirth < ageLimit {
	//	return false, "You must be at least 18 years old."
	//}

	// Date is valid
	return true, "Valid date of birth."
}
func validateEmailPattern(email string) (int, domain.ResponseT[interface{}]) {
	// validate email pattern
	regexPattern := `^[A-Za-z0-9_-]+(\.[A-Za-z0-9_-]+)*@[A-Za-z0-9-]+(\.[A-Za-z0-9-]+)*(\.[A-Za-z]{2,})$`
	regex, err := regexp.Compile(regexPattern)
	if err != nil {
		log.Println("Invalid regex pattern:", err)
		return http.StatusInternalServerError, domain.ResponseT[interface{}]{
			Code:    "E9999",
			Message: "The system has a problem. Please contact the system administrator.",
		}
	}
	matched := regex.MatchString(email)
	if !matched {
		return http.StatusBadRequest, domain.ResponseT[interface{}]{
			Code:    "I0004",
			Message: "Invalid email Please try again.",
		}
	}
	return http.StatusOK, domain.ResponseT[interface{}]{}
}
