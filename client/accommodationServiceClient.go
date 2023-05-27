package client

import (
	"bytes"
	"encoding/json"
	"github.com/windbnb/reservation-service/model"
	"github.com/windbnb/reservation-service/util"
	"net/http"
	"strconv"
)

func GetAccommodation(accommodationID uint) (model.AccommodationInfo, error) {
	url := util.BaseAccommodationServicePathRoundRobin.Next().Host + "/api/accomodation/" + strconv.FormatUint(uint64(accommodationID), 10)
	req, err := http.NewRequest("GET", url, nil)
	//req.Header.Add("Authorization", tokenString)
	client := &http.Client{}
	response, err := client.Do(req)

	if err != nil {
		return model.AccommodationInfo{}, err
	}

	var accommodationInfo model.AccommodationInfo
	json.NewDecoder(response.Body).Decode(&accommodationInfo)
	return accommodationInfo, nil
}

func CreateReservedTerm(reservationRequest model.ReservationRequest) (uint, error) {
	url := util.BaseAccommodationServicePathRoundRobin.Next().Host + "/api/accomodation/reservedTerm"

	reservedTerm := model.ReservedTermRequest{
		StartDate:      reservationRequest.StartDate,
		EndDate:        reservationRequest.EndDate,
		AccomodationID: reservationRequest.AccommodationID}

	marshalled, err := json.Marshal(reservedTerm)
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(marshalled))
	//req.Header.Add("Authorization", tokenString)
	client := &http.Client{}
	response, err := client.Do(req)

	if err != nil {
		return 0, err
	}

	var reservedTermResponse model.ReservedTermResponse
	json.NewDecoder(response.Body).Decode(&reservedTermResponse)
	return reservedTermResponse.Id, nil
}

func DeleteReservedTerm(reservedTermId uint) {
	url := util.BaseAccommodationServicePathRoundRobin.Next().Host + "/api/accomodation/reservedTerm" + strconv.FormatUint(uint64(reservedTermId), 10)
	req, _ := http.NewRequest("DELETE", url, nil)
	//req.Header.Add("Authorization", tokenString)
	client := &http.Client{}
	_, _ = client.Do(req)
}
