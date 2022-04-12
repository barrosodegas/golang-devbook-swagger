package routers

import (
	"api/src/controller"
	"net/http"
)

// loginRoute represents a login route.
var loginRoute = Router{
	URI:                    "/login",
	Method:                 http.MethodPost,
	Function:               controller.Login,
	RequiresAuthentication: false,
}
