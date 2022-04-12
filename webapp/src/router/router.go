package router

import (
	"webapp/src/router/routes"

	"github.com/gorilla/mux"
)

func GenerateRoutes() *mux.Router {
	router := mux.NewRouter()
	return routes.Configure(router)
}
