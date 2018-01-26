package main

import (
  "flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/hellerve/artifex/config"
  "github.com/hellerve/artifex/users"
)

func main() {
	configfile := flag.String("config", "./etc/arfx/server.conf", "Configuration file location")
	flag.Parse()
	router := mux.NewRouter()
	conf, err := config.ReadConfig(*configfile)
  log.Println(conf.Port)

	if err != nil {
		log.Fatal("Loading configuration failed: ", err)
	}

    router.HandleFunc("/users", users.GetUsers).Methods("GET")
    router.HandleFunc("/users", users.CreateUser).Methods("POST")
    router.HandleFunc("/users/{id}", users.GetUser).Methods("GET")
    router.HandleFunc("/users/{id}", users.UpdateUser).Methods("PATCH")
    router.HandleFunc("/users/{id}", users.DeleteUser).Methods("DELETE")
	log.Fatal(http.ListenAndServe(conf.Port, router))
}
