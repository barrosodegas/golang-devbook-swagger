package routers

import (
	"api/src/controller"
	"net/http"
)

// publicationsRouters represents the routes to publications.
var publicationsRouters = []Router{
	{
		URI:                    "/publications",
		Method:                 http.MethodPost,
		Function:               controller.CreatePublication,
		RequiresAuthentication: true,
	},
	{
		URI:                    "/publications",
		Method:                 http.MethodGet,
		Function:               controller.ListMyAndFollowPublications,
		RequiresAuthentication: true,
	},
	{
		URI:                    "/publications/{publicationId}",
		Method:                 http.MethodGet,
		Function:               controller.FindPublicationById,
		RequiresAuthentication: true,
	},
	{
		URI:                    "/publications/{publicationId}",
		Method:                 http.MethodPut,
		Function:               controller.UpdatePublicationById,
		RequiresAuthentication: true,
	},
	{
		URI:                    "/publications/{publicationId}",
		Method:                 http.MethodDelete,
		Function:               controller.DeletePublicationById,
		RequiresAuthentication: true,
	},
	{
		URI:                    "/publications/user/{userId}",
		Method:                 http.MethodGet,
		Function:               controller.ListPublicationsByUserId,
		RequiresAuthentication: true,
	},
	{
		URI:                    "/publications/{publicationId}/like",
		Method:                 http.MethodPost,
		Function:               controller.LikePublicationById,
		RequiresAuthentication: true,
	},
	{
		URI:                    "/publications/{publicationId}/unlike",
		Method:                 http.MethodPost,
		Function:               controller.UnlikePublicationById,
		RequiresAuthentication: true,
	},
}
