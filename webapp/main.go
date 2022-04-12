package main

import (
	"fmt"
	"log"
	"net/http"
	"webapp/src/config"
	"webapp/src/cookies"
	"webapp/src/router"
	"webapp/src/utils"
)

func init() {
	utils.LoadTemplates()
	config.LoadVars()
	cookies.ConfigureSecureCookie()
}

func main() {

	router := router.GenerateRoutes()

	fmt.Printf("Running web app on port: %d!\n\n", config.PortApp)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.PortApp), router))
}
