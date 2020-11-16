CREATE TABLE parking_lot(
parking_lot_id serial PRIMARY KEY,
slot_capacity INT NOT NULL 
);

CREATE TABLE slot(
slot_id serial PRIMARY KEY,
parking_lot_id INT NOT NULL REFERENCES parking_lot(parking_lot_id),
parking_status INT NOT NULL 
);

CREATE TABLE vehicle(
vehicle_id serial PRIMARY KEY,
slot_id INT REFERENCES slot(slot_id),
plate_number VARCHAR(15) NOT NULL,
color VARCHAR(15),
parking_status INT NOT NULL,
parking_time TIMESTAMP  
);