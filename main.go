package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"text/template"
	"time"

	"github.com/Alexiiisv/Project-Forum/v2/config"
)

var Data config.Account

var TName config.TName
var TContent config.TContent
var TopicsName config.Topics
var Dataarray config.AllAccount
var allTopics config.AllTopics
var UserActions config.UserActions
var Logged config.LoginYes
var state, Role, Name, Password, Email, TopicText, UUID, SetTopicsName, SetTopicsDescription, info, Category, pp_name, NewPassword1 string
var cookieonce, stateSingleTopics bool

var IdTopics, Likes int
var request *http.Request

func main() {
	cookieonce = true
	fmt.Println("Please connect to\u001b[31m localhost", config.LocalhostPort, "\u001b[0m")
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets")))) // Join Assets Directory to the server
	http.HandleFunc("/", index)
	http.HandleFunc("/accounts", ShowAccount)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/settings", Settings)
	http.HandleFunc("/logout", Logout)
	http.HandleFunc("/information", Info)
	http.HandleFunc("/register", Register)
	http.HandleFunc("/createAcc", CreateAccount)
	http.HandleFunc("/connect", LoggedOn)
	http.HandleFunc("/topics", AllTopics)
	http.HandleFunc("/singleTopics", singleTopics)
	http.HandleFunc("/DeletCom", Delcom)
	http.HandleFunc("/user_account", User_Info)
	http.HandleFunc("/updaccount", updaccount)
	http.HandleFunc("/updateaccount_by_user", UpdateAccountByUser)
	http.HandleFunc("/user_account_settings", User_Info_set)
	http.HandleFunc("/CreateTopicInfo", CreateTopicInfo)
	http.HandleFunc("/like", Like)
	http.HandleFunc("/upload_pp", uploadHandler)
	http.HandleFunc("/change_passwd", ChangePasswd)
	err := http.ListenAndServe(config.LocalhostPort, nil) // Set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// page account information
func Info(w http.ResponseWriter, r *http.Request) {
	Dataarray.Data = readuuid("LoggedOn")
	Logged.Account = Dataarray.Data[0]
	t := template.New("account-template")
	t = template.Must(t.ParseFiles("./tmpl/account.html", "./tmpl/header&footer.html"))
	t.ExecuteTemplate(w, "account", Logged)
}

func Settings(w http.ResponseWriter, r *http.Request) {
	Dataarray.Data = readuuid("LoggedOn")
	Logged.Account = Dataarray.Data[0]
	t := template.New("account-settings")
	t = template.Must(t.ParseFiles("./tmpl/account.html", "./tmpl/header&footer.html"))
	t.ExecuteTemplate(w, "settings", Logged)
}

// page index
func index(w http.ResponseWriter, r *http.Request) {
	request = r
	autoconnect()
	t := template.New("index-template")
	t = template.Must(t.ParseFiles("index.html", "./tmpl/header&footer.html"))
	t.ExecuteTemplate(w, "index", Logged)
}

// create account
func CreateAccount(w http.ResponseWriter, r *http.Request) {
	NameChoosen := r.FormValue("Name")
	Password := r.FormValue("Password")
	Email = r.FormValue("Email")
	if config.Verifmail(Email) {
		Dataarray.Data = append(Dataarray.Data, config.Account{Name: NameChoosen, Password: Password, Email: Email, Uuid: config.GetUUID()})
		saveUuid("accounts")
		ShowAccount(w, r)
	} else {
		Register(w, r)
	}
}

// create account
func updaccount(w http.ResponseWriter, r *http.Request) {
	NameUpdated := r.FormValue("Name")
	UUID = r.FormValue("Uuid")
	Role = r.FormValue("Role")
	fmt.Println(NameUpdated, Role, UUID)
	config.UpdateAccount(UUID, NameUpdated, Role)
	ShowAccount(w, r)
}

func UpdateAccountByUser(w http.ResponseWriter, r *http.Request) {
	NameUpdated := r.FormValue("Name")
	UUID = r.FormValue("Uuid")
	fmt.Println(NameUpdated, Role, UUID)
	config.UpdateAccount(UUID, NameUpdated, Role)
	ShowAccount(w, r)

}

func ChangePasswd(w http.ResponseWriter, r *http.Request) {
	state = r.FormValue("State")
	Password = r.FormValue("current_passwd")
	NewPassword1 = r.FormValue("new_passwd1")
	NewPassword2 := r.FormValue("new_passwd2")

	if NewPassword1 == NewPassword2 {
		saveUuid(state)
		Password = r.FormValue("new_passwd1")
	}
	http.Redirect(w, r, "/information", 301)
}

// page accounts
func ShowAccount(w http.ResponseWriter, r *http.Request) {
	request = r
	autoconnect()
	Dataarray.Data = readuuid("ShowAccount")
	Dataarray.Account = Logged
	t := template.New("account-template")
	t = template.Must(t.ParseFiles("./tmpl/account.html", "./tmpl/header&footer.html", "./tmpl/content.html"))
	t.ExecuteTemplate(w, "accounts", Dataarray)
}

// page accounts
func AllTopics(w http.ResponseWriter, r *http.Request) {
	request = r
	autoconnect()
	state := r.FormValue("State")
	if state == "CreateTopicInfo" {
		SetTopicsName = r.FormValue("Title")
		SetTopicsDescription = r.FormValue("Description")
		Category = config.GetCategory(r)
		config.SetTopicInfo(state, SetTopicsName, SetTopicsDescription, Category, Logged.Account.Uuid.String())
		http.Redirect(w, r, "/topics", 301)
	}
	allTopics.Name = readtopics()
	fmt.Println(allTopics.Name)
	if r.FormValue("ByLikeSub") == "ByLike" {
		fmt.Println("a")
		sort.SliceStable(allTopics.Name, func(i, j int) bool { return allTopics.Name[i].Like > allTopics.Name[j].Like })
		fmt.Println(allTopics.Name)
	}
	if r.FormValue("ByCreationDateSub") == "ByCreationDate" {
		sort.SliceStable(allTopics.Name, func(i, j int) bool {
			st, _ := time.Parse(time.ANSIC, allTopics.Name[i].CreationDate)
			nd, _ := time.Parse(time.ANSIC, allTopics.Name[j].CreationDate)
			fmt.Println(st, nd)
			return nd.Before(st)
		})
		fmt.Println(allTopics.Name)
	}
	allTopics.Login = Logged
	t := template.New("topics-template")
	t = template.Must(t.ParseFiles("./tmpl/topics.html", "./tmpl/header&footer.html", "./tmpl/content.html"))
	t.ExecuteTemplate(w, "topics", allTopics)
}

// page accounts
func CreateTopicInfo(w http.ResponseWriter, r *http.Request) {
	allTopics.Name = readtopics()
	allTopics.Login = Logged
	t := template.New("topics-template")
	t = template.Must(t.ParseFiles("./tmpl/topics.html", "./tmpl/header&footer.html", "./tmpl/content.html"))
	t.ExecuteTemplate(w, "CreateTopicInfo", allTopics)
}

// page accounts
func singleTopics(w http.ResponseWriter, r *http.Request) {
	state = r.FormValue("State")
	IdTopics, _ = strconv.Atoi(r.FormValue("IdTopics"))
	stateSingleTopics, _ = strconv.ParseBool(r.FormValue("StateBool"))
	if state == "PostTopic" {
		TopicText = r.FormValue("text")
		config.SetTopicText(IdTopics, Logged.Account.Uuid.String(), TopicText, "")
		urltest := "/singleTopics?IdTopics=" + strconv.Itoa(IdTopics) + "&State=SingleTopic"
		http.Redirect(w, r, urltest, 301)
	}
	TopicsName.Name = GetTopicsData()
	TopicsName.Name.Liked = config.SetLikerint(TopicsName.Name.Liker, TopicsName.Name.Disliker, UUID)
	TopicsName.Login = Logged
	UserActions.Account = readuuid("user_account")[0]
	fmt.Println(IdTopics)
	TopicsName.Content = GetTopicsContent()
	TopicsName.Accounts = readuuid("ShowAccount")
	if state == "SwitchMode" {
		TopicsName.Name.Pic = !TopicsName.Name.Pic
	}
	t := template.New("singleTopics-template")
	t = template.Must(t.ParseFiles("./tmpl/topics.html", "./tmpl/header&footer.html", "./tmpl/content.html"))
	t.ExecuteTemplate(w, "singleTopics", TopicsName)
}

func Delcom(w http.ResponseWriter, r *http.Request) {
	config.DeleteComment(r.FormValue("TimeStamps"))
	TopicsName.Content = GetTopicsContent()
	t := template.New("singleTopics-template")
	t = template.Must(t.ParseFiles("./tmpl/topics.html", "./tmpl/header&footer.html", "./tmpl/content.html"))
	t.ExecuteTemplate(w, "singleTopics", TopicsName)
}

// Display User Informations
func User_Info(w http.ResponseWriter, r *http.Request) {
	state = r.FormValue("state")

	if state == "user" {
		Name = r.FormValue("name")
	} else {
		UUID = r.FormValue("Uuid")
		UserActions.Login = Logged
	}
	UserActions.Account = readuuid(state)[0]
	UserActions.Commentaires = GetTopicsContent()
	if len(UserActions.Commentaires) > 5 {
		UserActions.Commentaires = UserActions.Commentaires[len(UserActions.Commentaires)-5:]
	}
	fmt.Println(state)
	t := template.New("account-template")
	t = template.Must(t.ParseFiles("./tmpl/account.html", "./tmpl/header&footer.html", "./tmpl/content.html"))
	t.ExecuteTemplate(w, "user_account", UserActions)
}

// Display User Informations
func User_Info_set(w http.ResponseWriter, r *http.Request) {
	t := template.New("account-template")
	t = template.Must(t.ParseFiles("./tmpl/account.html", "./tmpl/header&footer.html"))
	t.ExecuteTemplate(w, "user_account_settings", UserActions)
}

// page login
func Login(w http.ResponseWriter, r *http.Request) {
	request = r
	autoconnect()
	t := template.New("account-template")
	t = template.Must(t.ParseFiles("./tmpl/login&register.html", "./tmpl/header&footer.html"))
	if state == "login" {
		t.ExecuteTemplate(w, "login", config.Erreur{Connected: false, Miss: true})
		state = ""
	} else {

		t.ExecuteTemplate(w, "login", config.Erreur{Connected: false, Miss: false})
	}
}

// page register
func Register(w http.ResponseWriter, r *http.Request) {
	t := template.New("account-template")
	t = template.Must(t.ParseFiles("./tmpl/login&register.html", "./tmpl/header&footer.html"))
	t.ExecuteTemplate(w, "register", Dataarray)
}

// Connection to account
func LoggedOn(w http.ResponseWriter, r *http.Request) {
	Email = r.FormValue("Email")
	Password = r.FormValue("Password")
	Dataarray.Data = readuuid("LoggedOn")
	if len(Dataarray.Data) == 1 {
		Logged.Account = Dataarray.Data[0]
		Logged.Connected = true
		UUID = Dataarray.Data[0].Uuid.String()
		cookie(w, "Uuid", Dataarray.Data[0].Uuid.String(), 86400)
		http.Redirect(w, r, "/information", 301)

	} else {
		state = "login"
		Login(w, r)
	}
}

func cookie(w http.ResponseWriter, name string, value string, age int) {
	cookie := &http.Cookie{
		Name:   name,
		Value:  value,
		MaxAge: age,
	}
	fmt.Println(cookie)
	http.SetCookie(w, cookie)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie(w, "Uuid", getUuid(r), -1)
	Logged.Account = config.Account{}
	Logged.Connected = false
	Login(w, r)
}

func getUuid(r *http.Request) string {
	a, _ := r.Cookie("Uuid")
	return a.Value
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 32 MB is the default used by FormFile
	if err := r.ParseMultipartForm(20 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// get a reference to the fileHeaders
	files := r.MultipartForm.File["AddPP"]

	for _, fileHeader := range files {
		if fileHeader.Size > config.MAX_UPLOAD_SIZE {
			http.Error(w, fmt.Sprintf("The uploaded image is too big: %s. Please use an image less than 1MB in size", fileHeader.Filename), http.StatusBadRequest)
			return
		}

		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer file.Close()

		buff := make([]byte, 512)
		_, err = file.Read(buff)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		filetype := http.DetectContentType(buff)
		if filetype != "image/png" {
			http.Error(w, "The provided file format is not allowed. Please upload a JPEG or PNG image", http.StatusBadRequest)
			return
		}

		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var path = "./assets/image/" + r.FormValue("path")
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		pp_name = strconv.FormatInt(time.Now().UnixNano(), 10)
		f, err := os.Create(fmt.Sprintf("./assets/image/%s/%s%s", r.FormValue("path"), pp_name, filepath.Ext(fileHeader.Filename)))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		defer f.Close()

		pr := &config.Progress{
			TotalSize: fileHeader.Size,
		}

		_, err = io.Copy(f, io.TeeReader(file, pr))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	if r.FormValue("State") == "PostTopic" {
		var id, _ = strconv.Atoi(r.FormValue("IdTopics"))
		config.SetTopicText(id, Logged.Account.Uuid.String(), r.FormValue("text"), pp_name+".png")
		state = ""
		singleTopics(w, r)
	} else {
		state := "pp"
		saveUuid(state)
		Info(w, r)
	}
}

// read database/store value from database to go code
func readuuid(state string) []config.Account {
	db, err := sql.Open("sqlite3", "./Database/User.db")
	if err != nil {
		log.Fatal(err)
	}
	sql_readall := `SELECT Name, Password, Email, Uuid, Profile_Picture, Role FROM Accounts`

	rows, err := db.Query(sql_readall)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var result []config.Account
	for rows.Next() {
		rows.Scan(&Data.Name, &Data.Password, &Data.Email, &Data.Uuid, &Data.Profile_Picture, &Data.Role)
		if state == "LoggedOn" && config.CheckPasswordHash(Password, Data.Password) && Data.Email == Email {
			Data.Password = Password
			result = append(result, Data)
			break
		} else if state == "ShowAccount" {
			result = append(result, Data)
		} else if state == "user_account" && UUID == Data.Uuid.String() {
			result = append(result, Data)
			break
		} else if state == "user" && Name == Data.Name {
			result = append(result, Data)
			break
		}
	}
	return result
}

// read database/store value from database to go code
func readtopics() []config.TName {
	db, err := sql.Open("sqlite3", "./Database/Topics.db")
	if err != nil {
		log.Fatal(err)
	}
	sql_readall := `SELECT Id, Title, Description, Creation_Date, Category, Like, Creator FROM Topics_Name;`

	rows, err := db.Query(sql_readall)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var result []config.TName
	for rows.Next() {
		rows.Scan(&TName.Id, &TName.Title, &TName.Desc, &TName.CreationDate, &TName.Category, &TName.Like, &TName.Creator)
		TName.Creator = config.GetName(TName.Creator)
		result = append(result, TName)
	}
	return result
}

// Get from a database the information about a topic
func GetTopicsData() config.TName {
	db, err := sql.Open("sqlite3", "./Database/Topics.db")
	if err != nil {
		log.Fatal(err)
	}
	sql_readall := `SELECT Id, Title, Description, Category, Like, Liker, Disliker, Creator FROM Topics_Name;`

	rows, err := db.Query(sql_readall)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var result config.TName
	for rows.Next() {
		rows.Scan(&TName.Id, &TName.Title, &TName.Desc, &TName.Category, &TName.Like, &TName.Liker, &TName.Disliker, &TName.Creator)
		if TName.Id == IdTopics {
			TName.Pic = stateSingleTopics
			TName.Liked = 0
			result = TName
			break
		}
	}
	return result
}

// read database/store value from database to go code
func GetTopicsContent() []config.TContent {
	db, err := sql.Open("sqlite3", "./Database/Topics.db")
	if err != nil {
		log.Fatal(err)
	}
	sql_readall := `SELECT Id, Uuid, Text, Written, Picture FROM Topics`

	rows, err := db.Query(sql_readall)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var result []config.TContent
	Compte := readuuid("ShowAccount")
	for rows.Next() {
		if state == "SingleTopic" || state == "PostTopic" || state == "SwitchMode" {
			rows.Scan(&TContent.Id, &TContent.Uuid, &TContent.Text, &TContent.Written, &TContent.Picture)
			for i := 0; i < len(Compte); i++ {
				if TContent.Uuid == Compte[i].Uuid.String() {
					TContent.Name = Compte[i].Name
					break
				}
			}
			if TContent.Id == IdTopics {
				result = append(result, TContent)
			} else {
				continue
			}
		} else {
			rows.Scan(&TContent.Id, &TContent.Uuid, &TContent.Text, &TContent.Written, &TContent.Picture)
			if TContent.Uuid == UserActions.Account.Uuid.String() {
				TContent.Name = config.GetName(TContent.Uuid)
				result = append(result, TContent)
			}
		}

	}
	return result
}

// TODO: Why don't push into DB ?
// write in a database
func saveUuid(state string) {
	db, err := sql.Open("sqlite3", "./Database/User.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if state == "accounts" {
		stmt, err := db.Prepare("insert into Accounts(Name, Password, Email, Uuid, Profile_Picture) values(?, ?, ?, ?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		if !config.EmailExists(db, Email) {
			for index := range Dataarray.Data {
				if config.UserExists(db, Dataarray.Data[index].Uuid.String()) {
					continue
				}
				stmt.Exec(Dataarray.Data[index].Name, string(config.HashPassword(Dataarray.Data[len(Dataarray.Data)-1].Password)), Dataarray.Data[index].Email, Dataarray.Data[index].Uuid.String(), "Standard_Pic.png")
			}
		}
	} else if state == "pp" {
		if err != nil {
			panic(err)
		}
		stmt, err := db.Prepare("update Accounts set Profile_Picture = ? where Uuid = ?")
		if err != nil {
			log.Fatal(err)
		}
		if Logged.Account.Profile_Picture != "Standard_Pic.png" {
			toremove := "./assets/image/Account_pp/"
			toremove += Logged.Account.Profile_Picture
			os.Remove(toremove)
		}

		link := pp_name
		link += ".png"
		Logged.Account.Profile_Picture = link
		stmt.Exec(link, Logged.Account.Uuid.String())
	} else if state == "changepasswd" {
		stmt, _ := db.Prepare("update Accounts set Password = ? where Uuid = ?")
		NewPassword1 = string(config.HashPassword(NewPassword1))
		stmt.Exec(NewPassword1, Logged.Account.Uuid.String())

	}
}

func Like(w http.ResponseWriter, r *http.Request) {
	Likes, _ = strconv.Atoi(r.FormValue("Likes"))
	BtnStatus := r.FormValue("BtnStatus")
	fmt.Println(BtnStatus)
	if BtnStatus == "👎" {
		fmt.Println("a")
		Disliker := r.FormValue("Disliker")
		config.SetDisliker(IdTopics, UUID, Likes, Disliker)
	} else if BtnStatus == "👍" {
		Liker := r.FormValue("Liker")
		config.SetLiker(IdTopics, UUID, Likes, Liker)
	}
	TopicsName.Name = GetTopicsData()
	TopicsName.Name.Liked = config.SetLikerint(TopicsName.Name.Liker, TopicsName.Name.Disliker, UUID)
	t := template.New("like-template")
	t = template.Must(t.ParseFiles("./tmpl/topics.html", "./tmpl/header&footer.html", "./tmpl/content.html"))
	t.ExecuteTemplate(w, "singleTopics", TopicsName)
}

func autoconnect() {
	if cookieonce {
		cookieonce = false
		r := request
		fmt.Println("Cookies in API Call:")

		tokenCookie, err := r.Cookie("Uuid")
		if err != nil {
			fmt.Println("Ca marche !")
		} else {
			UUID = tokenCookie.Value
			var compte = readuuid("user_account")
			Logged.Account = compte[0]
			Logged.Connected = true
			UserActions.Login = Logged
			fmt.Println("le compte connecté c'est le suivant\n\n", compte)
		}
	}

}
