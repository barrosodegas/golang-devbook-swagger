package router

import (
	"api/src/router/routers"

	"github.com/gorilla/mux"

	httpSwagger "github.com/swaggo/http-swagger"
)

// Generate responsible for generating and configuring API routes.
func Generate() *mux.Router {
	r := mux.NewRouter()
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	return routers.Configure(r)
}
