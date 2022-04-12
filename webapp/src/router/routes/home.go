package routes

import (
	"net/http"
	"webapp/src/controller"
)

var homeRoute = Route{
	URI:                   "/home",
	Method:                http.MethodGet,
	Function:              controller.LoadHome,
	RequireAuthentication: true,
}
