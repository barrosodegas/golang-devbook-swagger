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
