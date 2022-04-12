package routes

import (
	"net/http"
	"webapp/src/middleware"

	"github.com/gorilla/mux"
)

type Route struct {
	URI                   string
	Method                string
	Function              func(w http.ResponseWriter, r *http.Request)
	RequireAuthentication bool
}

func Configure(router *mux.Router) *mux.Router {
	routes := getRoutes()
	for _, route := range routes {

		if route.RequireAuthentication {
			router.HandleFunc(
				route.URI,
				middleware.Logger(middleware.Authenticate(route.Function)),
			).Methods(route.Method)
		} else {
			router.HandleFunc(route.URI, middleware.Logger(route.Function)).Methods(route.Method)
		}
	}

	fileServer := http.FileServer(http.Dir("./assets/"))
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fileServer))

	return router
}

func getRoutes() []Route {
	routes := loginRoutes
	routes = append(routes, usersRoutes...)
	routes = append(routes, homeRoute)
	routes = append(routes, publicationsRoutes...)
	routes = append(routes, logoutRoute)
	return routes
}
