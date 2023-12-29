package db

import "time"

const timestampFormat = "2006-01-02 15:04:05"

func (repo *Repository) InsertSensorValue(sensorValue SensorValue) error {
	_, err := repo.db.Exec(queryInsertSensorValue, sensorValue.DeviceID, sensorValue.SensorType, sensorValue.Value, sensorValue.Time)
	return err
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
