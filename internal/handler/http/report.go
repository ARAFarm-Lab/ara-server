package http

import (
	"ara-server/internal/usecase"

	"github.com/gin-gonic/gin"
)

func (h *handler) HandleGetSensorChart(c *gin.Context) error {
	var param GetSensorChartParam
	if err := c.ShouldBindJSON(&param); err != nil {
		WriteJson(c, nil, errInvalidRequest)
		return errInvalidRequest
	}

	chart, err := h.usecase.GetSensorChart(usecase.GetSensorChartParam{
		DeviceID:   param.DeviceID,
		StartTime:  param.StartTime,
		EndTime:    param.EndTime,
		SensorType: param.SensorType,
	})
	if err != nil {
		WriteJson(c, nil, err)
		return err
	}

	WriteJson(c, chart, nil)
	return nil
}
