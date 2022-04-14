package controller

import (
	"api/src/model"
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

// TestCreatePublicationWithSuccess
// Ensures that a publication will be created for the logged in user when the request is
// authenticated and the request body is valid.
func TestCreatePublicationWithSuccess(t *testing.T) {

	body := []byte(`{
		"title": "Pub Test - 1",
		"content": "Pub to test - 1"
	}`)

	req, _ := http.NewRequest("POST", "/publications", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+userToken)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var newPublication model.Publication
	if error := json.Unmarshal(response.Body.Bytes(), &newPublication); error != nil {
		t.Errorf("Expected an new Publication model. Got %s", response.Body)
	}

	if newPublication.ID != 4 {
		t.Errorf("Expected an new Publication model with ID: 4. Got %d", newPublication.ID)
	}
}

// TestCreatePublicationWithUnauthorizedError
// It guarantees that a publication will not be created when the request is not authenticated.
func TestCreatePublicationWithUnauthorizedError(t *testing.T) {

	body := []byte(`{
		"title": "Pub Test - 1",
		"content": "Pub to test - 1"
	}`)

	req, _ := http.NewRequest("POST", "/publications", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	response := executeRequest(req)

	checkResponseCode(t, http.StatusUnauthorized, response.Code)
}

// TestCreatePublicationWithBadRequestErrorWhenIsInvalidJson
// It guarantees that a publication will not be created when the request is authenticated but the request body is invalid JSON.
func TestCreatePublicationWithBadRequestErrorWhenIsInvalidJson(t *testing.T) {

	body := []byte(``)

	req, _ := http.NewRequest("POST", "/publications", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+userToken)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

// TestCreatePublicationWithBadRequestErrorWhenAFieldIsMandatory
// It guarantees that a publication will not be created when the request is authenticated, the request body
//is valid and a mandatory field is not informed.
func TestCreatePublicationWithBadRequestErrorWhenAFieldIsMandatory(t *testing.T) {

	body := []byte(`{
		"title": "Pub Test - 1"
	}`)

	req, _ := http.NewRequest("POST", "/publications", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+userToken)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

// TestListMyAndFollowPublicationsWithSuccess
// Ensures a list containing the publications of the logged in user and the publications of the users he follows.
func TestListMyAndFollowPublicationsWithSuccess(t *testing.T) {

	req, _ := http.NewRequest("GET", "/publications", nil)
	req.Header.Set("Authorization", "Bearer "+userToken)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var publications []model.Publication
	if error := json.Unmarshal(response.Body.Bytes(), &publications); error != nil {
		t.Errorf("Expected a publication list. Got %s", response.Body)
	}

	if len(publications) < 1 {
		t.Errorf("Expected a Publication list. Got a empty list.")
	}

	// Follower user
	var hasPublicationToUser1 bool
	// Followed user
	var hasPublicationToUser2 bool

	for _, p := range publications {
		if p.ID == uint64(1) {
			hasPublicationToUser1 = true
		}

		if p.ID == uint64(2) {
			hasPublicationToUser2 = true
		}
	}

	if !hasPublicationToUser1 || !hasPublicationToUser2 {
		t.Errorf("Expected a Publication list contains my and follow publications. Got %q", publications)
	}
}

// TestListMyAndFollowPublicationsWithEmptyListSuccess
// Guarantees that an empty list will be returned when the user has no publications.
func TestListMyAndFollowPublicationsWithEmptyListSuccess(t *testing.T) {

	// Create a new user withou publications.
	body := []byte(`{
		"name": "User Test 2",
		"nick": "user@test2",
		"email": "user.test2@gmail.com",
		"password": "test2"
	}`)

	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	response := executeRequest(req)

	// Get token of new user
	token := GetUserToken("user.test2@gmail.com", "test2")

	// Get publications
	req, _ = http.NewRequest("GET", "/publications", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var publications []model.Publication
	if error := json.Unmarshal(response.Body.Bytes(), &publications); error != nil {
		t.Errorf("Expected a publication list. Got %s", response.Body)
	}

	if len(publications) != 0 {
		t.Errorf("Expected a empty Publication list. Got %q.", publications)
	}
}

// TestListMyAndFollowPublicationsWithUnauthorizedError
// It guarantees that it will not return the list of publications when the request is not authenticated.
func TestListMyAndFollowPublicationsWithUnauthorizedError(t *testing.T) {

	req, _ := http.NewRequest("GET", "/publications", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusUnauthorized, response.Code)
}
