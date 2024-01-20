package mq

import (
	"ara-server/util/log"
	"context"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func (h *handler) HandleInitiateDeviceState(ctx context.Context, client mqtt.Client, msg mqtt.Message) {
	log.Info(ctx, nil, nil, "handling initiate device state")

	deviceID := getDeviceID(msg.Topic())
	if deviceID < 0 {
		return
	}

	h.usecase.InitiateDeviceState(ctx, deviceID)
}
