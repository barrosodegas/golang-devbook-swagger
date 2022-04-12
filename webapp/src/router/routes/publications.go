package routes

import (
	"net/http"
	"webapp/src/controller"
)

var publicationsRoutes = []Route{
	{
		URI:                   "/publications",
		Method:                http.MethodPost,
		Function:              controller.CreatePublication,
		RequireAuthentication: true,
	},
	{
		URI:                   "/publications/{publicationId}/like",
		Method:                http.MethodPost,
		Function:              controller.LikePublication,
		RequireAuthentication: true,
	},
	{
		URI:                   "/publications/{publicationId}/unlike",
		Method:                http.MethodPost,
		Function:              controller.UnlikePublication,
		RequireAuthentication: true,
	},
	{
		URI:                   "/publications/{publicationId}/edit",
		Method:                http.MethodGet,
		Function:              controller.LoadEditPublicationView,
		RequireAuthentication: true,
	},
	{
		URI:                   "/publications/{publicationId}",
		Method:                http.MethodPut,
		Function:              controller.UpdatePublication,
		RequireAuthentication: true,
	},
	{
		URI:                   "/publications/{publicationId}",
		Method:                http.MethodDelete,
		Function:              controller.DeletePublication,
		RequireAuthentication: true,
	},
}
