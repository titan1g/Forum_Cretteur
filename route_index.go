package main

import (
	"Forum/data"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)



func index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	
	threads, err := data.Threads()
	if err != nil {
		error_message(w,r,ps,"Cannot get threads")
	} else {
		s, err := session(w,r)
		loginUser , _ := s.User()
	
		if err != nil {
			fetch := data.DisplayData{
				AllThreads: threads,
			}
			generateHTML(w, &fetch, "layout", "public.navbar", "index")
		} else {
			
			if err != nil {
				return
			}
			fetch := data.DisplayData{
				User: loginUser,
				AllThreads: threads,
			}
			fmt.Println(fetch.User.UserName)
			generateHTML(w, &fetch, "layout", "private.navbar", "index")
		}
	}
}

func err(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	vals := r.URL.Query()
	generateHTML(w, vals.Get("msg"), "layout", "public.navbar", "error")

}
