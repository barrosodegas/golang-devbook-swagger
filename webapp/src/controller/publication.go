package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"webapp/src/config"
	"webapp/src/request"
	"webapp/src/responses"

	"github.com/gorilla/mux"
)

func CreatePublication(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	publication, error := json.Marshal(map[string]string{
		"title":   r.FormValue("title"),
		"content": r.FormValue("content"),
	})
	if error != nil {
		responses.JSON(w, http.StatusBadRequest, responses.ResponseError{Error: error.Error()})
		return
	}

	publicationUrl := fmt.Sprintf("%s/publications", config.UrlApi)
	response, error := request.ExecuteRequestWithAuthentication(r, http.MethodPost, publicationUrl, bytes.NewBuffer(publication))
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

func LikePublication(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	publicationId, error := strconv.ParseUint(parameters["publicationId"], 10, 64)
	if error != nil {
		responses.JSON(w, http.StatusBadRequest, responses.ResponseError{Error: error.Error()})
		return
	}

	likePublicationUrl := fmt.Sprintf("%s/publications/%d/like", config.UrlApi, publicationId)
	response, error := request.ExecuteRequestWithAuthentication(r, http.MethodPost, likePublicationUrl, nil)
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

func UnlikePublication(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	publicationId, error := strconv.ParseUint(parameters["publicationId"], 10, 64)
	if error != nil {
		responses.JSON(w, http.StatusBadRequest, responses.ResponseError{Error: error.Error()})
		return
	}

	unlikePublicationUrl := fmt.Sprintf("%s/publications/%d/unlike", config.UrlApi, publicationId)
	response, error := request.ExecuteRequestWithAuthentication(r, http.MethodPost, unlikePublicationUrl, nil)
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

func UpdatePublication(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	publicationId, error := strconv.ParseUint(parameters["publicationId"], 10, 64)
	if error != nil {
		responses.JSON(w, http.StatusBadRequest, responses.ResponseError{Error: error.Error()})
		return
	}

	r.ParseForm()

	publication, error := json.Marshal(map[string]string{
		"title":   r.FormValue("title"),
		"content": r.FormValue("content"),
	})
	if error != nil {
		responses.JSON(w, http.StatusBadRequest, responses.ResponseError{Error: error.Error()})
		return
	}

	updatePublicationUrl := fmt.Sprintf("%s/publications/%d", config.UrlApi, publicationId)
	response, error := request.ExecuteRequestWithAuthentication(r, http.MethodPut, updatePublicationUrl, bytes.NewBuffer(publication))
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

func DeletePublication(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	publicationId, error := strconv.ParseUint(parameters["publicationId"], 10, 64)
	if error != nil {
		responses.JSON(w, http.StatusBadRequest, responses.ResponseError{Error: error.Error()})
		return
	}

	deletePublicationUrl := fmt.Sprintf("%s/publications/%d", config.UrlApi, publicationId)
	response, error := request.ExecuteRequestWithAuthentication(r, http.MethodDelete, deletePublicationUrl, nil)
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
