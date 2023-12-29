package http

import (
	"ara-server/internal/constants"
	"time"
)

type GetSensorChartParam struct {
	StartTime  time.Time            `json:"start_time"`
	EndTime    time.Time            `json:"end_time,omitempty"`
	SensorType constants.SensorType `json:"sensor_type"`
	DeviceID   int64                `json:"device_id"`
}
