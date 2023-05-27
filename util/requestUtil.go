package util

import (
	"net/url"

	roundrobin "github.com/hlts2/round-robin"
)

var BaseUserServicePathRoundRobin, _ = roundrobin.New(
	&url.URL{Host: "http://localhost:8081"},
)

var BaseAccommodationServicePathRoundRobin, _ = roundrobin.New(
	&url.URL{Host: "http://localhost:8082"},
)
