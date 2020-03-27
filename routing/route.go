package routing

type Point struct {
	Latitude  float64
	Longitude float64
}

type RoutingData interface {
	GetDistanceTime(loc1, loc2 Point) (float64, float64, error)
}
