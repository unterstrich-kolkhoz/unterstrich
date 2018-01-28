package main

import (
	"flag"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/hellerve/artifex/artworks"
	"github.com/hellerve/artifex/config"
	"github.com/hellerve/artifex/db"
	"github.com/hellerve/artifex/users"
)

func main() {
	configfile := flag.String("config", "./etc/arfx/server.conf", "Configuration file location")
	flag.Parse()
	conf, err := config.ReadConfig(*configfile)

	if err != nil {
		log.Fatal("Loading configuration failed: ", err)
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	dbconn, err := db.Create(conf.SQLDialect, conf.SQLName)
	defer dbconn.Close()

	if err != nil {
		log.Fatal("Connecting to database failed: ", err)
	}

	authfun := users.InitializeAuth(dbconn, router)
	users.Initialize(dbconn, router, authfun)
	artworks.Initialize(dbconn, router, authfun)
	router.Run()
}
