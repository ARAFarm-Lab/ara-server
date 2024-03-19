package db

import (
	"context"
	"database/sql"
	"fmt"
)

func (repo *Repository) GetActionHistories(ctx context.Context, deviceID int64) ([]ActionHistory, error) {
	var actionHistories []ActionHistory
	err := repo.db.SelectContext(ctx, &actionHistories, queryGetActionHistory, deviceID)
	if err != nil && err != sql.ErrNoRows {
		return actionHistories, err
	}

	return actionHistories, nil
}

func (repo *Repository) GetLastActions(deviceID int64) ([]ActionHistory, error) {
	query := fmt.Sprintf(queryGetLastAction, "device_id = $1 ")

	var actionHistories []ActionHistory
	err := repo.db.Select(&actionHistories, query, deviceID)
	if err != nil && err != sql.ErrNoRows {
		return actionHistories, err
	}

	return actionHistories, nil
}

func (repo *Repository) GetLastActionByActuatorID(deviceID int64, actuatorID int64) (ActionHistory, error) {
	query := fmt.Sprintf(queryGetLastAction, "device_id = $1 AND actuator_id = $2 ")

	var actionHistory ActionHistory
	err := repo.db.Get(&actionHistory, query, deviceID, actuatorID)
	if err != nil && err != sql.ErrNoRows {
		return actionHistory, err
	}

	return actionHistory, nil
}

func (repo *Repository) InsertActionLog(param ActionHistory) error {
	_, err := repo.db.Exec(queryInsertActionLog, param.ActuatorID, param.Value, param.ActionBy, param.ActionAt)
	return err
}
