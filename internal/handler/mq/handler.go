package mq

import (
	"ara-server/internal/usecase"
	"ara-server/util/log"
	"strconv"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type handler struct {
	mqtt mqtt.Client

	usecase *usecase.Usecase
}

func InitHandler(usecase *usecase.Usecase, mqtt mqtt.Client) *handler {
	handler := &handler{usecase: usecase, mqtt: mqtt}
	handler.registerMQHandler()

	return handler
}

func (h *handler) registerMQHandler() {
	h.mqtt.Subscribe("sensor-read/#", 1, h.HandleSensorRead)
	h.mqtt.Subscribe("device-initial-state-request/#", 1, h.HandleInitiateDeviceState)
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
