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

type SensorValueTimeSeriesItem struct {
	Value int       `db:"value"`
	Time  time.Time `db:"time"`
}

type GetSensorValueTimeSeriesParam struct {
	DeviceID   int64                `db:"device_id"`
	StartTime  time.Time            `db:"start_time"`
	EndTime    time.Time            `db:"end_time"`
	SensorType constants.SensorType `db:"sensor_type"`
}
