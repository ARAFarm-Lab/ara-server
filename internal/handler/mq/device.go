package mq

import (
	"ara-server/util/log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func (h *handler) HandleInitiateDeviceState(client mqtt.Client, msg mqtt.Message) {
	log.Info("handling initiate device state")

	deviceID := getDeviceID(msg.Topic())
	if deviceID < 0 {
		return
	}

	h.usecase.InitiateDeviceState(deviceID)
}
