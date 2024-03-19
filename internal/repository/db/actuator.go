package db

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

func (repo *Repository) GetActiveActuators(ctx context.Context, deviceID int64) ([]Actuator, error) {
	var result []Actuator
	if err := repo.db.SelectContext(ctx, &result, queryGetActiveActuators, deviceID); err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return result, nil
}

func (repo *Repository) GetActuatorByID(ctx context.Context, id int64) (Actuator, error) {
	var result Actuator
	if err := repo.db.GetContext(ctx, &result, queryGetActuatorByID, id); err != nil && err != sql.ErrNoRows {
		return Actuator{}, err
	}

	return result, nil
}

func (repo *Repository) InsertActuator(ctx context.Context, actuator Actuator) error {
	query, args, err := sqlx.Named(queryInsertActuator, actuator)
	if err != nil {
		return err
	}

	if _, err = repo.db.ExecContext(ctx, repo.db.Rebind(query), args...); err != nil {
		return err
	}

	return nil
}
