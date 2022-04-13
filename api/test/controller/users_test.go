package controller

import (
	"api/src/model"
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

// TestCreateUserSuccess Ensures that a user is created correctly when required data is entered.
func TestCreateUserWithSuccess(t *testing.T) {

	body := []byte(`{
		"name": "User Test",
		"nick": "user@test",
		"email": "user.test@gmail.com",
		"password": "test"
	}`)

	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var newUser model.User
	if error := json.Unmarshal(response.Body.Bytes(), &newUser); error != nil {
		t.Errorf("Expected an new User model. Got %s", body)
	}

	if newUser.Name != "User Test" {
		t.Errorf("Expected an new User model with name: User Test. Got %s", newUser.Name)
	}
}

// TestCreateUserWithBadRequestError Ensures that a user will not be created when the request body is invalid.
func TestCreateUserWithBadRequestError(t *testing.T) {

	body := []byte(``)

	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
}
