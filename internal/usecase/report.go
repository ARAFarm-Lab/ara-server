package usecase

import (
	"ara-server/internal/repository/db"
	"math"
)

func (uc *Usecase) GetSensorChart(param GetSensorChartParam) (SensorChartResponse, error) {
	values, err := uc.db.GetSensorValueTimeSeries(db.GetSensorValueTimeSeriesParam{
		DeviceID:   param.DeviceID,
		StartTime:  param.StartTime,
		EndTime:    param.EndTime,
		SensorType: param.SensorType,
	})
	if err != nil {
		return SensorChartResponse{}, err
	}

	chartConfig := uc.infra.GetConfig().Chart.SoilMoisture
	if len(values) == 0 {
		return SensorChartResponse{
			Data:          []SensorChartItem{},
			MaxPercentage: 100,
			MinPercentage: 0,
			MaxValue:      chartConfig.Max,
			MinValue:      chartConfig.Min,
		}, nil
	}

	maxValue := values[0].Value
	minValue := values[0].Value

	items := make([]SensorChartItem, len(values))
	for index, v := range values {
		maxValue = int(math.Max(float64(v.Value), float64(maxValue)))
		minValue = int(math.Min(float64(v.Value), float64(minValue)))

		// Greater the value, the less the percentage, means less water
		items[index] = SensorChartItem{
			Time:            v.Time,
			Value:           v.Value,
			ValuePercentage: 100 - calculatePercentage(v.Value, chartConfig.Min, chartConfig.Max),
		}
	}

	return SensorChartResponse{
		Data:          items,
		MinPercentage: 100 - calculatePercentage(minValue, chartConfig.Min, chartConfig.Max),
		MaxPercentage: 100 - calculatePercentage(maxValue, chartConfig.Min, chartConfig.Max),
		MinValue:      minValue,
		MaxValue:      maxValue,
	}, nil
}

func calculatePercentage(value, min, max int) int {
	if value < min {
		return 0
	}

	if value > max {
		return 100
	}

	return int(float64(value-min) / float64(max-min) * 100)
}
