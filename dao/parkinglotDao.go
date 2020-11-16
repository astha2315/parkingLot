package dao

import (
	"fmt"
	"parkingLot/db"
	"parkingLot/models"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

type parkingLotDao struct{}

type ParkingLotDaoIF interface {
	CreateParkingLot(parkingLot *models.ParkingLot) (int, error)
	AddSlots(slots []*models.Slot) (int, error)

	GetFreeSlot(parkingLotId int) (*models.Slot, error)

	UpdateSlotForParking(slotId int, parkingStatus int) (int, error)

	AddVehicleForParking(vehicle *models.Vehicle) (int, error)

	GetVehicleBySlotId(slotId int) (*models.Vehicle, error)

	UpdateVehicleForParking(vehicleId int, parkingStatus int) (int, error)

	GetVehicleByPlateNo(plateNumber string) (*models.Vehicle, error)

	GetSlotsForParkingId(parkingLotId int) ([]*models.Slot, error)
}

func ParkingLotDao() ParkingLotDaoIF {
	return &parkingLotDao{}
}

func (self *parkingLotDao) CreateParkingLot(parkingLot *models.ParkingLot) (int, error) {

	db, errs := db.DBConnect()

	if errs != nil {
		fmt.Printf("The error is", errs)
		return 0, errs
	}

	sqlStatement := `INSERT INTO parking_lot
	(slot_capacity) 
	VALUES ($1)
	 RETURNING parking_lot_id`

	parkingLotId := 0
	err := db.QueryRow(sqlStatement, parkingLot.SlotCapacity).Scan(&parkingLotId)

	if err != nil {
		return 0, err
	}

	defer db.Close()
	return parkingLotId, err
}

func (self *parkingLotDao) AddSlots(slots []*models.Slot) (int, error) {

	db, errs := db.SqlxConnect()
	if errs != nil {
		fmt.Printf("The error is", errs)
		return 0, errs
	}

	vals := []interface{}{}
	for _, row := range slots {
		vals = append(vals, row.ParkingLotId, row.ParkingStatus)
	}

	sqlStr := `INSERT INTO slot (parking_lot_id,parking_status) VALUES %s`
	sqlStr = ReplaceSQL(sqlStr, "(?,?)", len(slots))

	fmt.Println(sqlStr)
	stmt, statError := db.Prepare(sqlStr)
	if statError != nil {
		return 0, statError
	}
	fmt.Println("****")
	res, excErr := stmt.Exec(vals...)
	// fmt.Println("%v", vals)
	if excErr != nil {
		fmt.Println(res, excErr.Error())
		return 0, excErr
	}

	defer db.Close()
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return 0, err
	}
	return int(rowsAffected), errs
}

func ReplaceSQL(stmt, pattern string, len int) string {
	pattern += ","
	stmt = fmt.Sprintf(stmt, strings.Repeat(pattern, len))
	n := 0
	for strings.IndexByte(stmt, '?') != -1 {
		n++
		param := "$" + strconv.Itoa(n)
		stmt = strings.Replace(stmt, "?", param, 1)
	}
	return strings.TrimSuffix(stmt, ",")
}

func (self *parkingLotDao) GetFreeSlot(parkingLotId int) (*models.Slot, error) {

	slots := []*models.Slot{}
	// db, errs := db.DBConnect()
	db, errs := db.SqlxConnect()
	if errs != nil {
		fmt.Printf("The error is", errs)
		return nil, errs
	}

	fmt.Println("the error is")
	fmt.Println(errs)
	sqlStatement := `SELECT
	                 slot_id,parking_lot_id,parking_status	
				     FROM
				     slot 
				     WHERE parking_status=0 and parking_lot_id=$1 order by slot_id asc limit 1
				`
	err := db.Select(&slots, sqlStatement, parkingLotId)

	if err != nil {
		fmt.Printf("The error is", err)
		return nil, err
	}

	// fmt.Printf("materials are", materials)
	defer db.Close()

	if slots != nil && len(slots) == 1 {
		return slots[0], err
	}
	return nil, err
}

func (self *parkingLotDao) UpdateSlotForParking(slotId int, parkingStatus int) (int, error) {

	db, errs := db.DBConnect()
	var slotInfo []models.Slot

	if errs != nil {
		fmt.Printf("The error is", errs)
		return 0, errs
	}

	sqlStatement := `	UPDATE slot set parking_status=$1 where slot_id=$2;
   `

	query, args, err := sqlx.In(sqlStatement, parkingStatus, slotId)

	if err != nil {
		return 0, err
	}

	query = db.Rebind(query)
	err = db.Select(&slotInfo, query, args...)
	if err != nil {
		return 0, err
	}

	defer db.Close()
	return 0, err
}

func (self *parkingLotDao) AddVehicleForParking(vehicle *models.Vehicle) (int, error) {

	db, errs := db.DBConnect()

	if errs != nil {
		fmt.Printf("The error is", errs)
		return 0, errs
	}

	sqlStatement := `INSERT INTO vehicle
	(slot_id,plate_number,color,parking_status,parking_time) 
	VALUES ($1,$2,$3,$4,$5)
	 RETURNING vehicle_id`

	parkingLotId := 0
	err := db.QueryRow(sqlStatement, vehicle.SlotId, vehicle.PlateNumber, vehicle.Color, vehicle.ParkingStatus, vehicle.ParkingTime).Scan(&parkingLotId)

	if err != nil {
		return 0, err
	}

	defer db.Close()
	return parkingLotId, err

}

func (self *parkingLotDao) GetVehicleBySlotId(slotId int) (*models.Vehicle, error) {

	var vehicles []*models.Vehicle
	// db, errs := db.DBConnect()
	db, errs := db.SqlxConnect()
	fmt.Println("the error is")
	fmt.Println(errs)
	sqlStatement := `SELECT
	                 slot_id,vehicle_id	
				     FROM
				     vehicle 
				     WHERE parking_status=1 and slot_id=$1
				`
	err := db.Select(&vehicles, sqlStatement, slotId)

	if err != nil {
		fmt.Printf("The error is", err)
		return nil, err
	}

	// fmt.Printf("materials are", materials)
	defer db.Close()

	if vehicles != nil && len(vehicles) == 1 {
		return vehicles[0], err
	}
	return nil, err

}

func (self *parkingLotDao) UpdateVehicleForParking(vehicleId int, parkingStatus int) (int, error) {

	db, errs := db.DBConnect()
	var vehicleInfo []models.Vehicle

	if errs != nil {
		fmt.Printf("The error is", errs)
		return 0, errs
	}

	sqlStatement := `	UPDATE vehicle set parking_status=$1 where vehicle_id=$2;
   `

	query, args, err := sqlx.In(sqlStatement, parkingStatus, vehicleId)

	if err != nil {
		return 0, err
	}

	query = db.Rebind(query)
	err = db.Select(&vehicleInfo, query, args...)
	if err != nil {
		return 0, err
	}

	defer db.Close()
	return 0, err
}

func (self *parkingLotDao) GetVehicleByPlateNo(plateNumber string) (*models.Vehicle, error) {

	var vehicles []*models.Vehicle
	// db, errs := db.DBConnect()
	db, errs := db.SqlxConnect()
	fmt.Println("the error is")
	fmt.Println(errs)
	sqlStatement := `SELECT
	                 slot_id,vehicle_id	
				     FROM
				     vehicle 
				     WHERE parking_status=1 and plate_number=$1
				`
	err := db.Select(&vehicles, sqlStatement, plateNumber)

	if err != nil {
		fmt.Printf("The error is", err)
		return nil, err
	}

	// fmt.Printf("materials are", materials)
	defer db.Close()

	if vehicles != nil && len(vehicles) == 1 {
		return vehicles[0], err
	}
	return nil, err

}

func (self *parkingLotDao) GetSlotsForParkingId(parkingLotId int) ([]*models.Slot, error) {

	var slots []*models.Slot
	// db, errs := db.DBConnect()
	db, errs := db.SqlxConnect()
	fmt.Println("the error is")
	fmt.Println(errs)
	sqlStatement := `SELECT
	                 slot_id	
				     FROM
				     slot 
				     WHERE parking_status=0 and parking_lot_id=$1
				`
	err := db.Select(&slots, sqlStatement, parkingLotId)

	if err != nil {
		fmt.Printf("The error is", err)
		return nil, err
	}

	// fmt.Printf("materials are", materials)
	defer db.Close()

	return slots, err

}
