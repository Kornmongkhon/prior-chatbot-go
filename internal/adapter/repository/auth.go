package repository

import (
	"database/sql"
	_ "github.com/lib/pq" // เพื่อให้การใช้ pq driver ทำงาน
	"log"
	"prior-chat-bot/configs"
	"prior-chat-bot/internal/adapter/api/model"
	"time"
)

type AuthRepository struct {
	DB *sql.DB
}

func (repo *AuthRepository) FindUserByEmail(email string) (model.UserLoginRequest, error) {
	var user model.UserLoginRequest
	findUserQuery := `SELECT user_id, email, password FROM pc_ti_user WHERE email = $1 AND is_deleted = 'N'`
	row := repo.DB.QueryRow(findUserQuery, email)
	err := row.Scan(&user.UserId, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("User not found:", email)
		} else {
			log.Println("Error querying database:", err)
		}
		return user, err
	}
	return user, nil

}

func (repo *AuthRepository) FindUserById(userId float64) (model.UserLoginModel, error) {
	var user model.UserLoginModel
	findUserQuery := `SELECT user_id, email, mobile, dob, sex
						FROM pc_ti_user WHERE user_id = $1 AND is_deleted = 'N'`
	row := repo.DB.QueryRow(findUserQuery, userId)
	err := row.Scan(&user.UserId, &user.Email, &user.Mobile, &user.Dob, &user.Sex)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("User not found:", userId)
		} else {
			log.Println("Error querying database:", err)
		}
		return user, err
	}
	return user, nil
}
func (repo *AuthRepository) SignUp(request model.UserSignUpRequest, password string) error {
	tx, err := repo.DB.Begin()
	if err != nil {
		log.Println("Error starting transaction:", err)
		return err
	}
	signUpQuery := `INSERT INTO pc_ti_user (email, password, mobile, dob, sex, created_at )
					VALUES ($1, $2, $3, $4, $5, $6)`
	currentTime := time.Now()
	formattedTime := configs.FormatTime(currentTime)
	_, err = tx.Exec(signUpQuery, request.Email, password, request.Mobile, request.Dob, request.Sex, formattedTime)
	if err != nil {
		tx.Rollback()
		log.Println("Error inserting user:", err)
		return err
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
