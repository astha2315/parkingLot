package route

import (
	"parkingLot/controllers"

	"github.com/labstack/echo"
)

func ParkingLotRouteService(e *echo.Echo) {
	e.POST("/parkinglot/create", controllers.CreateParkingLot)

	e.POST("/parkinglot/park", controllers.ParkVehicle)
	e.GET("/parkinglot/leave/slot/:slotId", controllers.LeaveParkingBySlotId)
	e.GET("/parkinglot/leave/plateNo/:plateNo", controllers.LeaveParkingByPlateNo)
	e.GET("/parkinglot/freespace/:parkingLotId", controllers.GetFreeSpaceForParking)

}
