package routing

import (
	"encoding/json"
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	distanceURL = "https://maps.googleapis.com/maps/api/distancematrix/json?key=%s&origins=%f,%f&destinations=%f,%f&departure_time=now"
	ok          = "OK"
)

var (
	ErrorGoogleMapDistance = status.Error(codes.Unavailable, "Failed to get distance from google map")
)

//Field is distance/duration value
type Field struct {
	Value float64 `json:"value,omitempty"`
}

//Element is information about each origin-destination pairing distance/duration
type Element struct {
	Distance          Field  `json:"distance,omitempty"`
	Duration          Field  `json:"duration,omitempty"`
	DurationInTraffic Field  `json:"duration_in_traffic,omitempty"`
	Status            string `json:"status,omitempty"`
}

//Rows is array of element
type Rows struct {
	Elements []Element `json:"elements,omitempty"`
}

//Response is google map distance matrix response
type Response struct {
	Rows []Rows `json:"rows,omitempty"`
}

type locator struct {
	gmapKey string
}

//NewLocator return new redis cache
func NewLocator(key string) RoutingData {
	return &locator{gmapKey: key}
}

func (r *locator) GetDistanceTime(loc1, loc2 Point) (float64, float64, error) {
	resp, err := http.Get(fmt.Sprintf(distanceURL, r.gmapKey, loc1.Latitude, loc1.Longitude, loc2.Latitude, loc2.Longitude))
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()
	var response Response
	json.NewDecoder(resp.Body).Decode(&response)

	rows := response.Rows
	if len(rows) > 0 {
		elements := rows[0].Elements
		if len(elements) > 0 {
			element := elements[0]
			if element.Status == ok {
				if element.DurationInTraffic.Value > 0 {
					return element.Distance.Value, element.DurationInTraffic.Value, nil
				}
				return element.Distance.Value, element.Duration.Value, nil
			}
		}
	}

	return 0, 0, ErrorGoogleMapDistance
}
