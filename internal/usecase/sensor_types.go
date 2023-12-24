package usecase

type StoreSensorValueParam struct {
	DeviceID     int64
	SensorValues []SensorValue
}

type SensorValue struct {
	SensorType int
	Value      int
}
