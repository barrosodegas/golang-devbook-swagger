package routes

import (
	"net/http"
	"webapp/src/controller"
)

var logoutRoute = Route{
	URI:                   "/logout",
	Method:                http.MethodGet,
	Function:              controller.Logout,
	RequireAuthentication: true,
}
