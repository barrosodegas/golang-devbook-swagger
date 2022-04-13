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
		t.Errorf("Expected an new User model. Got %s", response.Body)
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

// TestListUsersByFilterByNameWithSuccess
// It guarantees that it will return users that contain the same name or part of the name as entered.
func TestListUsersByFilterByNameWithSuccess(t *testing.T) {

	userName := "User Test"

	req, _ := http.NewRequest("GET", "/users?user="+userName, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+userToken)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var users []model.User
	if error := json.Unmarshal(response.Body.Bytes(), &users); error != nil {
		t.Errorf("Expected a user list. Got %s", response.Body)
	}

	if len(users) == 0 {
		t.Errorf("Expected a user list. Got a empty list")
	}
}

// TestListUsersByFilterByNickWithSuccess
// It guarantees that users that contain the nick or part of the nick equal to the one informed will be returned.
func TestListUsersByFilterByNickWithSuccess(t *testing.T) {

	userNick := "user@test"

	req, _ := http.NewRequest("GET", "/users?user="+userNick, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+userToken)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var users []model.User
	if error := json.Unmarshal(response.Body.Bytes(), &users); error != nil {
		t.Errorf("Expected a user list. Got %s", response.Body)
	}

	if len(users) == 0 {
		t.Errorf("Expected a user list. Got a empty list")
	}
}

// TestListUsersByFilterByEmptyParamWithSuccess It guarantees that it will return all users.
func TestListUsersByFilterByEmptyParamWithSuccess(t *testing.T) {

	emptyParam := ""

	req, _ := http.NewRequest("GET", "/users?user="+emptyParam, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+userToken)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var users []model.User
	if error := json.Unmarshal(response.Body.Bytes(), &users); error != nil {
		t.Errorf("Expected a user list. Got %s", response.Body)
	}

	if len(users) < 1 {
		t.Error("Expected a populated list. Got a empty list")
	}
}

// TestListUsersByFilterWithUnauthorizedError
// It guarantees that it will not return users when the request is not authenticated.
func TestListUsersByFilterWithUnauthorizedError(t *testing.T) {

	userName := "User Test"

	req, _ := http.NewRequest("GET", "/users?user="+userName, nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusUnauthorized, response.Code)
}
