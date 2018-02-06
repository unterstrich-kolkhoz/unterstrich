package main

import (
	"flag"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/hellerve/unterstrich/artworks"
	"github.com/hellerve/unterstrich/config"
	"github.com/hellerve/unterstrich/db"
	"github.com/hellerve/unterstrich/static"
	"github.com/hellerve/unterstrich/users"
)

func main() {
	configfile := flag.String("config", "./etc/_/server.conf", "Configuration file location")
	flag.Parse()
	conf, err := config.ReadConfig(*configfile)

	if err != nil {
		log.Fatal("Loading configuration failed: ", err)
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	dbconn, err := db.Create(conf.SQLDialect, conf.SQLName)
	defer func() {
		log.Fatal(dbconn.Close())
	}()

	if err != nil {
		log.Fatal("Connecting to database failed: ", err)
	}

	authfun := users.InitializeAuth(dbconn, router)
	users.Initialize(dbconn, router, authfun)
	artworks.Initialize(dbconn, router, authfun)
	static.Initialize(conf.Staticdir, router)
	log.Fatal(router.Run(conf.Port))
}
