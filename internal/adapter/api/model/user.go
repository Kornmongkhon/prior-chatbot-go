package model

type UserLoginRequest struct {
	UserId   int    `json:"userId"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginModel struct {
	UserId int    `json:"userId"`
	Email  string `json:"email"`
	Mobile string `json:"mobile"`
	Dob    string `json:"dob"`
	Sex    string `json:"sex"`
}

type UserSignUpRequest struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
	Mobile          string `json:"mobile"`
	Dob             string `json:"dob"`
	Sex             string `json:"sex"`
}
