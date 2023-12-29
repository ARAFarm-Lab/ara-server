package usecase

import (
	"ara-server/internal/constants"
	"time"
)

type GetSensorChartParam struct {
	StartTime  time.Time
	EndTime    time.Time
	SensorType constants.SensorType
	DeviceID   int64
}

type SensorChartResponse struct {
	Data          []SensorChartItem `json:"data"`
	MinPercentage int               `json:"min_percentage"`
	MaxPercentage int               `json:"max_percentage"`
	MinValue      int               `json:"min_value"`
	MaxValue      int               `json:"max_value"`
}

type SensorChartItem struct {
	Time            time.Time `json:"time"`
	ValuePercentage int       `json:"value_percentage"`
	Value           int       `json:"value"`
}
