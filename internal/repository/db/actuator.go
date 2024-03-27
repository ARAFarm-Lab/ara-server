package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

func (repo *Repository) GetActuatorsByFilter(ctx context.Context, filter []GetActuatorsFilter) ([]Actuator, error) {
	query, args := buildGetActuatorsByFilterQuery(filter)

	var result []Actuator
	if err := repo.db.SelectContext(ctx, &result, query, args...); err != nil && err != sql.ErrNoRows {
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

func (repo *Repository) UpdateActuator(ctx context.Context, actuator Actuator) error {
	query, args, err := sqlx.Named(queryUpdateActuatorByID, actuator)
	if err != nil {
		return err
	}

	if _, err = repo.db.ExecContext(ctx, repo.db.Rebind(query), args...); err != nil {
		return err
	}

	return nil
}

func buildGetActuatorsByFilterQuery(filter []GetActuatorsFilter) (string, []interface{}) {
	if len(filter) == 0 {
		return fmt.Sprintf(queryGetActuatorsByFilter, ""), nil
	}

	args := make([]interface{}, 0, len(filter))
	clauses := make([]string, 0, len(filter))
	for index, item := range filter {
		args = append(args, item.Value)
		clauses = append(clauses, fmt.Sprintf("%s = $%d", item.Name, index+1))
	}

	return fmt.Sprintf(queryGetActuatorsByFilter, "WHERE "+strings.Join(clauses, " AND ")+" "), args
}
