package controller

import (
	"api/src/authentication"
	"api/src/database"
	"api/src/model"
	"api/src/repository"
	"api/src/responses"
	"api/src/security"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// User godoc
//
// @Summary     CreateUser create a user.
// @Description CreateUser create a user and return the user created.
//
// @Tags users
//
// @Accept  json
// @Produce json
//
// @Param Authorization header string     true "Enter the content: Bearer and your access token."
// @Param user          body   model.User true " "
//
// @Success 201 {object} model.User "Success"
//
// @Failure 400 {string} string "Error: Bad Request"
// @Failure 401 {string} string "Error: Unauthorized"
// @Failure 422 {string} string "Error: Unprocessable Entity"
// @Failure 500 {string} string "Error: Internal Server"
//
// @Router /users [post]
func CreateUser(w http.ResponseWriter, r *http.Request) {
	requestBody, error := ioutil.ReadAll(r.Body)
	if error != nil {
		responses.Error(w, http.StatusUnprocessableEntity, error)
		return
	}

	var user model.User

	if error = json.Unmarshal(requestBody, &user); error != nil {
		responses.Error(w, http.StatusBadRequest, error)
		return
	}

	if error = user.Prepare("create"); error != nil {
		responses.Error(w, http.StatusBadRequest, error)
		return
	}

	db, error := database.Connect()
	if error != nil {
		responses.Error(w, http.StatusInternalServerError, error)
		return
	}

	defer db.Close()

	repository := repository.NewUsersRepository(db)

	user.ID, error = repository.CreateUser(user)
	if error != nil {
		responses.Error(w, http.StatusInternalServerError, error)
		return
	}

	responses.JSON(w, http.StatusCreated, user)
}

// User godoc
//
// @Summary     ListUsersByFilter list users by a filter.
// @Description ListUsersByFilter list users by a filter and returns the list of users found or an error.
//
// @Tags users
//
// @Accept  json
// @Produce json
//
// @Param Authorization header string true "Enter the content: Bearer and your access token."
// @Param user          query  string true "Name or nick of user."
//
// @Success 200 {array} model.User "Success"
//
// @Failure 401 {string}  string "Error: Unauthorized"
// @Failure 500 {string}  string "Error: Internal Server"
//
// @Router /users [get]
func ListUsersByFilter(w http.ResponseWriter, r *http.Request) {
	nameOrNick := strings.ToLower(r.URL.Query().Get("user"))

	db, error := database.Connect()
	if error != nil {
		responses.Error(w, http.StatusInternalServerError, error)
		return
	}

	defer db.Close()

	repository := repository.NewUsersRepository(db)

	users, error := repository.ListUsersByFilter(nameOrNick)
	if error != nil {
		responses.Error(w, http.StatusInternalServerError, error)
		return
	}

	responses.JSON(w, http.StatusOK, users)
}

// User godoc
//
// @Summary     FindUserById search for a user by the given ID.
// @Description FindUserById search for a user by the given ID and return a user or an error.
//
// @Tags users
//
// @Accept  json
// @Produce json
//
// @Param Authorization header string true "Enter the content: Bearer and your access token."
// @Param userId        path   int    true "User ID."
//
// @Success 200 {object} model.User "Success"
//
// @Failure 400 {string}  string "Error: Bad Request"
// @Failure 401 {string}  string "Error: Unauthorized"
// @Failure 404 {string}  string "Error: Not Found"
// @Failure 500 {string}  string "Error: Internal Server"
//
// @Router /users/{userId} [get]
func FindUserById(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)

	userId, error := strconv.ParseUint(parameters["userId"], 10, 64)
	if error != nil {
		responses.Error(w, http.StatusBadRequest, error)
		return
	}

	db, error := database.Connect()
	if error != nil {
		responses.Error(w, http.StatusInternalServerError, error)
		return
	}

	defer db.Close()

	repository := repository.NewUsersRepository(db)

	user, error := repository.FindUserById(userId)
	if error != nil {
		responses.Error(w, http.StatusInternalServerError, error)
		return
	}
	if user.ID == 0 {
		responses.Error(w, http.StatusNotFound, errors.New("User not found!"))
		return
	}

	responses.JSON(w, http.StatusOK, user)
}

// User godoc
//
// @Summary     UpdateUserById updates a user by the given ID.
// @Description UpdateUserById updates a user by the given ID and returns an error if unable to update the user.
//
// @Tags users
//
// @Accept  json
// @Produce json
//
// @Param Authorization header string     true "Enter the content: Bearer and your access token."
// @Param userId        path   int        true "User ID."
// @Param user          body   model.User true "User"
//
// @Success 204 "Success with no content"
//
// @Failure 400 {string}  string "Error: Bad Request"
// @Failure 401 {string}  string "Error: Unauthorized"
// @Failure 403 {string}  string "Error: Forbidden"
// @Failure 422 {string}  string "Error: Unprocessable Entity"
// @Failure 500 {string}  string "Error: Internal Server"
//
// @Router /users/{userId} [put]
func UpdateUserById(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)

	userId, error := strconv.ParseUint(parameters["userId"], 10, 64)
	if error != nil {
		responses.Error(w, http.StatusBadRequest, error)
		return
	}

	userIdOfToken, error := authentication.ExtractUserIdOfToken(r)
	if error != nil {
		responses.Error(w, http.StatusUnauthorized, error)
		return
	}
	if userIdOfToken != userId {
		responses.Error(w, http.StatusForbidden, errors.New("Invalid action!"))
		return
	}

	requestBody, error := ioutil.ReadAll(r.Body)
	if error != nil {
		responses.Error(w, http.StatusUnprocessableEntity, error)
		return
	}

	var user model.User

	if error = json.Unmarshal(requestBody, &user); error != nil {
		responses.Error(w, http.StatusBadRequest, error)
		return
	}

	if error = user.Prepare("update"); error != nil {
		responses.Error(w, http.StatusBadRequest, error)
		return
	}

	db, error := database.Connect()
	if error != nil {
		responses.Error(w, http.StatusInternalServerError, error)
		return
	}

	defer db.Close()

	repository := repository.NewUsersRepository(db)

	if error = repository.UpdateUserById(userId, user); error != nil {
		responses.Error(w, http.StatusInternalServerError, error)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

// User godoc
//
// @Summary     DeleteUserById delete a user by the given ID.
// @Description DeleteUserById delete a user by the given ID and returns an error if unable to delete the user.
//
// @Tags users
//
// @Accept  json
// @Produce json
//
// @Param Authorization header string true "Enter the content: Bearer and your access token."
// @Param userId        path   int    true "User ID."
//
// @Success 204 "Success with no content"
//
// @Failure 400 {string} string "Error: Bad Request"
// @Failure 401 {string} string "Error: Unauthorized"
// @Failure 403 {string} string "Error: Forbidden"
// @Failure 500 {string} string "Error: Internal Server"
//
// @Router /users/{userId} [delete]
func DeleteUserById(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)

	userId, error := strconv.ParseUint(parameters["userId"], 10, 64)
	if error != nil {
		responses.Error(w, http.StatusBadRequest, error)
		return
	}

	userIdOfToken, error := authentication.ExtractUserIdOfToken(r)
	if error != nil {
		responses.Error(w, http.StatusUnauthorized, error)
		return
	}
	if userIdOfToken != userId {
		responses.Error(w, http.StatusForbidden, errors.New("Invalid action!"))
		return
	}

	db, error := database.Connect()
	if error != nil {
		responses.Error(w, http.StatusInternalServerError, error)
		return
	}

	defer db.Close()

	repository := repository.NewUsersRepository(db)

	if error = repository.DeleteUserById(userId); error != nil {
		responses.Error(w, http.StatusInternalServerError, error)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

// User godoc
//
// @Summary     FollowUserById follows a user by the given ID.
// @Description FollowUserById follows a user by the given ID and returns an error if unable to follow the user.
//
// @Tags users
//
// @Accept  json
// @Produce json
//
// @Param Authorization header string true "Enter the content: Bearer and your access token."
// @Param userId        path   int    true "User ID."
//
// @Success 204 "Success with no content"
//
// @Failure 400 {string} string "Error: Bad Request"
// @Failure 401 {string} string "Error: Unauthorized"
// @Failure 403 {string} string "Error: Forbidden"
// @Failure 500 {string} string "Error: Internal Server"
//
// @Router /users/{userId}/follow [post]
func FollowUserById(w http.ResponseWriter, r *http.Request) {

	loggedUserIdOfToken, error := authentication.ExtractUserIdOfToken(r)
	if error != nil {
		responses.Error(w, http.StatusUnauthorized, error)
		return
	}

	parameters := mux.Vars(r)
	followedUserId, error := strconv.ParseUint(parameters["userId"], 10, 64)
	if error != nil {
		responses.Error(w, http.StatusBadRequest, error)
		return
	}
	if loggedUserIdOfToken == followedUserId {
		responses.Error(w, http.StatusForbidden, errors.New("Invalid action!"))
		return
	}

	db, error := database.Connect()
	if error != nil {
		responses.Error(w, http.StatusInternalServerError, error)
		return
	}

	defer db.Close()

	repository := repository.NewUsersRepository(db)

	if error = repository.FollowUserById(followedUserId, loggedUserIdOfToken); error != nil {
		responses.Error(w, http.StatusInternalServerError, error)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

// User godoc
//
// @Summary     UnfollowUserById unfollow a user by the given ID.
// @Description UnfollowUserById unfollow a user by the given ID and returns an error if unable to unfollow the user.
//
// @Tags users
//
// @Accept  json
// @Produce json
//
// @Param Authorization header string true "Enter the content: Bearer and your access token."
// @Param userId        path   int    true "User ID."
//
// @Success 204 "Success with no content"
//
// @Failure 400 {string} string "Error: Bad Request"
// @Failure 401 {string} string "Error: Unauthorized"
// @Failure 403 {string} string "Error: Forbidden"
// @Failure 500 {string} string "Error: Internal Server"
//
// @Router /users/{userId}/unfollow [post]
func UnfollowUserById(w http.ResponseWriter, r *http.Request) {

	loggedUserIdOfToken, error := authentication.ExtractUserIdOfToken(r)
	if error != nil {
		responses.Error(w, http.StatusUnauthorized, error)
		return
	}

	parameters := mux.Vars(r)
	followedUserId, error := strconv.ParseUint(parameters["userId"], 10, 64)
	if error != nil {
		responses.Error(w, http.StatusBadRequest, error)
		return
	}
	if loggedUserIdOfToken == followedUserId {
		responses.Error(w, http.StatusForbidden, errors.New("Invalid action!"))
		return
	}

	db, error := database.Connect()
	if error != nil {
		responses.Error(w, http.StatusInternalServerError, error)
		return
	}

	defer db.Close()

	repository := repository.NewUsersRepository(db)

	if error = repository.UnfollowUserById(followedUserId, loggedUserIdOfToken); error != nil {
		responses.Error(w, http.StatusInternalServerError, error)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

// User godoc
//
// @Summary     ListFollowersByFollowedUserId lists a user's followers.
// @Description ListFollowersByFollowedUserId lists a user's followers and returns a list of followers or an error.
//
// @Tags users
//
// @Accept  json
// @Produce json
//
// @Param Authorization header string true "Enter the content: Bearer and your access token."
// @Param userId        path   int    true "User ID."
//
// @Success 200 {array} model.User "Success"
//
// @Failure 400 {string} string "Error: Bad Request"
// @Failure 401 {string} string "Error: Unauthorized"
// @Failure 500 {string} string "Error: Internal Server"
//
// @Router /users/{userId}/followers [get]
func ListFollowersByFollowedUserId(w http.ResponseWriter, r *http.Request) {

	parameters := mux.Vars(r)
	followedUserId, error := strconv.ParseUint(parameters["userId"], 10, 64)
	if error != nil {
		responses.Error(w, http.StatusBadRequest, error)
		return
	}

	db, error := database.Connect()
	if error != nil {
		responses.Error(w, http.StatusInternalServerError, error)
		return
	}

	defer db.Close()

	repository := repository.NewUsersRepository(db)

	followers, error := repository.ListFollowersByFollowedUserId(followedUserId)
	if error != nil {
		responses.Error(w, http.StatusInternalServerError, error)
		return
	}

	responses.JSON(w, http.StatusOK, followers)
}

// User godoc
//
// @Summary     ListFollowedByFollowerId lists people the user follows.
// @Description ListFollowedByFollowerId lists people the user follows and returns a list of people the user follows or an error.
//
// @Tags users
//
// @Accept  json
// @Produce json
//
// @Param Authorization header string true "Enter the content: Bearer and your access token."
// @Param userId        path   int    true "User ID."
//
// @Success 200 {array} model.User "Success"
//
// @Failure 400 {string} string "Error: Bad Request"
// @Failure 401 {string} string "Error: Unauthorized"
// @Failure 500 {string} string "Error: Internal Server"
//
// @Router /users/{userId}/list-followed [get]
func ListFollowedByFollowerId(w http.ResponseWriter, r *http.Request) {

	parameters := mux.Vars(r)
	followerUserId, error := strconv.ParseUint(parameters["userId"], 10, 64)
	if error != nil {
		responses.Error(w, http.StatusBadRequest, error)
		return
	}

	db, error := database.Connect()
	if error != nil {
		responses.Error(w, http.StatusInternalServerError, error)
		return
	}

	defer db.Close()

	repository := repository.NewUsersRepository(db)

	followedList, error := repository.ListFollowedByFollowerId(followerUserId)
	if error != nil {
		responses.Error(w, http.StatusInternalServerError, error)
		return
	}

	responses.JSON(w, http.StatusOK, followedList)
}

// User godoc
//
// @Summary     UpdatePasswordByUserId updates the user's password by the given ID.
// @Description UpdatePasswordByUserId updates the user's password by the given ID and returns an error if unable to update the password.
//
// @Tags users
//
// @Accept  json
// @Produce json
//
// @Param Authorization header string         true "Enter the content: Bearer and your access token."
// @Param userId        path   int            true "User ID."
// @Param user          body   model.Password true "The current and new password of user."
//
// @Success 204 "Success with no content"
//
// @Failure 400 {string} string "Error: Bad Request"
// @Failure 401 {string} string "Error: Unauthorized"
// @Failure 403 {string} string "Error: Forbidden"
// @Failure 422 {string} string "Error: Unprocessable Entity"
// @Failure 500 {string} string "Error: Internal Server"
//
// @Router /users/{userId}/update-password [post]
func UpdatePasswordByUserId(w http.ResponseWriter, r *http.Request) {

	loggedUserIdOfToken, error := authentication.ExtractUserIdOfToken(r)
	if error != nil {
		responses.Error(w, http.StatusUnauthorized, error)
		return
	}

	parameters := mux.Vars(r)
	userId, error := strconv.ParseUint(parameters["userId"], 10, 64)
	if error != nil {
		responses.Error(w, http.StatusBadRequest, error)
		return
	}
	if loggedUserIdOfToken != userId {
		responses.Error(w, http.StatusForbidden, errors.New("Invalid action!"))
		return
	}

	requestBody, error := ioutil.ReadAll(r.Body)
	if error != nil {
		responses.Error(w, http.StatusUnprocessableEntity, error)
		return
	}

	var password model.Password

	if error = json.Unmarshal(requestBody, &password); error != nil {
		responses.Error(w, http.StatusBadRequest, error)
		return
	}

	if password.Current == password.New {
		responses.Error(w, http.StatusBadRequest, errors.New("Password already used!"))
		return
	}

	db, error := database.Connect()
	if error != nil {
		responses.Error(w, http.StatusInternalServerError, error)
		return
	}

	defer db.Close()

	repository := repository.NewUsersRepository(db)

	passwordHash, error := repository.FindPasswordByUserId(userId)
	if error != nil {
		responses.Error(w, http.StatusInternalServerError, error)
		return
	}

	if error = security.VerifyPassword(passwordHash, password.Current); error != nil {
		responses.Error(w, http.StatusUnauthorized, errors.New("Invalid Password!"))
		return
	}

	newPasswordHash, error := security.GenerateHash(password.New)
	if error != nil {
		responses.Error(w, http.StatusBadRequest, error)
		return
	}

	if error = repository.UpdatePasswordByUserId(userId, string(newPasswordHash)); error != nil {
		responses.Error(w, http.StatusInternalServerError, error)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}
