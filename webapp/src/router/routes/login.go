package routes

import (
	"net/http"
	"webapp/src/controller"
)

var loginRoutes = []Route{
	{
		URI:                   "/",
		Method:                http.MethodGet,
		Function:              controller.LoadLoginView,
		RequireAuthentication: false,
	},
	{
		URI:                   "/login",
		Method:                http.MethodGet,
		Function:              controller.LoadLoginView,
		RequireAuthentication: false,
	},
	{
		URI:                   "/login",
		Method:                http.MethodPost,
		Function:              controller.Login,
		RequireAuthentication: false,
	},
}
