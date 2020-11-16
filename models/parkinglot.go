package models

type ParkingLot struct {
	ParkingLotId int `db:"parking_lot_id" json:"parkingLotId"`
	SlotCapacity int `db:"slot_capacity" json:"slotCapacity"`
	Slots        []Slot
}

type Slot struct {
	SlotId        int `db:"slot_id" json:"slotId"`
	ParkingLotId  int `db:"parking_lot_id" json:"parkingLotId"`
	ParkingStatus int `db:"parking_status" json:"parkingStatus"`
}

type Vehicle struct {
	VehicleId     int    `db:"vehicle_id" json:"vehicleId"`
	SlotId        int    `db:"slot_id" json:"slotId"`
	PlateNumber   string `db:"plate_number" json:"plateNumber"`
	Color         string `db:"color" json:"color"`
	ParkingStatus int    `db:"parking_status" json:"parkingStatus"`
	ParkingTime   string `db:"parking_time" json:"parkingTime"`
}

type VehicleJSON struct {
	ParkingLotId  int    `db:"parking_lot_id" json:"parkingLotId"`
	VehicleId     int    `db:"vehicle_id" json:"vehicleId"`
	SlotId        int    `db:"slot_id" json:"slotId"`
	PlateNumber   string `db:"plate_number" json:"plateNumber"`
	Color         string `db:"color" json:"color"`
	ParkingStatus int    `db:"parking_status" json:"parkingStatus"`
	ParkingTime   string `db:"parking_time" json:"parkingTime"`
}
