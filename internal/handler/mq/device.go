package mq

import (
	"ara-server/util/log"
	"context"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func (h *handler) HandleInitiateDeviceState(ctx context.Context, msg mqtt.Message) {
	log.Info(ctx, nil, nil, "handling initiate device state")

	deviceID := getDeviceID(msg.Topic())
	if deviceID < 0 {
		return
	}

	h.usecase.InitiateDeviceState(ctx, deviceID)
}

func (h *handler) HandleHeartbeatResponse(ctx context.Context, msg mqtt.Message) {

}
