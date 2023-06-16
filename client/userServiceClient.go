package client

import (
	"encoding/json"
	"github.com/windbnb/reservation-service/model"
	"github.com/windbnb/reservation-service/util"
	"net/http"
)

func AuthorizeHost(tokenString string) (model.UserResponseDTO, error) {
	userUrl, _ := util.GetUserServicePathRoundRobin()
	url := userUrl.Next().Host + "/api/users/authorize/host"
	req, err := http.NewRequest("POST", url, nil)
	req.Header.Add("Authorization", tokenString)
	client := &http.Client{}
	response, err := client.Do(req)

	if err != nil {
		return model.UserResponseDTO{}, err
	}

	var userResponse model.UserResponseDTO
	json.NewDecoder(response.Body).Decode(&userResponse)
	return userResponse, nil
}

func AuthorizeGueest(tokenString string) (model.UserResponseDTO, error) {
	userUrl, _ := util.GetUserServicePathRoundRobin()
	url := userUrl.Next().Host + "/api/users/authorize/guest"
	req, err := http.NewRequest("POST", url, nil)
	req.Header.Add("Authorization", tokenString)
	client := &http.Client{}
	response, err := client.Do(req)

	if err != nil {
		return model.UserResponseDTO{}, err
	}

	var userResponse model.UserResponseDTO
	json.NewDecoder(response.Body).Decode(&userResponse)
	return userResponse, nil
}
