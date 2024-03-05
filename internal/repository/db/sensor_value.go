package db

import (
	"ara-server/internal/constants"
	"database/sql"
	"math/rand"
	"time"
)

const timestampFormat = "2006-01-02 15:04:05"

func (repo *Repository) InsertSensorValue(sensorValue SensorValue) error {
	_, err := repo.db.Exec(queryInsertSensorValue, sensorValue.DeviceID, sensorValue.SensorType, sensorValue.Value, sensorValue.Time)
	return err
}

func (repo *Repository) InsertDummySensorValueData(deviceID int64, sensorType constants.SensorType) error {
	var result SensorValue
	if err := repo.db.Get(&result, "SELECT * FROM sensor_values WHERE device_id = $1 AND sensor_type = $2 ORDER BY time DESC LIMIT 1", deviceID, sensorType); err != nil && err != sql.ErrNoRows {
		return err
	}

	result.Time = time.Now()
	result.DeviceID = deviceID
	result.SensorType = sensorType
	nextVal := rand.Intn(5)
	if result.Value-nextVal >= 0 {
		multiplier := 1
		if rand.Intn(101) <= result.Value {
			multiplier = -1
		}
		nextVal *= multiplier
	}
	result.Value += nextVal

	if _, err := repo.db.Exec(queryInsertSensorValue, result.DeviceID, result.SensorType, result.Value, result.Time); err != nil {
		return err
	}

	return nil
}

func (repo *Repository) GetSensorValueTimeSeries(param GetSensorValueTimeSeriesParam) ([]SensorValueTimeSeriesItem, error) {
	if param.EndTime.IsZero() {
		param.EndTime = time.Now()
	}

	var result []SensorValueTimeSeriesItem
	err := repo.db.Select(&result, queryGetSensorValuesTimeSeries, param.StartTime.Format(timestampFormat), param.EndTime.Format(timestampFormat), param.DeviceID, param.SensorType)
	if err != nil {
		return nil, err
	}

	return result, nil
}
