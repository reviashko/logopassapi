package models

import (
	"github.com/lib/pq"
)

// UserData structure
type UserData struct {
	UserID    int    `json:"user_id" db:"user_id"`
	IsActive  bool   `json:"is_active" db:"is_active"`
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
	Email     string `json:"email" db:"email"`
	PswdHashB []byte `json:"pswd_hash_bytes" db:"pswd_hash_bytes"`
}

// GetUserByAuth method
func (db *DB) GetUserByAuth(email string, pswdHashB []byte) (*UserData, pq.ErrorCode, error) {

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

// GetUserByEmail method
func (db *DB) GetUserByEmail(email string) (*UserData, pq.ErrorCode, error) {

	var errorCode pq.ErrorCode

	rows, err := db.Queryx("SELECT user_id, is_active, first_name, last_name, email from users.user_getByEmail($1)", email)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			errorCode = err.Code
		}
		return nil, errorCode, err
	}
	defer rows.Close()

	userList := make([]*UserData, 0)

	for rows.Next() {
		userData := new(UserData)
		err = rows.StructScan(&userData)
		if err != nil {
			return nil, errorCode, err
		}
		userList = append(userList, userData)
	}

	if len(userList) != 1 {
		return nil, errorCode, nil
	}
	return userList[0], errorCode, nil
}

// GetUser method
func (db *DB) GetUser(userID int) (*UserData, pq.ErrorCode, error) {

	var errorCode pq.ErrorCode

	rows, err := db.Queryx("SELECT user_id, is_active, first_name, last_name, email from users.user_get($1)", userID)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			errorCode = err.Code
		}
		return nil, errorCode, err
	}
	defer rows.Close()

	userList := make([]*UserData, 0)

	for rows.Next() {
		userData := new(UserData)
		err = rows.StructScan(&userData)
		if err != nil {
			return nil, errorCode, err
		}
		userList = append(userList, userData)
	}

	if len(userList) != 1 {
		return nil, errorCode, nil
	}
	return userList[0], errorCode, nil
}

// SaveUser method
func (db *DB) SaveUser(userData *UserData) (int, pq.ErrorCode, error) {

	var errorCode pq.ErrorCode

	rows, err := db.Queryx("select user_id from users.user_save($1, $2, $3, $4, $5, $6)", userData.UserID, userData.IsActive, userData.FirstName, userData.LastName, userData.Email, userData.PswdHashB)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			errorCode = err.Code
		}
		return 0, errorCode, err
	}
	defer rows.Close()

	userList := make([]*UserData, 0)

	for rows.Next() {
		uData := new(UserData)
		err = rows.StructScan(&uData)
		if err != nil {
			return 0, errorCode, err
		}
		userList = append(userList, uData)
	}

	if len(userList) != 1 {
		return 0, errorCode, nil
	}

	return userList[0].UserID, errorCode, nil
}
