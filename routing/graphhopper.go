package routing

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// CalcPoints is points for the route should be calculated at all printing out only distance and time
type CalcPoints bool

// Instructions is instruction direction way like to "keep this way" or "500m will turn right into another way"
type Instructions bool

// Vehicle is vehicle like car, scooter, foot, bike, etc
type Vehicle string

// PointsEncoded every point will be encoded or not encoded to string or lat,long
type PointsEncoded bool

// Elevation is will include polyline or geoJSON
type Elevation bool

const (
	// CarVehicle is Car mode
	CarVehicle Vehicle = "car"
	// VanVehicle is like Van or Small truck like a Mercedes Sprinter, Ford Transit or Iveco Daily which have height=2.7m, width=2+0.4m, length=5.5m, weight=2080+1400 kg
	VanVehicle Vehicle = "small_truck"
	// TruckVehicle is like Mercedes-Benz Actros or have height=3.7m, width=2.6+0.5m, length=12m, weight=13000 + 13000 kg, hgv=yes, 3 Axes
	TruckVehicle Vehicle = "truck"

	// ActiveInstruction is always give instructions
	ActiveInstruction Instructions = true
	// NonActiveInstruction is not give instructions
	NonActiveInstruction Instructions = false

	// ActiveCalcPoints will calculated all route and give total distance and time
	ActiveCalcPoints CalcPoints = true
	// NonActiveCalcPoints not calculated all route
	NonActiveCalcPoints CalcPoints = false

	// ActivePointsEncoded is the coordinates in point and snapped_waypoints are returned as array using the order [lon,lat,elevation] for every point.
	ActivePointsEncoded PointsEncoded = true
	// NonActivePointsEncoded is  the coordinates will be encoded as string leading to less bandwith usage. You'll need a special handling for the decoding of this string on the client-side.
	NonActivePointsEncoded PointsEncoded = false

	// ActiveElevation is will return polyline or geoJSON. BUT, pointsEncode MUST Non Active or FALSE
	ActiveElevation Elevation = true
	// NonActiveElevation is will not return polyline or geoJSON
	NonActiveElevation Elevation = false

	// GRAPHHOPER_ROUTE is API v.1to get route service
	GRAPHHOPER_ROUTE = `https://graphhopper.com/api/1/route?point=%f,%f&point=%f,%f`
	// GRAPHHOPER_ROUTE_OPTION is API v.1 to option for GRAPHHOPER_ROUTE
	GRAPHHOPER_ROUTE_OPTION = `&locale=%s&instructions=%v&vehicle=%s&elevation=%v&points_encoded=%v&calc_points=%v&key=%s`
)

type OptionGrapHhoper struct {
	Locale        string
	Instructions  bool
	Vehicle       string
	Elevation     bool
	PointsEncoded bool
	CalcPoints    bool
}

type graphhoper struct {
	link string
}

type graphhoperRoute struct {
	Hints struct {
		VisitedNodesAverage string `json:"visited_nodes.average"`
		VisitedNodesSum     string `json:"visited_nodes.sum"`
	} `json:"hints"`
	Info struct {
		Copyrights []string `json:"copyrights"`
		Took       int      `json:"took"`
	} `json:"info"`
	Paths []struct {
		Distance      float64       `json:"distance"`
		Weight        float64       `json:"weight"`
		Time          int           `json:"time"`
		Transfers     int           `json:"transfers"`
		PointsEncoded bool          `json:"points_encoded"`
		Bbox          []float64     `json:"bbox"`
		Points        string        `json:"points"`
		Legs          []interface{} `json:"legs"`
		Details       struct {
		} `json:"details"`
		Ascend           float64 `json:"ascend"`
		Descend          float64 `json:"descend"`
		SnappedWaypoints string  `json:"snapped_waypoints"`
	} `json:"paths"`
}

// NewGrapHhoperRoutingData is OSRM to get routing data
func NewGrapHhoperRoutingData(key string, option *OptionGrapHhoper) RoutingData {
	if option == nil {
		option = defaultOptionGrapHhoper()
	}
	ops := fmt.Sprintf(GRAPHHOPER_ROUTE_OPTION, option.Locale, option.Instructions, option.Vehicle, option.Elevation, option.PointsEncoded, option.CalcPoints, key)
	link := GRAPHHOPER_ROUTE + ops
	return &graphhoper{link: link}
}

// GetDistanceTime is return distance in meters and time in seconds
func (gp *graphhoper) GetDistanceTime(loc1, loc2 Point) (float64, float64, error) {
	url := fmt.Sprintf(gp.link, loc1.Latitude, loc1.Longitude, loc2.Latitude, loc2.Longitude)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("[ERROR GetDistanceTime Graphhoper] : ", err)
		return 0, 0, err
	}
	defer resp.Body.Close()

	readLoc, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("[ERROR GetDistanceTime Graphhoper] Can't Readall req body : ", err)
		return 0, 0, err
	}

	var route graphhoperRoute
	err = json.Unmarshal(readLoc, &route)
	if err != nil {
		log.Println("[ERROR GetDistanceTime Graphhoper] Unable to unmarshal response : ", err)
		return 0, 0, err
	}
	rt := route.Paths[0]

	// convert from milliseconds to seconds
	time := float64(rt.Time / 1000)

	return rt.Distance, time, nil
}

// defaultOptionGrapHhoper is default option for graphhoper v.1
func defaultOptionGrapHhoper() *OptionGrapHhoper {
	return &OptionGrapHhoper{
		Locale:        "en",
		Instructions:  bool(NonActiveInstruction),
		Vehicle:       string(CarVehicle),
		Elevation:     bool(NonActiveElevation),
		PointsEncoded: bool(ActivePointsEncoded),
		CalcPoints:    bool(ActiveCalcPoints),
	}
}
