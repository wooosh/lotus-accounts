package httpserver

import (
	"fmt"
	"log"
	"net/http"

	"lotusaccounts/backend"

	"github.com/julienschmidt/httprouter"
)

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

	log.Println("Attempt to create new token: username='" + username + "', password=REDACTED, location='" + location + "'")

	token, err := backend.CreateToken(username, password, location)
	if err != nil {
		log.Println("Token creation failed: ", err)
		// TODO: bad request vs forbidden
		w.WriteHeader(http.StatusForbidden)
		// TODO: show error only if it is an expected error
		fmt.Fprintln(w, "Forbidden:", err)
	} else {
		log.Println("Token creation succesful")
		fmt.Fprint(w, token)
	}
}

func Start() {
	// TODO: error recovery middleware
	log.Print("Starting HTTP server")
	router := httprouter.New()
	router.GET("/", index)
	router.POST("/api/new_token", newToken)

	// TODO: serve on unix domain socket
	// TODO: csrf??
	// TODO: reject not https or non localhost connections
	log.Fatal(http.ListenAndServe(":7655", router))
}
