package main

import (
	"log"

	dt "git.bluebird.id/bluebird/tracker/routing"
)

func main() {

	loc1 := dt.Point{Latitude: -6.260971, Longitude: 106.829552}
	loc2 := dt.Point{Latitude: -6.273751, Longitude: 106.831823}

	// Testing OSRM
	route := dt.NewOSRMRoutingData(nil)
	distance, time, err := route.GetDistanceTime(loc1, loc2)
	if err != nil {
		log.Println("Error : ", err)
		return
	}
	// Testing OSRM

	// Testing Graphhoper
	// route := dt.NewGrapHhoperRoutingData("0126024f-2098-485a-bd54-89ff74850d34", nil)
	// distance, time, err := route.GetDistanceTime(loc1, loc2)
	// if err != nil {
	// 	log.Println("Error : ", err)
	// 	return
	// }
	// Testing Graphhoper

	// Testing Google Map
	// route := dt.NewLocator("AIzaSyD9tm3UVfxRWeaOy_MQ7tsCj1fVCLfG8Bo")
	// distance, time, err := route.GetDistanceTime(loc1, loc2)
	// if err != nil {
	// 	log.Println("Error : ", err)
	// 	return
	// }
	// Testing Google Map

	log.Printf("Distance : %f meters\n", distance)
	log.Printf("Time : %f seconds\n", time)
}
