package repository

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/threadpulse/internal/config"
	"github.com/threadpulse/models"
)

type AuthRepo struct {
	DB *sqlx.DB
}

func NewAuthRepo() *AuthRepo {
	return &AuthRepo{
		DB: config.PostgisDB,
	}

}

func (r *AuthRepo) RegisterNewUserRepo(userRegister *models.User) error {

	query := `insert into users (name,email,hashed_pass) values($1,$2,$3)`
	_, err := r.DB.Exec(query, userRegister.Name, userRegister.Email, userRegister.HashPassword)
	if err != nil {
		return err
	}
	return nil

}

func (r *AuthRepo) VerifyByEmail(email string) (*models.User, error) {
	var user models.User
	query := `select * from users where email=$1`
	err := r.DB.Get(&user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return &user, nil
}
