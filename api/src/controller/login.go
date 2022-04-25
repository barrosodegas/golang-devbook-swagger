package controller

import (
	"api/src/authentication"
	"api/src/database"
	"api/src/model"
	"api/src/repository"
	"api/src/responses"
	"api/src/security"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Login 	 	godoc
//
// @Summary     Generates the system access token.
// @Description Generates the system access token and return an object containing the system access token and the user ID.
// @Description Use the token generated at login to run the other API endpoints through Authorization Bearer.
//
// @Tags login
//
// @Accept  json
// @Produce json
//
// @Param user body model.User true "Enter only the user's email and password."
//
// @Success 200 {object} model.DataAuthentication "Success"
//
// @Failure 400 {string}  string "Error: Bad Request"
// @Failure 401 {string}  string "Error: Unauthorized"
// @Failure 422 {string}  string "Error: Unprocessable Entity"
// @Failure 500 {string}  string "Error: Internal Server"
//
// @Router /login [post]
func Login(w http.ResponseWriter, r *http.Request) {
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

	db, error := database.Connect()
	if error != nil {
		responses.Error(w, http.StatusInternalServerError, error)
		return
	}

	defer db.Close()

	repository := repository.NewUsersRepository(db)

	entity, error := repository.FindUserByEmail(user.Email)
	if error != nil {
		responses.Error(w, http.StatusInternalServerError, error)
		return
	}

	if error = security.VerifyPassword(entity.Password, user.Password); error != nil {
		responses.Error(w, http.StatusUnauthorized, error)
		return
	}

	token, error := authentication.CreateToken(entity.ID)
	if error != nil {
		responses.Error(w, http.StatusInternalServerError, error)
		return
	}

	userId := strconv.FormatUint(entity.ID, 10)

	responses.Success(w, http.StatusOK, model.DataAuthentication{ID: userId, Token: token})
}
