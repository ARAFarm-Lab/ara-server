package handler

import (
	"ara-server/internal/usecase"

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
	h.mqtt.Subscribe("mq-iot-poc", 1, nil)
}

func (h *handler) ToggleLamp(c *gin.Context) {
	p := c.DefaultQuery("value", "false")
	h.mqtt.Publish("iot-poc-topic", 1, false, p)
	c.JSON(200, gin.H{
		"success": true,
	})
}
