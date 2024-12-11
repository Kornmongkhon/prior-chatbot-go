package repository

import (
	"database/sql"
	_ "github.com/lib/pq" // เพื่อให้การใช้ pq driver ทำงาน
	"log"
	"prior-chat-bot/internal/adapter/api/model"
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
