package controller

import (
	"api/src/authentication"
	"api/src/database"
	"api/src/model"
	"api/src/repository"
	"api/src/responses"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Publication godoc
//
// @Summary     CreatePublication create a publication.
// @Description CreatePublication create a publication and returns the publication created or an error.
//
// @Tags publications
//
// @Accept  json
// @Produce json
//
// @Param Authorization header string            true "Enter the content: Bearer and your access token."
// @Param publication   body   model.Publication true "Enter only the title and content of the publication."
//
// @Success 201 {object} model.Publication "Success"
//
// @Failure 400 {string}  string "Error: Bad Request"
// @Failure 401 {string}  string "Error: Unauthorized"
// @Failure 422 {string}  string "Error: Unprocessable Entity"
// @Failure 500 {string}  string "Error: Internal Server"
//
// @Router /publications [post]
func CreatePublication(w http.ResponseWriter, r *http.Request) {
	loggedUserId, error := authentication.ExtractUserIdOfToken(r)
	if error != nil {
		responses.Error(w, http.StatusUnauthorized, error)
		return
	}

	requestBody, error := ioutil.ReadAll(r.Body)
	if error != nil {
		responses.Error(w, http.StatusUnprocessableEntity, error)
		return
	}

	var publication model.Publication

	if error = json.Unmarshal(requestBody, &publication); error != nil {
		responses.Error(w, http.StatusBadRequest, error)
		return
	}

	if error = publication.Prepare(); error != nil {
		responses.Error(w, http.StatusBadRequest, error)
		return
	}

	db, error := database.Connect()
	if error != nil {
		responses.Error(w, http.StatusInternalServerError, error)
		return
	}

	publication.AuthorId = loggedUserId

	defer db.Close()

	repository := repository.NewPublicationsRepository(db)

	publication.ID, error = repository.CreatePublication(publication)
	if error != nil {
		responses.Error(w, http.StatusInternalServerError, error)
		return
	}

	responses.Success(w, http.StatusCreated, publication)
}

// Publication godoc
//
// @Summary     ListMyAndFollowPublications lists the publications of the logged in user and the publications they follow.
// @Description ListMyAndFollowPublications lists the publications of the logged in user and the publications they follow
// @Description and returns a list containing the posts of the logged in user and the posts of whom he follows or an error.
//
// @Tags publications
//
// @Accept  json
// @Produce json
//
// @Param Authorization header string true "Enter the content: Bearer and your access token."
//
// @Success 200 {array} model.Publication "Success"
//
// @Failure 401 {string}  string "Error: Unauthorized"
// @Failure 500 {string}  string "Error: Internal Server"
//
// @Router /publications [get]
func ListMyAndFollowPublications(w http.ResponseWriter, r *http.Request) {
	loggedUserId, error := authentication.ExtractUserIdOfToken(r)
	if error != nil {
		responses.Error(w, http.StatusUnauthorized, error)
		return
	}

	db, error := database.Connect()
	if error != nil {
		responses.Error(w, http.StatusInternalServerError, error)
		return
	}

	defer db.Close()

	repository := repository.NewPublicationsRepository(db)

	publications, error := repository.ListMyAndFollowPublications(loggedUserId)
	if error != nil {
		responses.Error(w, http.StatusInternalServerError, error)
		return
	}

	responses.Success(w, http.StatusOK, publications)
}

// Publication godoc
//
// @Summary     FindPublicationById search for a publication by the given ID.
// @Description FindPublicationById search for a publication by the given ID
// @Description and returns the publication or an error if unable to create the publication.
//
// @Tags publications
//
// @Accept  json
// @Produce json
//
// @Param Authorization header string true "Enter the content: Bearer and your access token."
// @Param publicationId path   int    true "Publication ID"
//
// @Success 200 {object} model.Publication "Success"
//
// @Failure 400 {string}  string "Error: Bad Request"
// @Failure 401 {string}  string "Error: Unauthorized"
// @Failure 404 {string}  string "Error: Not Found"
// @Failure 500 {string}  string "Error: Internal Server"
//
// @Router /publications/{publicationId} [get]
func FindPublicationById(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)

	publicationId, error := strconv.ParseUint(parameters["publicationId"], 10, 64)
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

	repository := repository.NewPublicationsRepository(db)

	publication, error := repository.FindPublicationById(publicationId)
	if error != nil {
		responses.Error(w, http.StatusInternalServerError, error)
		return
	}
	if publication.ID == 0 {
		responses.Error(w, http.StatusNotFound, errors.New("Publication not found!"))
		return
	}

	responses.Success(w, http.StatusOK, publication)
}

// Publication godoc
//
// @Summary     UpdatePublicationById updates a publication by the given ID..
// @Description UpdatePublicationById updates a publication by the given ID
// @Description and returns an error if unable to update the publication.
//
// @Tags publications
//
// @Accept  json
// @Produce json
//
// @Param Authorization header string            true "Enter the content: Bearer and your access token."
// @Param publicationId path   int               true "Publication ID"
// @Param publication   body   model.Publication true "Enter only the title and content of the publication."
//
// @Success 204 "Success with no content"
//
// @Failure 400 {string} string "Error: Bad Request"
// @Failure 401 {string} string "Error: Unauthorized"
// @Failure 403 {string} string "Error: Forbidden"
// @Failure 422 {string} string "Error: Unprocessable Entity"
// @Failure 500 {string} string "Error: Internal Server"
//
// @Router /publications/{publicationId} [put]
func UpdatePublicationById(w http.ResponseWriter, r *http.Request) {
	loggedUserId, error := authentication.ExtractUserIdOfToken(r)
	if error != nil {
		responses.Error(w, http.StatusUnauthorized, error)
		return
	}

	parameters := mux.Vars(r)

	publicationId, error := strconv.ParseUint(parameters["publicationId"], 10, 64)
	if error != nil {
		responses.Error(w, http.StatusBadRequest, error)
		return
	}

	requestBody, error := ioutil.ReadAll(r.Body)
	if error != nil {
		responses.Error(w, http.StatusUnprocessableEntity, error)
		return
	}

	var publication model.Publication

	if error = json.Unmarshal(requestBody, &publication); error != nil {
		responses.Error(w, http.StatusBadRequest, error)
		return
	}

	if error = publication.Prepare(); error != nil {
		responses.Error(w, http.StatusBadRequest, error)
		return
	}

	db, error := database.Connect()
	if error != nil {
		responses.Error(w, http.StatusInternalServerError, error)
		return
	}

	defer db.Close()

	repository := repository.NewPublicationsRepository(db)

	publicationEntity, error := repository.FindPublicationById(publicationId)
	if error != nil {
		responses.Error(w, http.StatusInternalServerError, error)
		return
	}

	if publicationEntity.AuthorId != loggedUserId {
		responses.Error(w, http.StatusForbidden, errors.New("Invalid action!"))
		return
	}

	if error = repository.UpdatePublicationById(publicationId, publication); error != nil {
		responses.Error(w, http.StatusInternalServerError, error)
		return
	}

	responses.Success(w, http.StatusNoContent, nil)
}

// Publication godoc
//
// @Summary     DeletePublicationById deletes a publication by the given ID.
// @Description DeletePublicationById deletes a publication by the given ID
// @Description and returns an error if unable to delete the publication.
//
// @Tags publications
//
// @Accept  json
// @Produce json
//
// @Param Authorization header string true "Enter the content: Bearer and your access token."
// @Param publicationId path   int    true "Publication ID"
//
// @Success 204 "Success with no content"
//
// @Failure 400 {string} string "Error: Bad Request"
// @Failure 401 {string} string "Error: Unauthorized"
// @Failure 403 {string} string "Error: Forbidden"
// @Failure 500 {string} string "Error: Internal Server"
//
// @Router /publications/{publicationId} [delete]
func DeletePublicationById(w http.ResponseWriter, r *http.Request) {
	loggedUserId, error := authentication.ExtractUserIdOfToken(r)
	if error != nil {
		responses.Error(w, http.StatusUnauthorized, error)
		return
	}

	parameters := mux.Vars(r)

	publicationId, error := strconv.ParseUint(parameters["publicationId"], 10, 64)
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

	repository := repository.NewPublicationsRepository(db)

	publicationEntity, error := repository.FindPublicationById(publicationId)
	if error != nil {
		responses.Error(w, http.StatusInternalServerError, error)
		return
	}

	if publicationEntity.AuthorId != loggedUserId {
		responses.Error(w, http.StatusForbidden, errors.New("Invalid action!"))
		return
	}

	if error = repository.DeletePublicationById(publicationId); error != nil {
		responses.Error(w, http.StatusInternalServerError, error)
		return
	}

	responses.Success(w, http.StatusNoContent, nil)
}

// Publication godoc
//
// @Summary     ListPublicationsByUserId lists a user's publications.
// @Description ListPublicationsByUserId lists a user's publications
// @Description and returns a list of publications or an error.
//
// @Tags publications
//
// @Accept  json
// @Produce json
//
// @Param Authorization header string true "Enter the content: Bearer and your access token."
// @Param userId        path   int    true "User ID"
//
// @Success 200 {array} model.Publication "Success"
//
// @Failure 400 {string} string "Error: Bad Request"
// @Failure 401 {string} string "Error: Unauthorized"
// @Failure 403 {string} string "Error: Forbidden"
// @Failure 500 {string} string "Error: Internal Server"
//
// @Router /publications/user/{userId} [get]
func ListPublicationsByUserId(w http.ResponseWriter, r *http.Request) {

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

	repository := repository.NewPublicationsRepository(db)

	publications, error := repository.ListPublicationsByUserId(userId)
	if error != nil {
		responses.Error(w, http.StatusInternalServerError, error)
		return
	}

	responses.Success(w, http.StatusOK, publications)
}

// Publication godoc
//
// @Summary     LikePublicationById likes a publication by the given ID.
// @Description LikePublicationById likes a publication by the given ID
// @Description and returns an error if unable to like the publication.
//
// @Tags publications
//
// @Accept  json
// @Produce json
//
// @Param Authorization header string true "Enter the content: Bearer and your access token."
// @Param publicationId path   int    true "Publication ID"
//
// @Success 204 "Success with no content"
//
// @Failure 400 {string} string "Error: Bad Request"
// @Failure 401 {string} string "Error: Unauthorized"
// @Failure 500 {string} string "Error: Internal Server"
//
// @Router /publications/{publicationId}/like [post]
func LikePublicationById(w http.ResponseWriter, r *http.Request) {

	parameters := mux.Vars(r)

	publicationId, error := strconv.ParseUint(parameters["publicationId"], 10, 64)
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

	repository := repository.NewPublicationsRepository(db)

	error = repository.LikePublicationById(publicationId)
	if error != nil {
		responses.Error(w, http.StatusInternalServerError, error)
		return
	}

	responses.Success(w, http.StatusNoContent, nil)
}

// Publication godoc
//
// @Summary     UnlikePublicationById unlike a publication by the given ID.
// @Description UnlikePublicationById unlike a publication by the given ID
// @Description and returns an error if unable to unlike the publication.
//
// @Tags publications
//
// @Accept  json
// @Produce json
//
// @Param Authorization header string true "Enter the content: Bearer and your access token."
// @Param publicationId path   int    true "Publication ID"
//
// @Success 204 "Success with no content"
//
// @Failure 400 {string} string "Error: Bad Request"
// @Failure 401 {string} string "Error: Unauthorized"
// @Failure 500 {string} string "Error: Internal Server"
//
// @Router /publications/{publicationId}/unlike [post]
func UnlikePublicationById(w http.ResponseWriter, r *http.Request) {

	parameters := mux.Vars(r)

	publicationId, error := strconv.ParseUint(parameters["publicationId"], 10, 64)
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

	repository := repository.NewPublicationsRepository(db)

	error = repository.UnLikePublicationById(publicationId)
	if error != nil {
		responses.Error(w, http.StatusInternalServerError, error)
		return
	}

	responses.Success(w, http.StatusNoContent, nil)
}
