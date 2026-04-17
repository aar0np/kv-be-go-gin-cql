package repository

import (
	"fmt"
	"killrvideo/go-backend-astra-cql/models"
	"time"

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

func (r *AuthDAL) ExistsByEmail(email string) bool {
	user, err := r.GetUserByEmail(email)
	if user != nil && err == nil {
		return true
	}
	return false
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

func (r *AuthDAL) UpdateUser(user models.User) {
	var email bool
	var firstname bool
	var lastname bool
	var accountStatus bool

	strCQL := `UPDATE users SET `

	if user.Email != "" {
		strCQL += `email=?`
		email = true
	}
	if user.FirstName != "" {
		if strCQL != "UPDATE users SET " {
			strCQL += ","
		}
		strCQL += `firstname=?`
		firstname = true
	}
	if user.LastName != "" {
		if strCQL != "UPDATE users SET " {
			strCQL += ","
		}
		strCQL += `lastname=?`
		lastname = true
	}
	if user.AccountStatus != "" {
		if strCQL != "UPDATE users SET " {
			strCQL += ","
		}
		strCQL += `account_status=?`
		accountStatus = true
	}

	strCQL += ` WHERE userid=?`

	if email && firstname && lastname && accountStatus {
		r.DB.Query(strCQL, user.Email, user.FirstName, user.LastName, user.AccountStatus, user.Userid).Exec()
	} else if email && firstname && lastname {
		r.DB.Query(strCQL, user.Email, user.FirstName, user.LastName, user.Userid).Exec()
	} else if email && firstname {
		r.DB.Query(strCQL, user.Email, user.FirstName, user.Userid).Exec()
	} else if email && lastname {
		r.DB.Query(strCQL, user.Email, user.LastName, user.Userid).Exec()
	} else if firstname && lastname {
		r.DB.Query(strCQL, user.FirstName, user.LastName, user.Userid).Exec()
	} else if email {
		r.DB.Query(strCQL, user.Email, user.Userid).Exec()
	} else if firstname {
		r.DB.Query(strCQL, user.FirstName, user.Userid).Exec()
	} else if lastname {
		r.DB.Query(strCQL, user.LastName, user.Userid).Exec()
	} else if accountStatus {
		r.DB.Query(strCQL, user.AccountStatus, user.Userid).Exec()
	} else if firstname && accountStatus {
		r.DB.Query(strCQL, user.FirstName, user.AccountStatus, user.Userid).Exec()
	} else if lastname && accountStatus {
		r.DB.Query(strCQL, user.LastName, user.AccountStatus, user.Userid).Exec()
	} else if firstname && lastname && accountStatus {
		r.DB.Query(strCQL, user.FirstName, user.LastName, user.AccountStatus, user.Userid).Exec()
	} else if email && accountStatus {
		r.DB.Query(strCQL, user.Email, user.AccountStatus, user.Userid).Exec()
	} else if email && firstname && accountStatus {
		r.DB.Query(strCQL, user.Email, user.FirstName, user.AccountStatus, user.Userid).Exec()
	} else if email && lastname && accountStatus {
		r.DB.Query(strCQL, user.Email, user.LastName, user.AccountStatus, user.Userid).Exec()
	}
}

func (r *AuthDAL) UpdatePassword(userCreds models.UserCredentials) {
	strCQL := `UPDATE user_credentials SET password=? WHERE email=?`
	r.DB.Query(strCQL, userCreds.Password, userCreds.Email).Exec()
}

func (r *AuthDAL) DeleteUserCreds(email string) {
	strCQL := `DELETE FROM user_credentials WHERE email=?`
	r.DB.Query(strCQL, email).Exec()
}

func (r *AuthDAL) RegisterLogin(userid apachegocql.UUID) {
	strCQL := `UPDATE users SET last_login_date=? WHERE userid=?`
	r.DB.Query(strCQL, time.Now, userid).Exec()
}
