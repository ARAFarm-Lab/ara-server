package db

const (
	queryInsertSensorValue = `
		INSERT INTO sensor_values (device_id, sensor_type, value, time)
		VALUES ($1, $2, $3, $4)
	`

	queryGetSensorValuesTimeSeries = `
		SELECT
			value,
			time
		FROM 
			sensor_values
		WHERE 
			time BETWEEN $1 AND $2 
		AND device_id = $3 AND sensor_type = $4
		ORDER BY time ASC
	`
)
