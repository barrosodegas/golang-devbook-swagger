package routes

import (
	"net/http"
	"webapp/src/controller"
)

var usersRoutes = []Route{
	{
		URI:                   "/create-user",
		Method:                http.MethodGet,
		Function:              controller.LoadCreateUserView,
		RequireAuthentication: false,
	},
	{
		URI:                   "/users",
		Method:                http.MethodPost,
		Function:              controller.CreateUser,
		RequireAuthentication: false,
	},
	{
		URI:                   "/users",
		Method:                http.MethodGet,
		Function:              controller.LoadUsersView,
		RequireAuthentication: true,
	},
	{
		URI:                   "/users/{userId}",
		Method:                http.MethodGet,
		Function:              controller.LoadUserProfileView,
		RequireAuthentication: true,
	},
	{
		URI:                   "/users/{userId}/unfollow",
		Method:                http.MethodPost,
		Function:              controller.UnfollowUserByUserId,
		RequireAuthentication: true,
	},
	{
		URI:                   "/users/{userId}/follow",
		Method:                http.MethodPost,
		Function:              controller.FollowUserByUserId,
		RequireAuthentication: true,
	},
	{
		URI:                   "/me",
		Method:                http.MethodGet,
		Function:              controller.LoadLoggedUserProfileView,
		RequireAuthentication: true,
	},
	{
		URI:                   "/me/edit",
		Method:                http.MethodGet,
		Function:              controller.LoadUpdateUserView,
		RequireAuthentication: true,
	},
	{
		URI:                   "/me",
		Method:                http.MethodPut,
		Function:              controller.UpdateUser,
		RequireAuthentication: true,
	},
	{
		URI:                   "/me/update-password",
		Method:                http.MethodGet,
		Function:              controller.LoadUpdatePasswordView,
		RequireAuthentication: true,
	},
	{
		URI:                   "/me/update-password",
		Method:                http.MethodPost,
		Function:              controller.UpdateUserPassword,
		RequireAuthentication: true,
	},
	{
		URI:                   "/users",
		Method:                http.MethodDelete,
		Function:              controller.DeleteUserById,
		RequireAuthentication: true,
	},
}
