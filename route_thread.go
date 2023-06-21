package main

import (
	"Forum/data"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/stoewer/go-strcase"
)

// GET /threads/new
// Show the new thread form page
func newThread(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	_, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	} else {

		generateHTML(w, nil, "layout", "private.navbar", "new.thread")
	}
}

// POST /thread/create

func createThread(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	} else {
		err = r.ParseForm()
		if err != nil {
			danger(err, "Cannot parse form")
		}
		user, err := sess.User()
		if err != nil {
			danger(err, "Cannot get user from session")
		}
		topic := r.PostFormValue("topic")
		uuid := data.CreateUUID()
		discussion := r.PostFormValue("discussion")
		slug := strcase.KebabCase(topic) + "." + uuid[0:4]
		if _, err := user.CreateThread(uuid, topic, slug, discussion); err != nil {
			danger(err, "Cannot create thread")
		}
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// GET /thread/read
// Show the details of the thread, including the posts and the form to write a post
func readThread(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	vals := r.URL.Query()
	uuid := vals.Get("id")

	thread, err := data.ThreadByUUID(uuid)
	if err != nil {
		error_message(w, r, ps, "Cannot read thread")
	} else {
		s, err := session(w, r)
		loginUser, _ := s.User()
		if err != nil {
			generateHTML(w, &thread, "layout", "public.navbar", "public.thread")
		} else {
			fetch_thread := data.DisplayData{
				User:         loginUser,
				SingleThread: thread,
			}
			generateHTML(w, &fetch_thread, "layout", "private.navbar", "private.thread")

		}
	}
}

// GET /t/params
// get the thread
func readThreadBySlug(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	slug := ps.ByName("params")
	thread, err := data.ThreadBySlug(slug)
	if err != nil {
		error_message(w, r, ps, "Cannot get the thread")
	} else {
		s, err := session(w, r)
		loginUser, _ := s.User()
		if err != nil {
			generateHTML(w, &thread, "layout", "public.navbar", "public.thread")
		} else {
			fetchThread := data.DisplayData{
				User:         loginUser,
				SingleThread: thread,
			}
			generateHTML(w, &fetchThread, "layout", "private.navbar", "private.thread")

		}
	}
}

// POST /thread/post
// Create the post
func postThread(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	} else {
		err = r.ParseForm()
		if err != nil {
			danger(err, "Cannot parse form")
		}
		user, err := sess.User()
		if err != nil {
			danger(err, "Cannot get user from session")
		}
		body := r.PostFormValue("body")
		uuid := r.PostFormValue("uuid")

		thread, err := data.ThreadByUUID(uuid)
		if err != nil {
			error_message(w, r, ps, "Cannot read thread")
		}
		if _, err := user.CreatePost(thread, body); err != nil {
			danger(err, "Cannot created post")
		}
		url := fmt.Sprint("/thread/read?id=", uuid)
		http.Redirect(w, r, url, http.StatusFound)
	}
}
