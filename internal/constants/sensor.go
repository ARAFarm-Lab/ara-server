package constants

type SensorType uint8

const (
	TemperatureSensor SensorType = iota + 1
	HumiditySensor
	SoilSensor
)
