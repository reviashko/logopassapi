package models

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq" // here
)

//ConnectionData struct
type ConnectionData struct {
	Host     string
	Port     int
	Dbname   string
	User     string
	Password string
}

//ToString function
func (c *ConnectionData) ToString() string {
	return fmt.Sprintf("host=%s port=%v dbname=%s user=%s password=%s sslmode=disable", c.Host, strconv.Itoa(c.Port), c.Dbname, c.User, c.Password)
}

//Datastore interface
type Datastore interface {
	SaveUser(*UserData) (int, pq.ErrorCode, error)
	GetUser(int) (*UserData, pq.ErrorCode, error)
	GetUserByAuth(string, []byte) (*UserData, pq.ErrorCode, error)
	GetUserByEmail(email string) (*UserData, pq.ErrorCode, error)
}

//DB struct
type DB struct {
	*sqlx.DB
}

// InitDB function
func InitDB(connectionData ConnectionData) (*DB, error) {

	db, err := sqlx.Connect("postgres", connectionData.ToString())
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	//connection pool default settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err = db.Ping(); err != nil {
		log.Panic(err)
		return nil, err
	}

	return &DB{db}, nil
}
