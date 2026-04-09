package repository

import (
	"fmt"
	"killrvideo/go-backend-astra-cql/models"

	apachegocql "github.com/apache/cassandra-gocql-driver/v2"
)

type AuthDAL struct {
	DB *apachegocql.Session
}

func NewAuthDAL(session *apachegocql.Session) *AuthDAL {
	return &AuthDAL{
		DB: session,
	}
}

func (r *AuthDAL) GetUserById(id apachegocql.UUID) (*models.User, error) {
	user := &models.User{Userid: id}

	err1 := r.DB.Query(
		"SELECT userid, account_status, created_date, email, firstname, lastname, last_login_date FROM users WHERE userid = ?",
		id,
	).Scan(&user.Userid, &user.AccountStatus, &user.CreatedDate, &user.Email, &user.FirstName, &user.LastName, &user.LastLoginDate)

	if err1 != nil {
		return nil, fmt.Errorf("query has failed: %w", err1)
	}

	return user, nil
}

func (r *AuthDAL) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{Email: email}

	err1 := r.DB.Query(
		"SELECT userid, account_status, created_date, email, firstname, lastname, last_login_date FROM users WHERE email = ?",
		email,
	).Scan(&user.Userid, &user.AccountStatus, &user.CreatedDate, &user.Email, &user.FirstName, &user.LastName, &user.LastLoginDate)

	if err1 != nil {
		return nil, fmt.Errorf("query has failed: %w", err1)
	}

	return user, nil
}

func (r *AuthDAL) GetUserCredsByEmail(email string) (*models.UserCredentials, error) {
	user := &models.UserCredentials{Email: email}

	err1 := r.DB.Query(
		"SELECT userid, email, password, account_locked FROM user_credentials WHERE email = ?",
		email,
	).Scan(&user.Userid, &user.Email, &user.Password, &user.AccountLocked)

	if err1 != nil {
		return nil, fmt.Errorf("query has failed: %w", err1)
	}

	return user, nil
}

func (r *AuthDAL) SaveUser(user models.User) {
	r.DB.Query(`INSERT INTO users (userid, email, firstname, lastname, account_status, created_date, last_login_date) VALUES(?,?,?,?)`,
		user.Userid, user.Email).Exec()
}

func (r *AuthDAL) SaveUserCreds(userCreds models.UserCredentials) {
	r.DB.Query(`INSERT INTO user_credentials (userid, email, password, account_locked) VALUES(?,?,?,?)`,
		userCreds.Userid, userCreds.Email, userCreds.Password, userCreds.AccountLocked).Exec()
}
