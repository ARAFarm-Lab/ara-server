package http

type DispatcherRequest struct {
	DeviceID   int64       `json:"device_id"`
	ActuatorID int64       `json:"actuator_id"`
	Value      interface{} `json:"value"`
}
