package db

const (
	queryInsertSensorValue = `
		INSERT INTO sensor_values (device_id, sensor_type, value, time)
		VALUES ($1, $2, $3, $4)
	`
)
