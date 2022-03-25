package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func HandleReq() {
	log.Println("Start development server localhost:9999")

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", HomePage)
	myRouter.HandleFunc("/user", CreateUser).Methods("OPTIONS", "POST")
	myRouter.HandleFunc("/users", ListUsers).Methods("OPTIONS", "GET")
	myRouter.HandleFunc("/user/{id}", DetailUser).Methods("OPTIONS", "GET")
	// myRouter.HandleFunc("/article/{id}", UpdateArticle).Methods("OPTIONS", "PUT")
	// myRouter.HandleFunc("/article/{id}", DeleteArticle).Methods("OPTIONS", "Delete")

	log.Fatal(http.ListenAndServe(":5002", myRouter))
}
