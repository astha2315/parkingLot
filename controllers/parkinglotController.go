package controllers

import (
	"fmt"
	"net/http"
	"parkingLot/models"
	"parkingLot/services"
	"strconv"

	"github.com/labstack/echo"
)

func CreateParkingLot(c echo.Context) error {

	fmt.Println("In create parking lot api")

	parkingLotJSON := new(models.ParkingLot)
	err := c.Bind(parkingLotJSON)
	if err != nil {
		return err
	}

	status, response := services.ParkingLotService().CreateParkingLot(parkingLotJSON)

	fmt.Println(status)
	return c.String(http.StatusOK, response)

}

func ParkVehicle(c echo.Context) error {

	fmt.Println("In park vehicle api")

	vehicle := new(models.VehicleJSON)
	err := c.Bind(vehicle)
	if err != nil {
		return err
	}

	status, response := services.ParkingLotService().ParkVehicle(vehicle)

	fmt.Println(status)
	return c.String(http.StatusOK, response)

}

func LeaveParkingBySlotId(c echo.Context) error {

	fmt.Println("In leave parking by slot number api")

	slotId, err := strconv.Atoi(c.Param("slotId"))

	if err != nil {
		return err
	}

	status, response := services.ParkingLotService().LeaveParkingBySlotId(slotId)

	fmt.Println(status)
	return c.String(http.StatusOK, response)

}

func LeaveParkingByPlateNo(c echo.Context) error {

	fmt.Println("In leave parking by plateNo api")

	plateNo := c.Param("plateNo")

	status, response := services.ParkingLotService().LeaveParkingByPlateNo(plateNo)

	fmt.Println(status)
	return c.String(http.StatusOK, response)

}

func GetFreeSpaceForParking(c echo.Context) error {

	fmt.Println("In get free space for parking lot api")

	parkingLotId, err := strconv.Atoi(c.Param("parkingLotId"))
	if err != nil {
		return err
	}

	status, response := services.ParkingLotService().GetFreeSpaceForParking(parkingLotId)

	fmt.Println(status)
	return c.String(http.StatusOK, response)

}
