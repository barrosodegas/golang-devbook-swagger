package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"webapp/src/config"
	"webapp/src/model"
	"webapp/src/request"
)

func LoadUserDataByUserId(user *model.User, r *http.Request) error {

	channelFollowers := make(chan []model.User)
	channelFollowing := make(chan []model.User)
	channelPublications := make(chan []model.Publication)

	go findFollowersByUserId(channelFollowers, user.ID, r)
	go findFollowingByUserId(channelFollowing, user.ID, r)
	go findPublicationsByUserId(channelPublications, user.ID, r)

	for i := 0; i < 3; i++ {
		select {
		case followers := <-channelFollowers:
			if followers == nil {
				return errors.New("Fail to find followers!")
			}
			user.Followers = followers
		case following := <-channelFollowing:
			if following == nil {
				return errors.New("Fail to find following!")
			}
			user.Following = following
		case publications := <-channelPublications:
			if publications == nil {
				return errors.New("Fail to find publications!")
			}
			user.Publications = publications
		}
	}

	return nil
}

func findFollowersByUserId(channel chan<- []model.User, userId uint64, r *http.Request) {

	userUrl := fmt.Sprintf("%s/users/%d/followers", config.UrlApi, userId)
	response, error := request.ExecuteRequestWithAuthentication(r, http.MethodGet, userUrl, nil)
	if error != nil {
		fmt.Printf("Error to find followers of user: %d with error: %s", userId, error.Error())
		channel <- nil
		return
	}

	defer response.Body.Close()

	if response.StatusCode >= 400 {
		fmt.Printf("Error to find followers of user: %d with error: %s", userId, error.Error())
		channel <- nil
		return
	}

	var followers []model.User

	if error = json.NewDecoder(response.Body).Decode(&followers); error != nil {
		fmt.Printf("Error to head followers of user: %d with error: %s", userId, error.Error())
		channel <- nil
		return
	}

	if followers == nil {
		channel <- []model.User{}
	} else {
		channel <- followers
	}
}

func findFollowingByUserId(channel chan<- []model.User, userId uint64, r *http.Request) {

	userUrl := fmt.Sprintf("%s/users/%d/list-followed", config.UrlApi, userId)
	response, error := request.ExecuteRequestWithAuthentication(r, http.MethodGet, userUrl, nil)
	if error != nil {
		fmt.Printf("Error to find following of user: %d with error: %s", userId, error.Error())
		channel <- nil
		return
	}

	defer response.Body.Close()

	if response.StatusCode >= 400 {
		fmt.Printf("Error to find following of user: %d with error: %s", userId, error.Error())
		channel <- nil
		return
	}

	var following []model.User

	if error = json.NewDecoder(response.Body).Decode(&following); error != nil {
		fmt.Printf("Error to head following of user: %d with error: %s", userId, error.Error())
		channel <- nil
		return
	}

	if following == nil {
		channel <- []model.User{}
	} else {
		channel <- following
	}
}

func findPublicationsByUserId(channel chan<- []model.Publication, userId uint64, r *http.Request) {

	userUrl := fmt.Sprintf("%s/publications/user/%d", config.UrlApi, userId)
	response, error := request.ExecuteRequestWithAuthentication(r, http.MethodGet, userUrl, nil)
	if error != nil {
		fmt.Printf("Error to find publications of user: %d with error: %s", userId, error.Error())
		channel <- nil
		return
	}

	defer response.Body.Close()

	if response.StatusCode >= 400 {
		fmt.Printf("Error to find publications of user: %d with error: %s", userId, error.Error())
		channel <- nil
		return
	}

	var publications []model.Publication

	if error = json.NewDecoder(response.Body).Decode(&publications); error != nil {
		fmt.Printf("Error to head publications of user: %d with error: %s", userId, error.Error())
		channel <- nil
		return
	}

	if publications == nil {
		channel <- []model.Publication{}
	} else {
		channel <- publications
	}
}
