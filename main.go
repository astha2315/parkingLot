package main

import (
	"fmt"
	"parkingLot/db"
	route "parkingLot/routes"

	"github.com/labstack/echo"
)

func main() {
	fmt.Println("Hello")

	_, err := db.DBConnect()
	if err != nil {
		panic(err)
	}

	// http.HandleFunc("/parkinglot/list", enf(&enforcer, controllers.GetMovies))

	// log.Fatal(http.ListenAndServe(":8080", nil))

	e := echo.New()

	route.ParkingLotRouteService(e)

	e.Logger.Fatal(e.Start(":8080"))

}
