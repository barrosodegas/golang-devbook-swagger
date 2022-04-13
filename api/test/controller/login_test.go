package controller

import (
	"api/src/model"
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

// TestLoginWithSuccess Ensures that the login is successful when the user enters a correct email and password.
func TestLoginWithSuccess(t *testing.T) {

	body := []byte(`{"email": "user1@gmail.com", "password": "alb1234"}`)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var dataAuthentication model.DataAuthentication
	if error := json.Unmarshal(response.Body.Bytes(), &dataAuthentication); error != nil {
		t.Errorf("Expected an DataAuthentication model. Got %s", body)
	}
}

// TestLoginWithBadRequestError Ensures login will not be successful when user submits invalid json.
func TestLoginWithBadRequestError(t *testing.T) {

	body := []byte(``)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

// TestLoginWithUnauthorizedError It guarantees that the login will not be successful when the user does not provide a correct email and password.
func TestLoginWithUnauthorizedError(t *testing.T) {

	body := []byte(`{"email": "user1@gmail.com", "password": "alb12"}`)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusUnauthorized, response.Code)
}
