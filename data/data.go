package data

import (
	"crypto/sha1"
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

var Db *sql.DB

func init() {
	var err error
	Db, err = sql.Open("postgres", "user=postgres password=1 dbname=Forum sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
}

type DisplayData struct{
	User User
	SingleThread Thread
	AllThreads []Thread
}

func CreateUUID() (client_uuid string) {
	client_uuid = uuid.New().String()
	return client_uuid
}

func Encrypt(plaintext string) (cryptext string) {
	cryptext = fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
	return
}
