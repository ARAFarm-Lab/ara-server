package handler

import (
	"ara-server/internal/usecase"
	"ara-server/util/log"
	"strconv"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
)

type handler struct {
	mqtt mqtt.Client

	usecase *usecase.Usecase
}

func NewHandler(usecase *usecase.Usecase, mqtt mqtt.Client) *handler {
	handler := &handler{usecase: usecase, mqtt: mqtt}
	handler.registerMQHandler()

	return handler
}

func (h *handler) RegisterHTTPHandler(router *gin.Engine) {
	router.POST("/toggle", h.ToggleLamp)
}

func (h *handler) registerMQHandler() {
	h.mqtt.Subscribe("sensor-read/#", 1, h.HandleSensorRead)
}

func getDeviceID(topic string) int64 {
	topicSegments := strings.Split(topic, "/")
	if len(topicSegments) < 1 {
		return -1
	}
	deviceID, err := strconv.ParseInt(topicSegments[len(topicSegments)-1], 10, 64)
	if err != nil {
		log.Error(err, "error parsing device id", topicSegments)
		return -1
	}

	return deviceID
}

func (h *handler) ToggleLamp(c *gin.Context) {
	p := c.DefaultQuery("value", "false")
	h.mqtt.Publish("iot-poc-topic", 1, false, p)
	c.JSON(200, gin.H{
		"success": true,
	})
}
