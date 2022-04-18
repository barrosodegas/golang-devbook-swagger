package controller

import (
	"api/src/model"
	"encoding/json"
	"net/http"
	"testing"
)

// TestLoginWithSuccess Ensures that the login is successful when the user enters a correct email and password.
func TestLoginWithSuccess(t *testing.T) {

	body := []byte(`{"email": "user1@gmail.com", "password": "alb1234"}`)

	response := executeRequest("POST", "/login", "", body, true)

	checkResponseCode(t, http.StatusOK, response.Code)

	var dataAuthentication model.DataAuthentication
	if error := json.Unmarshal(response.Body.Bytes(), &dataAuthentication); error != nil {
		t.Errorf("Expected an DataAuthentication model. Got %s", response.Body)
	}
}

// TestLoginWithBadRequestError Ensures login will not be successful when user submits invalid json.
func TestLoginWithBadRequestError(t *testing.T) {

	body := []byte(``)

	response := executeRequest("POST", "/login", "", body, true)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

// TestLoginWithUnauthorizedError It guarantees that the login will not be successful when the user does not provide a correct email and password.
func TestLoginWithUnauthorizedError(t *testing.T) {

	body := []byte(`{"email": "user1@gmail.com", "password": "alb12"}`)

	response := executeRequest("POST", "/login", "", body, true)

	checkResponseCode(t, http.StatusUnauthorized, response.Code)
}
