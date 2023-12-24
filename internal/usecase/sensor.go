package usecase

import (
	"ara-server/internal/constants"
	"ara-server/internal/repository/db"
	"fmt"
	"time"
)

func (uc *Usecase) StoreSensorValue(param StoreSensorValueParam) error {
	var errs []error

	timeNow := time.Now()
	for _, value := range param.SensorValues {
		if err := uc.db.InsertSensorValue(db.SensorValue{
			DeviceID:   param.DeviceID,
			SensorType: constants.SensorType(value.SensorType),
			Value:      value.Value,
			Time:       timeNow,
		}); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("error storing sensor values: %v", errs)
	}

	return nil
}
