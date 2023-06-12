package config

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	uuid "github.com/gofrs/uuid"

	hash "golang.org/x/crypto/bcrypt"
)

// get all category a topic have and write it in a long string
func GetCategory(r *http.Request) string {
	var CategoryName = []string{"Info", "Video Games", "Music", "Design", "Communication", "Animation3D", "NSFW", "Anime", "Manga"}
	var result = ""
	for _, v := range CategoryName {
		a := r.FormValue(v)
		b, _ := strconv.ParseBool(a)
		if b {
			result += v
			result += ","
		}
	}
	return result[:len(result)-1]
}

// Set in a database the information about a topic
func SetTopicInfo(state string, SetTopicsName string, SetTopicsDescription string, Category string, Uuid string) {
	db, err := sql.Open("sqlite3", "./Database/Topics.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if state == "CreateTopicInfo" {
		stmt, err := db.Prepare("insert into Topics_Name(Title, Description, Category, Creation_Date, Creator, Liker, Disliker) values(?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		stmt.Exec(SetTopicsName, SetTopicsDescription, Category, time.Now().Format(time.ANSIC), Uuid, "", "")
	}
}

// Set in a database the text written in a topic
func SetTopicText(IdTopics int, Uuid string, TopicText string, pp string) {
	db, err := sql.Open("sqlite3", "./Database/Topics.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt, err := db.Prepare("insert into Topics(Id, Uuid, Text, Written, Picture) values(?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	stmt.Exec(IdTopics, Uuid, TopicText, time.Now().Format(time.ANSIC), pp)
}

// give a unique uuid to a user
func GetUUID() uuid.UUID {
	// Creating UUID Version 4
	// panic on error
	return uuid.Must(uuid.NewV4())
}

// check if an account exist
func UserExists(db *sql.DB, uuid string) bool {
	sqlStmt := `SELECT Uuid FROM Accounts WHERE Uuid = ?`
	err := db.QueryRow(sqlStmt, uuid).Scan(&uuid)
	if err != nil {
		if err != sql.ErrNoRows {
			// a real error happened! you should change your function return
			// to "(bool, error)" and return "false, err" here
			log.Print(err)
		}
		fmt.Println("faux ", uuid)
		return false
	}
	fmt.Println("vrai ", uuid)
	return true
}

// check if an account exist
func EmailExists(db *sql.DB, Email string) bool {
	sqlStmt := `SELECT Email FROM Accounts WHERE Email = ?`
	err := db.QueryRow(sqlStmt, Email).Scan(&Email)
	if err != nil {
		if err != sql.ErrNoRows {
			// a real error happened! you should change your function return
			// to "(bool, error)" and return "false, err" here
			log.Print(err)
		}
		fmt.Println("faux ", Email)
		return false
	}
	fmt.Println("vrai ", Email)
	return true
}

// crypt password into hash
func HashPassword(passwd string) []byte {
	has, _ := hash.GenerateFromPassword([]byte(passwd), DefaultCost)
	return has
}

// compare given hash and password
func CheckPasswordHash(password string, hashpass string) bool {
	var err error = hash.CompareHashAndPassword([]byte(hashpass), []byte(password))
	return err == nil
}

const MAX_UPLOAD_SIZE = 1024 * 1024 // 1MB

// Progress is used to track the progress of a file upload.
// It implements the io.Writer interface so it can be passed
// to an io.TeeReader()
type Progress struct {
	TotalSize int64
	BytesRead int64
}

// Write is used to satisfy the io.Writer interface.
// Instead of writing somewhere, it simply aggregates
// the total bytes on each read
func (pr *Progress) Write(p []byte) (n int, err error) {
	n, err = len(p), nil
	pr.BytesRead += int64(n)
	//pr.Print()
	return
}

// set the number of like in a topic
func SetLiker(IdTopics int, UUID string, Like int, Liker string) {
	db, err := sql.Open("sqlite3", "./Database/Topics.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	stmt, err := db.Prepare("update Topics_Name set Liker = ?, Like = ? where Id = ?")
	if err != nil {
		log.Fatal(err)
	}
	var liked, result = CheckUuid(Liker, UUID)
	if liked {
		stmt.Exec(result, Like+1, IdTopics)
	} else if !liked {
		stmt.Exec(result, Like-1, IdTopics)
	}
}

// set the number of like in a topic
func SetDisliker(IdTopics int, UUID string, Like int, Liker string) {
	db, err := sql.Open("sqlite3", "./Database/Topics.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	stmt, err := db.Prepare("update Topics_Name set Disliker = ?, Like = ? where Id = ?")
	if err != nil {
		log.Fatal(err)
	}
	var liked, result = CheckUuid(Liker, UUID)
	fmt.Println(result)
	if liked {
		stmt.Exec(result, Like-1, IdTopics)
	} else if !liked {
		stmt.Exec(result, Like+1, IdTopics)
	}
}

// check if uuid is in string
func CheckUuid(str string, UUID string) (bool, string) {
	if str == "" {
		str = UUID
	} else {
		var strarr = strings.Split(str, ",")
		for index, v := range strarr {
			if UUID == v {
				fmt.Println(strarr)
				strarr = RemoveIndex(strarr, index)
				str = strings.Join(strarr, ",")
				return false, str
			}
		}
		str += ","
		str += UUID
	}
	return true, str
}

// remove element in array with index
func RemoveIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}

func SetLikerint(Liker string, Disliker string, UUID string) int {
	b, _ := CheckUuid(Liker, UUID)
	c, _ := CheckUuid(Disliker, UUID)
	if b && c {
		return 0
	} else if c && !b {
		return 1
	}
	return 2
}

// read database/store value from database to go code
func GetName(UUID string) string {
	db, err := sql.Open("sqlite3", "./Database/User.db")
	if err != nil {
		log.Fatal(err)
	}
	sql_readall := `SELECT Name, Uuid FROM Accounts`

	rows, err := db.Query(sql_readall)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var result Account
	for rows.Next() {
		rows.Scan(&result.Name, &result.Uuid)
		if UUID == result.Uuid.String() {
			return result.Name
		}
	}
	return ""
}

// set the number of like in a topic
func UpdateAccount(UUID string, Name string, Role string) {
	db, err := sql.Open("sqlite3", "./Database/User.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	stmt, err := db.Prepare("update Accounts set Name = ?, Role = ? where Uuid = ?")
	if err != nil {
		log.Fatal(err)
	}
	stmt.Exec(Name, Role, UUID)
}

// set the number of like in a topic
func DeleteComment(TimeStamp string) {
	db, err := sql.Open("sqlite3", "./Database/Topics.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.Exec("delete from Topics where Written = ?", TimeStamp)
}
func Verifmail(email string) bool {
	apiKey := "2783ca62-8461-42a4-a05e-816779066859"

	url := "https://isitarealemail.com/api/email/validate?email=" + url.QueryEscape(email)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "bearer "+apiKey)

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println(err)
		return false
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("error %v", err)
		return false
	}
	if res.StatusCode != 200 {
		fmt.Printf("unexpected result, check your api key. %v", res.Status)
		return false
	}

	var myJson RealEmailResponse
	json.Unmarshal(body, &myJson)

	fmt.Printf("status for %v is %v", email, myJson.Status)
	return myJson.Status == "valid"
}
