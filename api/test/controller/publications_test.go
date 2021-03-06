package controller

import (
	"api/src/model"
	"encoding/json"
	"fmt"
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

	response := executeRequest("POST", "/publications", userToken, body, true)

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

	response := executeRequest("POST", "/publications", "", body, true)

	checkResponseCode(t, http.StatusUnauthorized, response.Code)
}

// TestCreatePublicationWithBadRequestErrorWhenIsInvalidJson
// It guarantees that a publication will not be created when the request is authenticated but the request body is invalid JSON.
func TestCreatePublicationWithBadRequestErrorWhenIsInvalidJson(t *testing.T) {

	body := []byte(``)

	response := executeRequest("POST", "/publications", userToken, body, true)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

// TestCreatePublicationWithBadRequestErrorWhenAFieldIsMandatory
// It guarantees that a publication will not be created when the request is authenticated, the request body
//is valid and a mandatory field is not informed.
func TestCreatePublicationWithBadRequestErrorWhenAFieldIsMandatory(t *testing.T) {

	body := []byte(`{
		"title": "Pub Test - 1"
	}`)

	response := executeRequest("POST", "/publications", userToken, body, true)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

// TestListMyAndFollowPublicationsWithSuccess
// Ensures a list containing the publications of the logged in user and the publications of the users he follows.
func TestListMyAndFollowPublicationsWithSuccess(t *testing.T) {

	response := executeRequest("GET", "/publications", userToken, nil, false)

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

	// Create a new user without publications.
	body := []byte(`{
		"name": "User Test 2",
		"nick": "user@test2",
		"email": "user.test2@gmail.com",
		"password": "test2"
	}`)

	response := executeRequest("POST", "/users", "", body, true)

	// Get token of new user
	token := GetUserToken("user.test2@gmail.com", "test2")

	// Get publications
	response = executeRequest("GET", "/publications", token, nil, false)

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

	response := executeRequest("GET", "/publications", "", nil, false)

	checkResponseCode(t, http.StatusUnauthorized, response.Code)
}

// TestFindPublicationByIdSuccess
// It guarantees that a publication will be returned when the request is authenticated and the user entered is valid.
func TestFindPublicationByIdSuccess(t *testing.T) {

	var publicationId = 1

	response := executeRequest("GET", fmt.Sprintf("/publications/%d", publicationId), userToken, nil, false)

	checkResponseCode(t, http.StatusOK, response.Code)

	var publication model.Publication
	if error := json.Unmarshal(response.Body.Bytes(), &publication); error != nil {
		t.Errorf("Expected a publication. Got %s", response.Body)
	}

	if publication.ID == 0 {
		t.Errorf("Expected a Publication with ID: 1. Got %d.", publication.ID)
	}
}

// TestFindPublicationByIdBadRequestError
// It guarantees that a publication will not be returned when the request is authenticated and the user entered is invalid.
func TestFindPublicationByIdBadRequestError(t *testing.T) {

	var publicationId = "d"

	response := executeRequest("GET", fmt.Sprintf("/publications/%s", publicationId), userToken, nil, false)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

// TestFindPublicationByIdUnauthorizedError
// It guarantees that a post will not be returned when the request is not authenticated.
func TestFindPublicationByIdUnauthorizedError(t *testing.T) {

	var publicationId = 1

	response := executeRequest("GET", fmt.Sprintf("/publications/%d", publicationId), "", nil, false)

	checkResponseCode(t, http.StatusUnauthorized, response.Code)
}

// TestFindPublicationByIdNotFoundError
// Guarantees that a post will not be returned when the request is authenticated and the reported post does not exist.
func TestFindPublicationByIdNotFoundError(t *testing.T) {

	var publicationId = 1000

	response := executeRequest("GET", fmt.Sprintf("/publications/%d", publicationId), userToken, nil, false)

	checkResponseCode(t, http.StatusNotFound, response.Code)
}

// TestUpdatePublicationByIdWithSuccess
// It guarantees that a publication will be updated when the given id is valid, the request
// is authenticated and the request body is valid.
func TestUpdatePublicationByIdWithSuccess(t *testing.T) {

	var publicationId = 1

	// Find publication after is updated.
	response := executeRequest("GET", fmt.Sprintf("/publications/%d", publicationId), userToken, nil, false)

	checkResponseCode(t, http.StatusOK, response.Code)

	var publicationOld model.Publication
	if error := json.Unmarshal(response.Body.Bytes(), &publicationOld); error != nil {
		t.Errorf("Expected a publication. Got %s", response.Body)
	}

	// Update the publication
	body := []byte(`{
		"title": "Pub Test - 1",
		"content": "Pub to test - 1..."
	}`)

	response = executeRequest("PUT", fmt.Sprintf("/publications/%d", publicationId), userToken, body, true)

	checkResponseCode(t, http.StatusNoContent, response.Code)

	// Find the updated publication.
	response = executeRequest("GET", fmt.Sprintf("/publications/%d", publicationId), userToken, nil, false)

	checkResponseCode(t, http.StatusOK, response.Code)

	var publicationUpdated model.Publication
	if error := json.Unmarshal(response.Body.Bytes(), &publicationUpdated); error != nil {
		t.Errorf("Expected a publication. Got %s", response.Body)
	}

	if publicationUpdated.Content == publicationOld.Content {
		t.Errorf("Expected a Publication is updated with Content: Pub to test - 1... . Got %s", publicationUpdated.Content)
	}
}

// TestUpdatePublicationByIdWithUnauthorizedError
// It guarantees that a publication will not be updated when the given id is valid and the request is not authenticated.
func TestUpdatePublicationByIdWithUnauthorizedError(t *testing.T) {

	var publicationId = 1

	body := []byte(`{
		"title": "Pub Test - 1",
		"content": "Pub to test - 1..."
	}`)

	response := executeRequest("PUT", fmt.Sprintf("/publications/%d", publicationId), "", body, true)

	checkResponseCode(t, http.StatusUnauthorized, response.Code)
}

// TestUpdatePublicationByIdWithBadRequestErrorWhenPublicationIdIsInvalid
// It guarantees that a publication will not be updated when the id entered is invalid and the request is authenticated.
func TestUpdatePublicationByIdWithBadRequestErrorWhenPublicationIdIsInvalid(t *testing.T) {

	var publicationId = "d"

	body := []byte(`{
		"title": "Pub Test - 1",
		"content": "Pub to test - 1..."
	}`)

	response := executeRequest("PUT", fmt.Sprintf("/publications/%s", publicationId), userToken, body, true)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

// TestUpdatePublicationByIdWithBadRequestErrorWhenBodyIsInvalid
// It guarantees that a publication will not be updated when the given id is valid, the request
// is authenticated and the request body is invalid.
func TestUpdatePublicationByIdWithBadRequestErrorWhenBodyIsInvalid(t *testing.T) {

	var publicationId = 1

	body := []byte(``)

	response := executeRequest("PUT", fmt.Sprintf("/publications/%d", publicationId), userToken, body, true)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

// TestUpdatePublicationByIdWithBadRequestErrorWhenFieldIsMandatory
// It guarantees that a publication will not be updated when the given id is valid, the request is authenticated and
// the request body does not have a mandatory field.
func TestUpdatePublicationByIdWithBadRequestErrorWhenFieldIsMandatory(t *testing.T) {

	var publicationId = 1

	body := []byte(`{
		"title": "Pub Test - 1"
	}`)

	response := executeRequest("PUT", fmt.Sprintf("/publications/%d", publicationId), userToken, body, true)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

// TestUpdatePublicationByIdWithForbiddenError
// It guarantees that a publication will not be updated when the id entered is valid, the request
// is authenticated but the publication belongs to another user.
func TestUpdatePublicationByIdWithForbiddenError(t *testing.T) {

	var publicationId = 2

	body := []byte(`{
		"title": "Pub Test - 1",
		"content": "Pub to test - 1..."
	}`)

	response := executeRequest("PUT", fmt.Sprintf("/publications/%d", publicationId), userToken, body, true)

	checkResponseCode(t, http.StatusForbidden, response.Code)
}

// TestUpdatePublicationByIdWithSuccess
// It guarantees that a publication will be deleted when the id entered is valid, the request is
// authenticated but the publication belongs to the user.
func TestDeletePublicationByIdWithSuccess(t *testing.T) {

	var publicationId = 4

	response := executeRequest("DELETE", fmt.Sprintf("/publications/%d", publicationId), userToken, nil, false)

	checkResponseCode(t, http.StatusNoContent, response.Code)
}

// TestDeletePublicationByIdWithUnauthorizedError
// It guarantees that a publication will not be deleted when the given id is valid and the request is not authenticated.
func TestDeletePublicationByIdWithUnauthorizedError(t *testing.T) {

	var publicationId = 1

	response := executeRequest("DELETE", fmt.Sprintf("/publications/%d", publicationId), "", nil, false)

	checkResponseCode(t, http.StatusUnauthorized, response.Code)
}

// TestDeletePublicationByIdWithBadRequestError
// It guarantees that a publication will not be deleted when the id entered is invalid and the request is authenticated.
func TestDeletePublicationByIdWithBadRequestError(t *testing.T) {

	var publicationId = "d"

	response := executeRequest("DELETE", fmt.Sprintf("/publications/%s", publicationId), userToken, nil, false)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

// TestDeletePublicationByIdWithForbiddenError
// Ensures that a post will not be deleted when the post belongs to another user and the request is authenticated.
func TestDeletePublicationByIdWithForbiddenError(t *testing.T) {

	var publicationId = 2

	response := executeRequest("DELETE", fmt.Sprintf("/publications/%d", publicationId), userToken, nil, false)

	checkResponseCode(t, http.StatusForbidden, response.Code)
}

// TestListPublicationsByUserIdWithSuccess
// It guarantees that the list of publications of the informed user will be returned.
func TestListPublicationsByUserIdWithSuccess(t *testing.T) {

	var userId = 1

	response := executeRequest("GET", fmt.Sprintf("/publications/user/%d", userId), userToken, nil, false)

	checkResponseCode(t, http.StatusOK, response.Code)

	var publications []model.Publication
	if error := json.Unmarshal(response.Body.Bytes(), &publications); error != nil {
		t.Errorf("Expected a publication list. Got %s", response.Body)
	}

	if len(publications) < 1 {
		t.Errorf("Expected a Publication list. Got a empty list.")
	}
}

// TestListPublicationsByUserIdWithEmptyListSuccessWhenUserHasNoPublications
// It guarantees that the empty post list will be returned when the user has no publications.
func TestListPublicationsByUserIdWithEmptyListSuccessWhenUserHasNoPublications(t *testing.T) {

	// Get token of new user
	token := GetUserToken("user.test2@gmail.com", "test2")

	var userId = 4

	response := executeRequest("GET", fmt.Sprintf("/publications/user/%d", userId), token, nil, false)

	checkResponseCode(t, http.StatusOK, response.Code)

	var publications []model.Publication
	if error := json.Unmarshal(response.Body.Bytes(), &publications); error != nil {
		t.Errorf("Expected a empty publication list. Got %s", response.Body)
	}

	if len(publications) > 0 {
		t.Errorf("Expected a empty Publication list. Got %q", publications)
	}
}

// TestListPublicationsByUserIdWithUnauthorizedError
// It guarantees that the list of publications will not be returned when the informed user is valid and the request is not authenticated.
func TestListPublicationsByUserIdWithUnauthorizedError(t *testing.T) {

	var userId = 1

	response := executeRequest("GET", fmt.Sprintf("/publications/user/%d", userId), "", nil, false)

	checkResponseCode(t, http.StatusUnauthorized, response.Code)
}

// TestListPublicationsByUserIdWithBadRequestError
// It guarantees that the list of publications will not be returned when the user entered is invalid.
func TestListPublicationsByUserIdWithBadRequestError(t *testing.T) {

	var userId = "d"

	response := executeRequest("GET", fmt.Sprintf("/publications/user/%s", userId), userToken, nil, false)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

// TestLikePublicationByIdWithSuccess
// Ensures post will be liked when post id is valid and request is authenticated.
func TestLikePublicationByIdWithSuccess(t *testing.T) {

	var publicationId = 1

	// Get publication before like.
	response := executeRequest("GET", fmt.Sprintf("/publications/%d", publicationId), userToken, nil, false)

	checkResponseCode(t, http.StatusOK, response.Code)

	var publicationBeforeLike model.Publication
	if error := json.Unmarshal(response.Body.Bytes(), &publicationBeforeLike); error != nil {
		t.Errorf("Expected a publication. Got %s", response.Body)
	}

	// Like publication
	response = executeRequest("POST", fmt.Sprintf("/publications/%d/like", publicationId), userToken, nil, false)

	checkResponseCode(t, http.StatusNoContent, response.Code)

	// Get publication after like.
	response = executeRequest("GET", fmt.Sprintf("/publications/%d", publicationId), userToken, nil, false)

	checkResponseCode(t, http.StatusOK, response.Code)

	var publicationAfterLike model.Publication
	if error := json.Unmarshal(response.Body.Bytes(), &publicationAfterLike); error != nil {
		t.Errorf("Expected a publication. Got %s", response.Body)
	}

	if publicationAfterLike.Likes <= publicationBeforeLike.Likes {
		t.Errorf("Expected a after publication like is greater than publication before like.")
	}
}

// TestLikePublicationByIdWithBadRequestError
// It guarantees that the post will not be liked when the post id is invalid and the request is authenticated.
func TestLikePublicationByIdWithBadRequestError(t *testing.T) {

	var publicationId = "d"

	response := executeRequest("POST", fmt.Sprintf("/publications/%s/like", publicationId), userToken, nil, false)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

// TestLikePublicationByIdWithUnauthorizedError
// It guarantees that the post will not be liked when the post id is valid and the request is not authenticated.
func TestLikePublicationByIdWithUnauthorizedError(t *testing.T) {

	var publicationId = 1

	response := executeRequest("POST", fmt.Sprintf("/publications/%d/like", publicationId), "", nil, false)

	checkResponseCode(t, http.StatusUnauthorized, response.Code)
}

// TestLikePublicationByIdWithSuccess
// Ensures post will be liked when post id is valid and request is authenticated.
func TestUnlikePublicationByIdWithSuccess(t *testing.T) {

	var publicationId = 1

	// Get publication before unlike.
	response := executeRequest("GET", fmt.Sprintf("/publications/%d", publicationId), userToken, nil, false)

	checkResponseCode(t, http.StatusOK, response.Code)

	var publicationBeforeUnlike model.Publication
	if error := json.Unmarshal(response.Body.Bytes(), &publicationBeforeUnlike); error != nil {
		t.Errorf("Expected a publication. Got %s", response.Body)
	}

	// Unlike publication
	response = executeRequest("POST", fmt.Sprintf("/publications/%d/unlike", publicationId), userToken, nil, false)

	checkResponseCode(t, http.StatusNoContent, response.Code)

	// Get publication after unlike.
	response = executeRequest("GET", fmt.Sprintf("/publications/%d", publicationId), userToken, nil, false)

	checkResponseCode(t, http.StatusOK, response.Code)

	var publicationAfterUnlike model.Publication
	if error := json.Unmarshal(response.Body.Bytes(), &publicationAfterUnlike); error != nil {
		t.Errorf("Expected a publication. Got %s", response.Body)
	}

	if publicationAfterUnlike.Likes >= publicationBeforeUnlike.Likes {
		t.Errorf("Expected a after publication unlike is smaller than publication before unlike.")
	}
}

// TestUnlikePublicationByIdWithUnauthorizedError
// It guarantees that the post will not be unliked when the post id is valid and the request is not authenticated.
func TestUnlikePublicationByIdWithUnauthorizedError(t *testing.T) {

	var publicationId = 1

	response := executeRequest("POST", fmt.Sprintf("/publications/%d/unlike", publicationId), "", nil, false)

	checkResponseCode(t, http.StatusUnauthorized, response.Code)
}

// TestUnlikePublicationByIdWithBadRequestError
// It guarantees that the post will not be unliked when the post id is invalid and the request is authenticated.
func TestUnlikePublicationByIdWithBadRequestError(t *testing.T) {

	var publicationId = "d"

	response := executeRequest("POST", fmt.Sprintf("/publications/%s/unlike", publicationId), userToken, nil, false)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
}
