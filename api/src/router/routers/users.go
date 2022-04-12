package routers

import (
	"api/src/controller"
	"net/http"
)

// usersRouters represents the routes for users.
var usersRouters = []Router{
	{
		URI:                    "/users",
		Method:                 http.MethodPost,
		Function:               controller.CreateUser,
		RequiresAuthentication: false,
	},
	{
		URI:                    "/users",
		Method:                 http.MethodGet,
		Function:               controller.ListUsersByFilter,
		RequiresAuthentication: true,
	},
	{
		URI:                    "/users/{userId}",
		Method:                 http.MethodGet,
		Function:               controller.FindUserById,
		RequiresAuthentication: true,
	},
	{
		URI:                    "/users/{userId}",
		Method:                 http.MethodPut,
		Function:               controller.UpdateUserById,
		RequiresAuthentication: true,
	},
	{
		URI:                    "/users/{userId}",
		Method:                 http.MethodDelete,
		Function:               controller.DeleteUserById,
		RequiresAuthentication: true,
	},
	{
		URI:                    "/users/{userId}/follow",
		Method:                 http.MethodPost,
		Function:               controller.FollowUserById,
		RequiresAuthentication: true,
	},
	{
		URI:                    "/users/{userId}/unfollow",
		Method:                 http.MethodPost,
		Function:               controller.UnfollowUserById,
		RequiresAuthentication: true,
	},
	{
		URI:                    "/users/{userId}/followers",
		Method:                 http.MethodGet,
		Function:               controller.ListFollowersByFollowedUserId,
		RequiresAuthentication: true,
	},
	{
		URI:                    "/users/{userId}/list-followed",
		Method:                 http.MethodGet,
		Function:               controller.ListFollowedByFollowerId,
		RequiresAuthentication: true,
	},
	{
		URI:                    "/users/{userId}/update-password",
		Method:                 http.MethodPost,
		Function:               controller.UpdatePasswordByUserId,
		RequiresAuthentication: true,
	},
}
