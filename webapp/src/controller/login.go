package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"webapp/src/config"
	"webapp/src/cookies"
	"webapp/src/model"
	"webapp/src/responses"
)

func Login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	user, error := json.Marshal(map[string]string{
		"email":    r.FormValue("email"),
		"password": r.FormValue("password"),
	})
	if error != nil {
		responses.JSON(w, http.StatusBadRequest, responses.ResponseError{Error: error.Error()})
		return
	}

	loginUrl := fmt.Sprintf("%s/login", config.UrlApi)
	response, error := http.Post(loginUrl, "application/json", bytes.NewBuffer(user))
	if error != nil {
		responses.JSON(w, http.StatusInternalServerError, responses.ResponseError{Error: error.Error()})
		return
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		responses.PrepareAndReturnError(w, response)
		return
	}

	var dataAuthentication model.DataAuthenticaion
	if error = json.NewDecoder(response.Body).Decode(&dataAuthentication); error != nil {
		responses.JSON(w, http.StatusUnprocessableEntity, responses.ResponseError{Error: error.Error()})
		return
	}

	if error = cookies.SaveCookies(w, dataAuthentication.ID, dataAuthentication.Token); error != nil {
		responses.JSON(w, http.StatusUnprocessableEntity, responses.ResponseError{Error: error.Error()})
		return
	}

	responses.JSON(w, http.StatusOK, nil)
}
