package main

import (
	"flag"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/unterstrich-kolkhoz/unterstrich/artworks"
	"github.com/unterstrich-kolkhoz/unterstrich/config"
	"github.com/unterstrich-kolkhoz/unterstrich/db"
	"github.com/unterstrich-kolkhoz/unterstrich/endpoints"
	"github.com/unterstrich-kolkhoz/unterstrich/static"
	"github.com/unterstrich-kolkhoz/unterstrich/subsite"
	"github.com/unterstrich-kolkhoz/unterstrich/users"
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

	context := endpoints.Context{dbconn, conf}

	authfun := users.InitializeAuth(dbconn, router)
	artworks.Initialize(&context, router, authfun)
	users.Initialize(&context, router, authfun)
	subsite.Initialize(&context, router, authfun)
	static.Initialize(conf.Staticdir, router)
	log.Fatal(router.Run(conf.Port))
}
