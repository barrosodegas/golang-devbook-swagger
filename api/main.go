package main

import (
	"api/src/config"
	"api/src/router"
	"fmt"
	"log"
	"net/http"

	_ "api/docs"
)

// @title DevBook API
// @version 1.0
// @description API responsible for CRUD and DevBook social network authentication.
// @host localhost:5000
// @BasePath /
func main() {
	config.Load()

	fmt.Printf("Escutando a porta %d...\n\n", config.Port)
	r := router.Generate()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r))
}
