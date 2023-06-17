package util

import (
	"net/url"
	"os"

	roundrobin "github.com/hlts2/round-robin"
)

func GetUserServicePathRoundRobin() (roundrobin.RoundRobin, error) {
	userServicePath, userServicePathFound := os.LookupEnv("USER_SERVICE_PATH")
	if !userServicePathFound {
		userServicePath = "http://localhost:8081"
	}

	return roundrobin.New(
		&url.URL{Host: userServicePath},
	)
}

func GetAccommodationServicePathRoundRobin() (roundrobin.RoundRobin, error) {
	accommodationServicePath, accommodationServicePathFound := os.LookupEnv("ACCOMMODATION_SERVICE_PATH")
	if !accommodationServicePathFound {
		accommodationServicePath = "http://localhost:8082"
	}

	return roundrobin.New(
		&url.URL{Host: accommodationServicePath},
	)
}
