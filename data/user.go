package data

import "time"

type User struct {
	Id        int
	Uuid      string
	Name      string
	Email     string
	UserName  string
	Password  string
	CreatedAt time.Time
}

type Session struct {
	Id        int
	Uuid      string
	Email     string
	Username  string
	UserId    int
	CreatedAt time.Time
}

//----------------------------------Session---------------------------------------------
// Create a new session for an existed user
func (user *User) CreateSession() (session Session, err error) {
	statement := "INSERT INTO sessions (uuid, email, user_name, user_id, created_at) " +
		"VALUES ($1, $2, $3, $4, $5) returning id,uuid,email,user_name,user_id,created_at"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	// query row to scan into sessions struct
	err = stmt.QueryRow(CreateUUID(), user.Email, user.UserName, user.Id, time.Now()).Scan(
		&session.Id, &session.Uuid, &session.Email, &session.Username, &session.UserId, &session.CreatedAt)
	return
}

// Send the session for user
func (user *User) Session() (session Session, err error) {
	session = Session{}
	err = Db.QueryRow("SELECT id,uuid,email,user_name,user_id,created_at FROM sessions WHERE user_id=$1", user.Id).
		Scan(&session.Id, &session.Uuid, &session.Email, &session.Username, &session.UserId, &session.CreatedAt)
	return
}

// Check if the session is valid in the database
func (session *Session) Check() (valid bool, err error) {
	err = Db.QueryRow("SELECT id, uuid, email,user_name, user_id, created_at FROM sessions WHERE uuid = $1", session.Uuid).
		Scan(&session.Id, &session.Uuid, &session.Email, &session.Username, &session.UserId, &session.CreatedAt)
	if err != nil {
		valid = false
		return
	}
	if session.Id != 0 {
		valid = true
	}
	return
}

// delete session from database
func (session *Session) DeleleByUUID() (err error) {
	statement := "DELETE from sessions WHERE uuid=$1"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(session.Uuid)
	return
}

// Get the user from session
func (session *Session) User() (user User, err error) {
	user = User{}
	err = Db.QueryRow("SELECT id,uuid,name,email,user_name,created_at FROM users WHERE id=$1", session.UserId).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.UserName, &user.CreatedAt)
	return
}

// Delete all sessions from database
func SessionDeleteAll() (err error) {
	statement := "delete from sessions"
	_, err = Db.Exec(statement)
	return
}

// --------------------------------------------END SESSION-------------------------------------

// Create new user
func (user *User) Create() (err error) {
	statement := "INSERT INTO users(uuid,name,email,user_name,password,created_at) " +
		" VALUES ($1, $2, $3, $4, $5, $6) returning id,uuid,created_at"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(CreateUUID(), user.Name, user.Email, user.UserName, Encrypt(user.Password), time.Now()).
		Scan(&user.Id, &user.Uuid, &user.CreatedAt)
	return
}

// Delete user from database
func (user *User) Delete() (err error) {
	statement := "delete from users where id=$1"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Id)
	return
}

// Update user information
func (user *User) Update() (err error) {
	statement := "UPDATE users set name=$2,email=$3 where id=$1"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(user.Id, user.Name, user.Email)
	return
}

// Delete all users from database
func UserDeleteAll() (err error) {
	statement := "delete from users"
	_, err = Db.Exec(statement)
	return
}

// Get all users in the database and returns it
func Users() (users []User, err error) {
	rows, err := Db.Query("SELECT id, uuid, name, email, user_name, password, created_at FROM users")
	if err != nil {
		return
	}
	for rows.Next() {
		user := User{}
		if err = rows.Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.UserName, &user.Password, &user.CreatedAt); err != nil {
			return
		}
		users = append(users, user)
	}
	rows.Close()
	return
}

// Get a single user given the email
func UserByEmail(email string) (user User, err error) {
	user = User{}
	err = Db.QueryRow("SELECT id, uuid, name, email,user_name ,password, created_at FROM users WHERE email = $1", email).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.UserName, &user.Password, &user.CreatedAt)
	return
}

// Get a single user given the UUID
func UserByUUID(uuid string) (user User, err error) {
	user = User{}
	err = Db.QueryRow("SELECT id, uuid, name, email, password, created_at FROM users WHERE uuid = $1", uuid).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.UserName, &user.Password, &user.CreatedAt)
	return
}
