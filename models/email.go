package models

import (
	"github.com/lib/pq"
)

// EmailTemplate structure
type EmailTemplate struct {
	UserID int `json:"user_id" db:"user_id"`
}

// GetUserByAuth22 method
func (db *DB) GetUserByAuth2(email string, pswdHashB []byte) (*UserData, pq.ErrorCode, error) {

	var errorCode pq.ErrorCode

	rows, err := db.Queryx("SELECT user_id, is_active, first_name, last_name, email from users.user_getByAuth($1, $2)", email, pswdHashB)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			errorCode = err.Code
		}
		return new(UserData), errorCode, err
	}
	defer rows.Close()

	userList := make([]*UserData, 0)

	for rows.Next() {
		userData := new(UserData)
		err = rows.StructScan(&userData)
		if err != nil {
			return new(UserData), errorCode, err
		}
		userList = append(userList, userData)
	}

	if len(userList) != 1 {
		return new(UserData), errorCode, nil
	}

	return userList[0], errorCode, nil
}
