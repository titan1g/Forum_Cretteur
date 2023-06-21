package main

import (
	"Forum/data"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// GET /login
// Show the login page
func login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	t := parseTemplateFiles("login.layout", "public.navbar", "login")
	t.Execute(w, nil)
}

// GET /signup
// Show the sign up page
func signup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	generateHTML(w, nil, "login.layout", "public.navbar", "signup")
}

// POST /signup
// create user account

func signupAccount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := r.ParseForm()
	if err != nil {
		danger(err, "Cannot get the data")
	}

	user := data.User{
		Name:     r.PostFormValue("name"),
		Email:    r.PostFormValue("email"),
		UserName: r.PostFormValue("username"),
		Password: r.PostFormValue("password"),
	}
	error := user.Create()
	if error != nil {
		danger(error, "Cannot creat user")
	}
	http.Redirect(w, r, "/login", http.StatusFound)

}

// Post /authenticate
// authenticate the user given the email or password
var logged_user data.User

//var user_session data.Session
func authenticate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := r.ParseForm()
	if err != nil {
		danger(err, "Cannot get the data")
	}
	user, err := data.UserByEmail(r.PostFormValue("email"))

	if err != nil {
		danger(err, "User dont exist")
	}

	if err != nil {
		danger(err, "cannot get the session")
	}
	// Check if the password match
	if user.Password == data.Encrypt(r.PostFormValue("password")) {
		// assign user to the loggeduser
		logged_user = user
		session, err := user.CreateSession()
		if err != nil {
			danger(err, "Cannot create session")
			panic(err.Error())
		}
		cookie := http.Cookie{
			Name:     "_cookie",
			Value:    session.Uuid,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// GET /logout

func logout(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cookie, err := r.Cookie("_cookie")

	if err != http.ErrNoCookie {
		session := data.Session{Uuid: cookie.Value}
		session.DeleleByUUID()
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
