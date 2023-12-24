package db

import (
	"ara-server/internal/constants"
	"time"
)

type SensorValue struct {
	DeviceID   int64                `db:"device_id"`
	SensorType constants.SensorType `db:"sensor_type"`
	Value      int                  `db:"value"`
	Time       time.Time            `db:"time"`
}
