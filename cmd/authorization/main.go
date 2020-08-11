package main

import (
	"fmt"

	"github.com/arturoguerra/go-logging"
	"github.com/destinyarena/authorization/internal/endpoints/oauth2"
	"github.com/destinyarena/authorization/internal/router"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	log := logging.New()

	rconfig, err := router.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	r := router.New(log, rconfig)

	oauth2Group := r.Group("/api/oauth", middleware.Logger())

	oauth2Config, err := oauth2.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	oauth2Handler, err := oauth2.New(log, oauth2Config)
	if err != nil {
		log.Fatal(err)
	}

	oauth2Handler.Register(oauth2Group)

	log.Infof("Running on %s:%s", rconfig.Host, rconfig.Port)
	r.Logger.Fatal(r.Start(fmt.Sprintf("%s:%s", rconfig.Host, rconfig.Port)))
}
