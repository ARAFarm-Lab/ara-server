package usecase

type deviceConfig struct {
	Actuators [][]interface{} `json:"a"` // {deviceID, pinNumber}
	Values    [][]interface{} `json:"v"` // {deviceID, actionValue}
}
