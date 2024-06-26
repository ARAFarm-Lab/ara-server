package usecase

import (
	"ara-server/internal/constants"
	"ara-server/internal/repository/db"
	"ara-server/util/log"
	"context"
)

func (uc *Usecase) GetActuators(ctx context.Context, deviceID int64) ([]Actuator, error) {
	actuators, err := uc.db.GetActuatorsByFilter(ctx, nil)
	if err != nil {
		log.Error(ctx, deviceID, err, "failed to get actuators")
		return nil, err
	}

	result := make([]Actuator, len(actuators))
	for index, actuator := range actuators {
		result[index] = Actuator(actuator)
	}

	return result, nil
}

func (uc *Usecase) InsertActuator(ctx context.Context, actuator Actuator) error {
	if err := uc.db.InsertActuator(ctx, db.Actuator(actuator)); err != nil {
		log.Error(ctx, actuator, err, "failed to insert actuator")
		return err
	}

	return nil
}

func (uc *Usecase) UpdateActuator(ctx context.Context, actuator Actuator) error {
	existing, err := uc.db.GetActuatorByID(ctx, actuator.ID)
	if err != nil {
		log.Error(ctx, actuator, err, "failed to get actuator")
		return err
	}

	if existing.ID == 0 {
		return errorActuatorNotFound
	}

	existing.Name = actuator.Name
	existing.Icon = actuator.Icon
	existing.IsActive = actuator.IsActive

	if role := ctx.Value(string(constants.CtxKeyUserRole)).(constants.UserRole); role == constants.RoleAdmin {
		existing.PinNumber = actuator.PinNumber
		existing.TerminalNumber = actuator.TerminalNumber
	}

	// If it is disabling the actuator, turn off the actuator
	if !actuator.IsActive {
		if err := uc.DispatchAction(ctx, DispatcherParam{
			DeviceID:   existing.DeviceID,
			ActuatorID: existing.ID,
			Value:      false,
			ActionBy:   constants.ActionSourceScheduler,
		}); err != nil {
			log.Error(ctx, existing, err, "failed to dispatch action")
			return err
		}
	}

	if err := uc.db.UpdateActuator(ctx, existing); err != nil {
		log.Error(ctx, actuator, err, "failed to update actuator")
		return err
	}

	return nil
}
