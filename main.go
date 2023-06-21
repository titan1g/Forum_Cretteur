package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

func main() {
	
	router := httprouter.New()
	router.ServeFiles("/static/*filepath",http.Dir("static"))

	
	// index
	router.GET("/", index)

	// error
	router.GET("/err", err)

	// route_auth.go
	router.GET("/login", login)
	router.GET("/signup", signup)
	router.POST("/signup_account", signupAccount)
	router.POST("/authenticate", authenticate)
	router.GET("/logout", logout)

	// route_thread
	router.GET("/thread/new", newThread)
	router.POST("/thread/create", createThread)
	router.GET("/thread/read", readThread)
	router.POST("/thread/post", postThread)

	router.GET("/t/:params", readThreadBySlug)


	router.GET("/hello/:name", Hello)
	server := &http.Server{
		Addr:           config.Address,
		Handler:        router,
		ReadTimeout:    time.Duration(config.ReadTimeout * int64(time.Second)),
		WriteTimeout:   time.Duration(config.WriteTimeout * int64(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(server.ListenAndServe())
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

