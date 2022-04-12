package routers

import (
	"api/src/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

// Router the representation of a route.
type Router struct {
	URI                    string
	Method                 string
	Function               func(w http.ResponseWriter, r *http.Request)
	RequiresAuthentication bool
}

// Configure responsible for configuring API routes and deciding whether access should be authenticated or not.
func Configure(r *mux.Router) *mux.Router {
	routers := getRoutes()

	for _, route := range routers {
		if route.RequiresAuthentication {
			r.HandleFunc(
				route.URI,
				middleware.Logger(middleware.Authenticate(route.Function)),
			).Methods(route.Method)
		} else {
			r.HandleFunc(route.URI, middleware.Logger(route.Function)).Methods(route.Method)
		}
	}
	return r
}

// getRoutes retrieve all routes that should be exposed by the API.
func getRoutes() []Router {
	routers := usersRouters
	routers = append(routers, loginRoute)
	routers = append(routers, publicationsRouters...)
	return routers
}
