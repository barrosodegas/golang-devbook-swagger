package controller

import (
	"api/src/model"
	"bytes"
	"encoding/json"
	"fmt"
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

// TestFindUserByIdWithSuccess It guarantees that it will return the consulted user when a valid ID is informed.
func TestFindUserByIdWithSuccess(t *testing.T) {

	userId := 1

	req, _ := http.NewRequest("GET", fmt.Sprintf("/users/%d", userId), nil)
	req.Header.Set("Authorization", "Bearer "+userToken)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var user model.User
	if error := json.Unmarshal(response.Body.Bytes(), &user); error != nil {
		t.Errorf("Expected an User model. Got %s", response.Body)
	}

	if user.ID != uint64(userId) {
		t.Errorf("Expected an User model with ID: 1. Got %d", user.ID)
	}
}

// TestFindUserByIdWithBadRequestError
// It guarantees that the queried user will not be returned when a valid ID is not provided.
func TestFindUserByIdWithBadRequestError(t *testing.T) {

	userId := "d"

	req, _ := http.NewRequest("GET", fmt.Sprintf("/users/%s", userId), nil)
	req.Header.Set("Authorization", "Bearer "+userToken)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

// TestFindUserByIdWithUnauthorizedError
// It guarantees that it will not return the queried user when it is not an authenticated request.
func TestFindUserByIdWithUnauthorizedError(t *testing.T) {

	userId := 1

	req, _ := http.NewRequest("GET", fmt.Sprintf("/users/%d", userId), nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusUnauthorized, response.Code)
}

// TestFindUserByIdWithNotFoundError
// It guarantees that it will not return the queried user when the user does not exist.
func TestFindUserByIdWithNotFoundError(t *testing.T) {

	userId := 0

	req, _ := http.NewRequest("GET", fmt.Sprintf("/users/%d", userId), nil)
	req.Header.Set("Authorization", "Bearer "+userToken)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)
}

// TestUpdateUserByIdWithSuccess
// Ensures that the user will be updated correctly when a valid User ID is entered and the data for the update is valid.
func TestUpdateUserByIdWithSuccess(t *testing.T) {

	userId := 1

	// Find user after be updated
	req, _ := http.NewRequest("GET", fmt.Sprintf("/users/%d", userId), nil)
	req.Header.Set("Authorization", "Bearer "+userToken)
	response := executeRequest(req)

	var user model.User
	if error := json.Unmarshal(response.Body.Bytes(), &user); error != nil {
		t.Errorf("Expected an User model. Got %s", response.Body)
	}

	// Update user
	body := []byte(`{
		"name": "User Test - 1",
		"nick": "user@test-1",
		"email": "user.test1@gmail.com"
	}`)

	reqUpdate, _ := http.NewRequest("PUT", fmt.Sprintf("/users/%d", userId), bytes.NewBuffer(body))
	reqUpdate.Header.Set("Authorization", "Bearer "+userToken)
	responseUpdate := executeRequest(reqUpdate)

	checkResponseCode(t, http.StatusNoContent, responseUpdate.Code)

	// Find user after be updated
	response = executeRequest(req)

	var updatedUser model.User
	if error := json.Unmarshal(response.Body.Bytes(), &updatedUser); error != nil {
		t.Errorf("Expected an updated User model. Got %s", response.Body)
	}

	if updatedUser.Name != "User Test - 1" && updatedUser.Name == user.Name {
		t.Errorf("Expected an updated User model with name: User Test - 1. Got %s", updatedUser.Name)
	}

	if updatedUser.Nick != "user@test-1" && updatedUser.Nick == user.Nick {

		t.Errorf("Expected an updated User model with nick: user@test-1. Got %s", updatedUser.Nick)
	}
	if updatedUser.Email != "user.test1@gmail.com" && updatedUser.Email == user.Email {
		t.Errorf("Expected an updated User model with name: user.test1@gmail.com. Got %s", updatedUser.Email)
	}
}

// TestUpdateUserByIdWithBadRequestErrorWhenUserIdIsInvalid
// It guarantees that the user will not be updated when a valid User ID is not provided.
func TestUpdateUserByIdWithBadRequestErrorWhenUserIdIsInvalid(t *testing.T) {

	userId := "d"

	body := []byte(`{
		"name": "User Test - 1",
		"nick": "user@test-1",
		"email": "user.test1@gmail.com"
	}`)

	reqUpdate, _ := http.NewRequest("PUT", fmt.Sprintf("/users/%s", userId), bytes.NewBuffer(body))
	reqUpdate.Header.Set("Authorization", "Bearer "+userToken)
	responseUpdate := executeRequest(reqUpdate)

	checkResponseCode(t, http.StatusBadRequest, responseUpdate.Code)
}

// TestUpdateUserByIdWithBadRequestErrorWhenBodyDataIsInvalid
// It guarantees that the user will not be updated when a valid request body is not informed.
func TestUpdateUserByIdWithBadRequestErrorWhenBodyDataIsInvalid(t *testing.T) {

	userId := 1

	body := []byte(``)

	reqUpdate, _ := http.NewRequest("PUT", fmt.Sprintf("/users/%d", userId), bytes.NewBuffer(body))
	reqUpdate.Header.Set("Authorization", "Bearer "+userToken)
	responseUpdate := executeRequest(reqUpdate)

	checkResponseCode(t, http.StatusBadRequest, responseUpdate.Code)
}

// TestUpdateUserByIdWithUnauthorizedError
// It guarantees that the user will not be updated when the request is not authenticated.
func TestUpdateUserByIdWithUnauthorizedError(t *testing.T) {

	userId := 1

	body := []byte(`{
		"name": "User Test - 1",
		"nick": "user@test-1",
		"email": "user.test1@gmail.com"
	}`)

	reqUpdate, _ := http.NewRequest("PUT", fmt.Sprintf("/users/%d", userId), bytes.NewBuffer(body))
	responseUpdate := executeRequest(reqUpdate)

	checkResponseCode(t, http.StatusUnauthorized, responseUpdate.Code)
}

// TestUpdateUserByIdWithForbiddenError
// It guarantees that the user will not be updated when the user ID in the request is different from the logged in user ID.
func TestUpdateUserByIdWithForbiddenError(t *testing.T) {

	userId := 2

	body := []byte(`{
		"name": "User Test - 1",
		"nick": "user@test-1",
		"email": "user.test1@gmail.com"
	}`)

	reqUpdate, _ := http.NewRequest("PUT", fmt.Sprintf("/users/%d", userId), bytes.NewBuffer(body))
	reqUpdate.Header.Set("Authorization", "Bearer "+userToken)
	responseUpdate := executeRequest(reqUpdate)

	checkResponseCode(t, http.StatusForbidden, responseUpdate.Code)
}

// TestDeleteUserByIdWithSuccess
// Ensures the user will be deleted when the user ID is valid and the user is authenticated.
func TestDeleteUserByIdWithSuccess(t *testing.T) {

	userId := 1

	// Delete user
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/users/%d", userId), nil)
	req.Header.Set("Authorization", "Bearer "+userToken)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNoContent, response.Code)

	// Find user after be delete and not found user.
	req, _ = http.NewRequest("GET", fmt.Sprintf("/users/%d", userId), nil)
	req.Header.Set("Authorization", "Bearer "+userToken)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)
}

// TestDeleteUserByIdWithBadRequestError
// Ensures the user will not be deleted when the user ID is invalid.
func TestDeleteUserByIdWithBadRequestError(t *testing.T) {

	userId := "d"

	// Find user after be try deleted
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/users/%s", userId), nil)
	req.Header.Set("Authorization", "Bearer "+userToken)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

// TestDeleteUserByIdWithUnauthorizedError
// Ensures that the user will not be deleted when the user ID is valid and the request is not authenticated.
func TestDeleteUserByIdWithUnauthorizedError(t *testing.T) {

	userId := 1

	// Try delete user
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/users/%d", userId), nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusUnauthorized, response.Code)
}

// TestDeleteUserByIdWithForbiddenError
// It guarantees that the user will not be deleted when the user ID is valid and different from the logged in user ID.
func TestDeleteUserByIdWithForbiddenError(t *testing.T) {

	userId := 2

	// Try delete user
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/users/%d", userId), nil)
	req.Header.Set("Authorization", "Bearer "+userToken)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusForbidden, response.Code)
}
