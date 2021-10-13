package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Print("request")
	fmt.Fprint(w, "Hello\n")
}

func newToken(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := r.ParseMultipartForm(4096)
	if err != nil {
		log.Fatal(err)
	}

	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	location := r.PostFormValue("location")
	if location == "" {
		location = r.RemoteAddr + " " + r.UserAgent()
	}

	log.Print("Attempt to create new token: username='" + username + "', password=REDACTED, location='" + location + "'")

	token, err := createNewToken(username, password, location)
	if err != nil {
		log.Print("Token creation failed: ", err)
		// TODO: bad request vs forbidden
		w.WriteHeader(http.StatusForbidden)
		// TODO: show error only if it is an expected error
		fmt.Fprintln(w, "Forbidden:", err)
	} else {
		fmt.Fprint(w, token)
	}
}

func httpServer() {
	// TODO: error recovery middleware
	log.Print("Starting HTTP server")
	router := httprouter.New()
	router.GET("/", Index)
	router.POST("/api/new_token", newToken)

	// TODO: serve on unix domain socket
	// TODO: csrf??
	// TODO: reject not https or non localhost connections
	log.Fatal(http.ListenAndServe(":7655", router))
}
