package mq

import (
	"ara-server/internal/constants"
	"ara-server/internal/usecase"
	"ara-server/util/log"
	"context"
	"strconv"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/xid"
)

type handler struct {
	mqtt mqtt.Client

	usecase *usecase.Usecase
}

func InitHandler(usecase *usecase.Usecase, mqtt mqtt.Client) *handler {
	handler := &handler{
		usecase: usecase,
		mqtt:    mqtt,
	}
	handler.registerMQHandler()

	return handler
}

func (h *handler) registerMQHandler() {
	h.mqtt.Subscribe("sensor-read/#", 1, h.mqWrapper(h.HandleSensorRead))
	h.mqtt.Subscribe("device-initial-state-request/#", 1, h.mqWrapper(h.HandleInitiateDeviceState))
}

func getDeviceID(topic string) int64 {
	topicSegments := strings.Split(topic, "/")
	if len(topicSegments) < 1 {
		return -1
	}
	deviceID, err := strconv.ParseInt(topicSegments[len(topicSegments)-1], 10, 64)
	if err != nil {
		log.Error(context.Background(), nil, err, "error parsing device id")
		return -1
	}

	return deviceID
}

func (h *handler) mqWrapper(handler func(ctx context.Context, client mqtt.Client, msg mqtt.Message)) mqtt.MessageHandler {
	return func(client mqtt.Client, message mqtt.Message) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, constants.CtxKeyCtxID, xid.New().String())
		handler(ctx, client, message)
	}
}
