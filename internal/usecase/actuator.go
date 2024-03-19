package usecase

import (
	"ara-server/internal/repository/db"
	"ara-server/util/log"
	"context"
)

func (uc *Usecase) InsertActuator(ctx context.Context, actuator Actuator) error {
	if err := uc.db.InsertActuator(ctx, db.Actuator(actuator)); err != nil {
		log.Error(ctx, actuator, err, "failed to insert actuator")
		return err
	}

	return nil
}
