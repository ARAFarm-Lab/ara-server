package usecase

import "ara-server/internal/constants"

type StoreSensorDataParam struct {
	SensorType constants.SensorType
	DeviceSN   string
	Timestamp  int64
	Value      string // JSON string, e.g. {"temperature": 25.5}
}
