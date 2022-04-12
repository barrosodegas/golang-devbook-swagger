package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"webapp/src/config"
	"webapp/src/cookies"
	"webapp/src/request"
	"webapp/src/responses"

	"github.com/gorilla/mux"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	user, error := json.Marshal(map[string]string{
		"name":     r.FormValue("name"),
		"nick":     r.FormValue("nick"),
		"email":    r.FormValue("email"),
		"password": r.FormValue("password"),
	})
	if error != nil {
		responses.JSON(w, http.StatusBadRequest, responses.ResponseError{Error: error.Error()})
		return
	}

	createUserUrl := fmt.Sprintf("%s/users", config.UrlApi)
	response, error := http.Post(createUserUrl, "application/json", bytes.NewBuffer(user))
	if error != nil {
		responses.JSON(w, http.StatusInternalServerError, responses.ResponseError{Error: error.Error()})
		return
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		responses.PrepareAndReturnError(w, response)
		return
	}

	responses.JSON(w, response.StatusCode, nil)
}

func UnfollowUserByUserId(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	userId, error := strconv.ParseUint(parameters["userId"], 10, 64)
	if error != nil {
		responses.JSON(w, http.StatusBadRequest, responses.ResponseError{Error: error.Error()})
		return
	}

	unfollowUrl := fmt.Sprintf("%s/users/%d/unfollow", config.UrlApi, userId)
	response, error := request.ExecuteRequestWithAuthentication(r, http.MethodPost, unfollowUrl, nil)
	if error != nil {
		responses.JSON(w, http.StatusInternalServerError, responses.ResponseError{Error: error.Error()})
		return
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		responses.PrepareAndReturnError(w, response)
		return
	}

	responses.JSON(w, response.StatusCode, nil)
}

func FollowUserByUserId(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	userId, error := strconv.ParseUint(parameters["userId"], 10, 64)
	if error != nil {
		responses.JSON(w, http.StatusBadRequest, responses.ResponseError{Error: error.Error()})
		return
	}

	followUrl := fmt.Sprintf("%s/users/%d/follow", config.UrlApi, userId)
	response, error := request.ExecuteRequestWithAuthentication(r, http.MethodPost, followUrl, nil)
	if error != nil {
		responses.JSON(w, http.StatusInternalServerError, responses.ResponseError{Error: error.Error()})
		return
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		responses.PrepareAndReturnError(w, response)
		return
	}

	responses.JSON(w, response.StatusCode, nil)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	cookies, _ := cookies.Read(r)
	loggedUserId, _ := strconv.ParseUint(cookies["id"], 10, 64)

	r.ParseForm()
	user, error := json.Marshal(map[string]string{
		"name":  r.FormValue("name"),
		"nick":  r.FormValue("nick"),
		"email": r.FormValue("email"),
	})
	if error != nil {
		responses.JSON(w, http.StatusBadRequest, responses.ResponseError{Error: error.Error()})
		return
	}

	userUrl := fmt.Sprintf("%s/users/%d", config.UrlApi, loggedUserId)
	response, error := request.ExecuteRequestWithAuthentication(r, http.MethodPut, userUrl, bytes.NewBuffer(user))
	if error != nil {
		responses.JSON(w, http.StatusInternalServerError, responses.ResponseError{Error: error.Error()})
		return
	}

	defer response.Body.Close()

	if response.StatusCode >= 400 {
		responses.PrepareAndReturnError(w, response)
		return
	}

	responses.JSON(w, response.StatusCode, nil)
}

func UpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	cookies, _ := cookies.Read(r)
	loggedUserId, _ := strconv.ParseUint(cookies["id"], 10, 64)

	r.ParseForm()
	userPasswords, error := json.Marshal(map[string]string{
		"current": r.FormValue("current"),
		"new":     r.FormValue("new"),
	})
	if error != nil {
		responses.JSON(w, http.StatusBadRequest, responses.ResponseError{Error: error.Error()})
		return
	}

	userUpdatePasswordUrl := fmt.Sprintf("%s/users/%d/update-password", config.UrlApi, loggedUserId)
	response, error := request.ExecuteRequestWithAuthentication(r, http.MethodPost, userUpdatePasswordUrl, bytes.NewBuffer(userPasswords))
	if error != nil {
		responses.JSON(w, http.StatusInternalServerError, responses.ResponseError{Error: error.Error()})
		return
	}

	defer response.Body.Close()

	if response.StatusCode >= 400 {
		responses.PrepareAndReturnError(w, response)
		return
	}

	responses.JSON(w, response.StatusCode, nil)
}

func DeleteUserById(w http.ResponseWriter, r *http.Request) {
	c, _ := cookies.Read(r)
	loggedUserId, _ := strconv.ParseUint(c["id"], 10, 64)

	deleteUserUrl := fmt.Sprintf("%s/users/%d", config.UrlApi, loggedUserId)
	response, error := request.ExecuteRequestWithAuthentication(r, http.MethodDelete, deleteUserUrl, nil)
	if error != nil {
		responses.JSON(w, http.StatusInternalServerError, responses.ResponseError{Error: error.Error()})
		return
	}

	defer response.Body.Close()

	if response.StatusCode >= 400 {
		responses.PrepareAndReturnError(w, response)
		return
	}

	cookies.ClearCookies(w)

	responses.JSON(w, response.StatusCode, nil)
}
