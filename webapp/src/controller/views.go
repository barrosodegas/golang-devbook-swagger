package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"webapp/src/config"
	"webapp/src/cookies"
	"webapp/src/model"
	"webapp/src/request"
	"webapp/src/responses"
	"webapp/src/service"
	"webapp/src/utils"

	"github.com/gorilla/mux"
)

func LoadLoginView(w http.ResponseWriter, r *http.Request) {

	cookies, _ := cookies.Read(r)
	if cookies["token"] != "" {
		http.Redirect(w, r, "/home", http.StatusFound)
		return
	}

	utils.ExecuteTemplate(w, "login.html", nil)
}

func LoadCreateUserView(w http.ResponseWriter, _ *http.Request) {
	utils.ExecuteTemplate(w, "create-user.html", nil)
}

func LoadHome(w http.ResponseWriter, r *http.Request) {
	publicationsUrl := fmt.Sprintf("%s/publications", config.UrlApi)
	response, error := request.ExecuteRequestWithAuthentication(r, http.MethodGet, publicationsUrl, nil)
	if error != nil {
		responses.JSON(w, http.StatusInternalServerError, responses.ResponseError{Error: error.Error()})
		return
	}

	defer response.Body.Close()

	if response.StatusCode >= 400 {
		responses.PrepareAndReturnError(w, response)
		return
	}

	var publications []model.Publication

	if error = json.NewDecoder(response.Body).Decode(&publications); error != nil {
		responses.JSON(w, http.StatusUnprocessableEntity, responses.ResponseError{Error: error.Error()})
		return
	}

	cookies, _ := cookies.Read(r)
	userId, _ := strconv.ParseUint(cookies["id"], 10, 64)

	utils.ExecuteTemplate(w, "home.html", struct {
		Publications []model.Publication
		UserId       uint64
	}{
		Publications: publications,
		UserId:       userId,
	})
}

func LoadEditPublicationView(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	publicationId, error := strconv.ParseUint(parameters["publicationId"], 10, 64)
	if error != nil {
		responses.JSON(w, http.StatusBadRequest, responses.ResponseError{Error: error.Error()})
		return
	}

	publicationUrl := fmt.Sprintf("%s/publications/%d", config.UrlApi, publicationId)
	response, error := request.ExecuteRequestWithAuthentication(r, http.MethodGet, publicationUrl, nil)
	if error != nil {
		responses.JSON(w, http.StatusInternalServerError, responses.ResponseError{Error: error.Error()})
		return
	}

	defer response.Body.Close()

	if response.StatusCode >= 400 {
		responses.PrepareAndReturnError(w, response)
		return
	}

	var publication model.Publication

	if error = json.NewDecoder(response.Body).Decode(&publication); error != nil {
		responses.JSON(w, http.StatusUnprocessableEntity, responses.ResponseError{Error: error.Error()})
		return
	}

	utils.ExecuteTemplate(w, "edit-publication.html", publication)
}

func LoadUsersView(w http.ResponseWriter, r *http.Request) {

	nameOrNick := strings.ToLower(r.URL.Query().Get("user"))

	listUsersUrl := fmt.Sprintf("%s/users?user=%s", config.UrlApi, nameOrNick)
	response, error := request.ExecuteRequestWithAuthentication(r, http.MethodGet, listUsersUrl, nil)
	if error != nil {
		responses.JSON(w, http.StatusInternalServerError, responses.ResponseError{Error: error.Error()})
		return
	}

	defer response.Body.Close()

	if response.StatusCode >= 400 {
		responses.PrepareAndReturnError(w, response)
		return
	}

	var users []model.User

	if error = json.NewDecoder(response.Body).Decode(&users); error != nil {
		responses.JSON(w, http.StatusUnprocessableEntity, responses.ResponseError{Error: error.Error()})
		return
	}

	utils.ExecuteTemplate(w, "users.html", users)
}

func LoadUserProfileView(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	userId, error := strconv.ParseUint(parameters["userId"], 10, 64)
	if error != nil {
		responses.JSON(w, http.StatusBadRequest, responses.ResponseError{Error: error.Error()})
		return
	}

	userUrl := fmt.Sprintf("%s/users/%d", config.UrlApi, userId)
	response, error := request.ExecuteRequestWithAuthentication(r, http.MethodGet, userUrl, nil)
	if error != nil {
		responses.JSON(w, http.StatusInternalServerError, responses.ResponseError{Error: error.Error()})
		return
	}

	defer response.Body.Close()

	if response.StatusCode >= 400 {
		responses.PrepareAndReturnError(w, response)
		return
	}

	cookies, _ := cookies.Read(r)
	loggedUserId, _ := strconv.ParseUint(cookies["id"], 10, 64)

	if userId == loggedUserId {
		http.Redirect(w, r, "/me", http.StatusFound)
		return
	}

	var user model.User

	if error = json.NewDecoder(response.Body).Decode(&user); error != nil {
		responses.JSON(w, http.StatusUnprocessableEntity, responses.ResponseError{Error: error.Error()})
		return
	}

	if error = service.LoadUserDataByUserId(&user, r); error != nil {
		responses.JSON(w, http.StatusInternalServerError, responses.ResponseError{Error: error.Error()})
		return
	}

	utils.ExecuteTemplate(w, "user-profile.html", struct {
		FollowedByLoggedUser bool
		User                 model.User
	}{
		FollowedByLoggedUser: followedByLoggedUser(user, loggedUserId),
		User:                 user,
	})
}

func followedByLoggedUser(user model.User, loggedUserId uint64) bool {

	followedByLoggedUser := false
	for _, follower := range user.Followers {
		if follower.ID == loggedUserId {
			followedByLoggedUser = true
			break
		}
	}
	return followedByLoggedUser
}

func LoadLoggedUserProfileView(w http.ResponseWriter, r *http.Request) {
	cookies, _ := cookies.Read(r)
	loggedUserId, _ := strconv.ParseUint(cookies["id"], 10, 64)

	userUrl := fmt.Sprintf("%s/users/%d", config.UrlApi, loggedUserId)
	response, error := request.ExecuteRequestWithAuthentication(r, http.MethodGet, userUrl, nil)
	if error != nil {
		responses.JSON(w, http.StatusInternalServerError, responses.ResponseError{Error: error.Error()})
		return
	}

	defer response.Body.Close()

	if response.StatusCode >= 400 {
		responses.PrepareAndReturnError(w, response)
		return
	}

	var user model.User

	if error = json.NewDecoder(response.Body).Decode(&user); error != nil {
		responses.JSON(w, http.StatusUnprocessableEntity, responses.ResponseError{Error: error.Error()})
		return
	}

	if error = service.LoadUserDataByUserId(&user, r); error != nil {
		responses.JSON(w, http.StatusInternalServerError, responses.ResponseError{Error: error.Error()})
		return
	}

	utils.ExecuteTemplate(w, "me.html", struct{ User model.User }{User: user})
}

func LoadUpdateUserView(w http.ResponseWriter, r *http.Request) {
	cookies, _ := cookies.Read(r)
	loggedUserId, _ := strconv.ParseUint(cookies["id"], 10, 64)

	userUrl := fmt.Sprintf("%s/users/%d", config.UrlApi, loggedUserId)
	response, error := request.ExecuteRequestWithAuthentication(r, http.MethodGet, userUrl, nil)
	if error != nil {
		responses.JSON(w, http.StatusInternalServerError, responses.ResponseError{Error: error.Error()})
		return
	}

	defer response.Body.Close()

	if response.StatusCode >= 400 {
		responses.PrepareAndReturnError(w, response)
		return
	}

	var user model.User

	if error = json.NewDecoder(response.Body).Decode(&user); error != nil {
		responses.JSON(w, http.StatusUnprocessableEntity, responses.ResponseError{Error: error.Error()})
		return
	}

	utils.ExecuteTemplate(w, "update-user.html", user)
}

func LoadUpdatePasswordView(w http.ResponseWriter, r *http.Request) {
	cookies, _ := cookies.Read(r)
	loggedUserId, _ := strconv.ParseUint(cookies["id"], 10, 64)

	userUrl := fmt.Sprintf("%s/users/%d", config.UrlApi, loggedUserId)
	response, error := request.ExecuteRequestWithAuthentication(r, http.MethodGet, userUrl, nil)
	if error != nil {
		responses.JSON(w, http.StatusInternalServerError, responses.ResponseError{Error: error.Error()})
		return
	}

	defer response.Body.Close()

	if response.StatusCode >= 400 {
		responses.PrepareAndReturnError(w, response)
		return
	}

	var user model.User

	if error = json.NewDecoder(response.Body).Decode(&user); error != nil {
		responses.JSON(w, http.StatusUnprocessableEntity, responses.ResponseError{Error: error.Error()})
		return
	}

	utils.ExecuteTemplate(w, "update-password.html", user)
}
