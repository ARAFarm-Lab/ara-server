package db

func (repo *Repository) InsertSensorValue(sensorValue SensorValue) error {
	_, err := repo.db.Exec(queryInsertSensorValue, sensorValue.DeviceID, sensorValue.SensorType, sensorValue.Value, sensorValue.Time)
	return err
}
