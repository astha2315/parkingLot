package services

import (
	"fmt"
	"net/http"
	"parkingLot/common"
	"parkingLot/dao"
	"parkingLot/models"
	"strconv"
	"time"
)

type parkingLotService struct{}

type ParkingLotServiceIF interface {
	CreateParkingLot(parkingLot *models.ParkingLot) (int, string)
	ParkVehicle(vehicle *models.VehicleJSON) (int, string)

	LeaveParkingBySlotId(slotId int) (int, string)

	LeaveParkingByPlateNo(plateNo string) (int, string)

	GetFreeSpaceForParking(parkingLotId int) (int, string)
}

func ParkingLotService() ParkingLotServiceIF {
	return &parkingLotService{}
}

func (self *parkingLotService) CreateParkingLot(parkingLot *models.ParkingLot) (int, string) {

	slotCapacity := parkingLot.SlotCapacity

	if slotCapacity <= 0 {

		return 0, common.MarshalJson(common.ResponseJson{0, nil, 0, "Capacity of parking lot cannot be less that 1"})

	}

	parkingLotId, err := dao.ParkingLotDao().CreateParkingLot(parkingLot)

	if err != nil {
		return 0, common.MarshalJson(common.ResponseJson{0, nil, 0, err.Error()})
	}

	// now adding the slots

	var slots []*models.Slot

	for i := 0; i < slotCapacity; i++ {

		slot := new(models.Slot)
		slot.ParkingLotId = parkingLotId
		slot.ParkingStatus = 0

		slots = append(slots, slot)
	}

	rowsAffected, err := dao.ParkingLotDao().AddSlots(slots)

	if err != nil {
		return 0, common.MarshalJson(common.ResponseJson{0, nil, 0, err.Error()})
	}

	fmt.Println("Rows Affected are", rowsAffected)

	parkingLotObj := new(models.ParkingLot)
	parkingLotObj.ParkingLotId = parkingLotId
	parkingLotObj.SlotCapacity = slotCapacity

	return 1, common.MarshalJson(common.ResponseJson{http.StatusOK, parkingLotObj, 1, "Successfully created parking lot"})
}

func (self *parkingLotService) ParkVehicle(vehicle *models.VehicleJSON) (int, string) {

	// getting the free slot available

	slotInfo, err := dao.ParkingLotDao().GetFreeSlot(vehicle.ParkingLotId)

	if err != nil {
		return 0, common.MarshalJson(common.ResponseJson{0, nil, 0, err.Error()})
	}

	if slotInfo == nil {
		return 0, common.MarshalJson(common.ResponseJson{0, nil, 0, "No slot found for the parking"})
	}

	// now updating the parking status for the slot

	_, err = dao.ParkingLotDao().UpdateSlotForParking(slotInfo.SlotId, 1)

	if err != nil {
		return 0, common.MarshalJson(common.ResponseJson{0, nil, 0, err.Error()})
	}

	v := &models.Vehicle{
		SlotId:        slotInfo.SlotId,
		ParkingTime:   time.Now().Format("2006-01-02 15:04:05"),
		Color:         vehicle.Color,
		ParkingStatus: 1,
		PlateNumber:   vehicle.PlateNumber,
	}

	vehicleId, err := dao.ParkingLotDao().AddVehicleForParking(v)

	if err != nil {
		return 0, common.MarshalJson(common.ResponseJson{0, nil, 0, err.Error()})
	}

	fmt.Println("Vehicle Id is", vehicleId)

	return 1, common.MarshalJson(common.ResponseJson{http.StatusOK, nil, 1, "Successfully parked the car"})

}

func (self *parkingLotService) LeaveParkingBySlotId(slotId int) (int, string) {

	vehicleInfo, err := dao.ParkingLotDao().GetVehicleBySlotId(slotId)

	if err != nil {
		return 0, common.MarshalJson(common.ResponseJson{0, nil, 0, err.Error()})
	}

	if vehicleInfo == nil {
		return 0, common.MarshalJson(common.ResponseJson{0, nil, 0, "No vehicle found for the slotId"})
	}

	// now freeing the slot for the slotId

	_, err = dao.ParkingLotDao().UpdateSlotForParking(slotId, 0)
	if err != nil {
		return 0, common.MarshalJson(common.ResponseJson{0, nil, 0, err.Error()})
	}

	// now updating the vehicle info and changing the parking status

	vehicleId := vehicleInfo.VehicleId

	_, err = dao.ParkingLotDao().UpdateVehicleForParking(vehicleId, 0)

	if err != nil {
		return 0, common.MarshalJson(common.ResponseJson{0, nil, 0, err.Error()})
	}

	return 1, common.MarshalJson(common.ResponseJson{http.StatusOK, nil, 0, "Successfully freed the slot"})

}

func (self *parkingLotService) LeaveParkingByPlateNo(plateNo string) (int, string) {

	vehicleInfo, err := dao.ParkingLotDao().GetVehicleByPlateNo(plateNo)

	if err != nil {
		return 0, common.MarshalJson(common.ResponseJson{0, nil, 0, err.Error()})
	}

	if vehicleInfo == nil {
		return 0, common.MarshalJson(common.ResponseJson{0, nil, 0, "No vehicle found for the plate number"})
	}

	vehicleId := vehicleInfo.VehicleId
	slotId := vehicleInfo.SlotId

	// now updating the slot and freeing it

	_, err = dao.ParkingLotDao().UpdateSlotForParking(slotId, 0)
	if err != nil {
		return 0, common.MarshalJson(common.ResponseJson{0, nil, 0, err.Error()})
	}

	// now updating the vehicle info and changing the parking status

	_, err = dao.ParkingLotDao().UpdateVehicleForParking(vehicleId, 0)

	if err != nil {
		return 0, common.MarshalJson(common.ResponseJson{0, nil, 0, err.Error()})
	}

	return 1, common.MarshalJson(common.ResponseJson{http.StatusOK, nil, 0, "Successfully freed the slot for the given plate number"})

}

func (self *parkingLotService) GetFreeSpaceForParking(parkingLotId int) (int, string) {

	slotInfo, err := dao.ParkingLotDao().GetSlotsForParkingId(parkingLotId)

	if err != nil {
		return 0, common.MarshalJson(common.ResponseJson{0, nil, 0, err.Error()})
	}

	if slotInfo == nil || (slotInfo != nil && len(slotInfo) == 0) {
		return 0, common.MarshalJson(common.ResponseJson{0, nil, 0, "No space available for parking"})
	}

	freeSpaceAvailable := len(slotInfo)

	result := strconv.Itoa(freeSpaceAvailable) + " Slots available for parking"

	return 1, common.MarshalJson(common.ResponseJson{http.StatusOK, nil, 0, result})
}
