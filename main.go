package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/hellerve/artifex/config"
	"github.com/hellerve/artifex/db"
	"github.com/hellerve/artifex/users"
)

func main() {
	configfile := flag.String("config", "./etc/arfx/server.conf", "Configuration file location")
	flag.Parse()
	router := mux.NewRouter()
	conf, err := config.ReadConfig(*configfile)

	if err != nil {
		log.Fatal("Loading configuration failed: ", err)
	}

	dbconn, err := db.Create(conf.SQLDialect, conf.SQLName)
	defer dbconn.Close()

	if err != nil {
		log.Fatal("Connecting to database failed: ", err)
	}

	users.Initialize(dbconn, router)
	log.Fatal(http.ListenAndServe(conf.Port, router))
}
