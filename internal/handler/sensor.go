package handler

import (
	"ara-server/internal/usecase"
	"encoding/json"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func (h *handler) HandleSensorRead(client mqtt.Client, msg mqtt.Message) {
	deviceID := getDeviceID(msg.Topic())
	if deviceID < 0 {
		return
	}

	var payload [][]int
	if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
		fmt.Println(err)
	}

	sensorData := make([]usecase.SensorValue, len(payload))
	for i, p := range payload {
		sensorData[i] = usecase.SensorValue{
			SensorType: p[0],
			Value:      p[1],
		}
	}

	h.usecase.StoreSensorValue(usecase.StoreSensorValueParam{
		DeviceID:     deviceID,
		SensorValues: sensorData,
	})
}
