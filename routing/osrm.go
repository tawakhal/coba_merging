package routing

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Alternatives is altenative routing, true, false, and number of maximal alternative.
type Alternatives string

// Steps is returned steps from routing
type Steps bool

// Geometries is returned route geometry format (influences overview and per step)
type Geometries string

// Overview is add overview geometry either full, simplified according to highest zoom level it could be display on, or not at all
type Overview string

// ContinueStraight is forces the route to keep going straight at waypoints constraining uturns there even if it would be faster. Default value depends on the profile.
type ContinueStraight string

const (
	// ActiveAlternative is always search alternative routing from location 1 to location 2
	ActiveAlternative Alternatives = "true"
	// NonActiveAlternative is not have alternative routing
	NonActiveAlternative Alternatives = "false"

	// ActiveStep is returned steps from the routing
	ActiveStep Steps = true
	// NonActiveStep is not returned steps from the routing
	NonActiveStep Steps = false

	// PolyLineGeometry is using polyline geometry format
	PolyLineGeometry Geometries = "polyline"
	// PolyLine6Geometry is using polyline6 geometry format
	PolyLine6Geometry Geometries = "polyline6"
	// GeoJSONGeometry is using GeoJson geometry format
	GeoJSONGeometry Geometries = "geojson"

	// SimplifiedOverview is simplified according to the highest zoom level it can still be displayed on full
	SimplifiedOverview Overview = "simplified"
	// fullOverview is not simplified.
	fullOverview Overview = "full"
	// NonActiveOverview is not added.
	NonActiveOverview Overview = "false"

	// DefaultContinueStraight is forces way point with default routing
	DefaultContinueStraight ContinueStraight = "default"
	// ActiveContinueStraight is always forces way point when continue straight
	ActiveContinueStraight ContinueStraight = "true"
	// NonActiveContinueStraight is not forces way point
	NonActiveContinueStraight ContinueStraight = "false"

	// OSRM_ROUTE is API v.5.10.0 to get route service
	OSRM_ROUTE = `http://router.project-osrm.org/route/v1/driving/%f,%f;%f,%f`
	// OSRM_ROUTE_OPTION is API v.5.10.0 to option for OSRM_ROUTE
	OSRM_ROUTE_OPTION = `?alternatives=%s&steps=%v&geometries=%s&overview=%s&continue_straight=%s`
)

type osrm struct {
	link string
}

// OptionOSRM is option for OSRM
type OptionOSRM struct {
	alternatives     string
	steps            bool
	geometries       string
	overview         string
	continueStraight string
}

type osrmRoute struct {
	Routes []struct {
		Geometry string `json:"geometry"`
		Legs     []struct {
			Summary  string        `json:"summary"`
			Weight   float64       `json:"weight"`
			Duration float64       `json:"duration"`
			Steps    []interface{} `json:"steps"`
			Distance float64       `json:"distance"`
		} `json:"legs"`
		WeightName string  `json:"weight_name"`
		Weight     float64 `json:"weight"`
		Duration   float64 `json:"duration"`
		Distance   float64 `json:"distance"`
	} `json:"routes"`
	Waypoints []struct {
		Hint     string    `json:"hint"`
		Distance float64   `json:"distance"`
		Name     string    `json:"name"`
		Location []float64 `json:"location"`
	} `json:"waypoints"`
	Code string `json:"code"`
}

// NewOSRMRoutingData is OSRM to get routing data
func NewOSRMRoutingData(opt *OptionOSRM) RoutingData {
	if opt == nil {
		opt = defaultOptionOSRM()
	}
	ops := fmt.Sprintf(OSRM_ROUTE_OPTION, opt.alternatives, opt.steps, opt.geometries, opt.overview, opt.continueStraight)
	link := OSRM_ROUTE + ops
	return &osrm{link: link}
}

// GetDistanceTime is return distance in meter and timeDuration in seconds
func (om *osrm) GetDistanceTime(loc1, loc2 Point) (float64, float64, error) {
	url := fmt.Sprintf(om.link, loc1.Longitude, loc1.Latitude, loc2.Longitude, loc2.Latitude)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("[ERROR GetDistanceTime OSRM] : ", err)
		return 0, 0, err
	}
	defer resp.Body.Close()

	readLoc, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("[ERROR GetDistanceTime OSRM] Can't Readall req body : ", err)
		return 0, 0, err
	}

	var route osrmRoute
	err = json.Unmarshal(readLoc, &route)
	if err != nil {
		log.Println("[ERROR GetDistanceTime OSRM] Unable to unmarshal response : ", err)
		return 0, 0, err
	}
	rt := route.Routes[0]

	return rt.Distance, rt.Duration, nil
}

// defaultOptionOSRM is default option for OSRM v.5.10.0
func defaultOptionOSRM() *OptionOSRM {
	return &OptionOSRM{
		alternatives:     string(NonActiveAlternative),
		steps:            bool(NonActiveStep),
		geometries:       string(PolyLineGeometry),
		overview:         string(SimplifiedOverview),
		continueStraight: string(DefaultContinueStraight),
	}
}
