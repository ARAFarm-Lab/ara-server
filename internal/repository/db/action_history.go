package db

import (
	"ara-server/internal/constants"
	"database/sql"
)

func (repo *Repository) GetLastAction(deviceID int64, actionType constants.ActionType) (ActionHistory, error) {
	var actionHistory ActionHistory
	err := repo.db.Get(&actionHistory, queryGetLastAction, deviceID, actionType)
	if err != nil && err != sql.ErrNoRows {
		return actionHistory, err
	}

	return actionHistory, nil
}

func (repo *Repository) InsertActionLog(param ActionHistory) error {
	_, err := repo.db.Exec(queryInsertActionLog, param.DeviceID, param.ActionType, param.Value, param.ActionBy, param.ActionAt)
	return err
}
