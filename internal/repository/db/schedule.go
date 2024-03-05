package db

import (
	"ara-server/internal/constants"
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

func (repo *Repository) BulkUpdateActionScheduleStatusPending(scheduleIDs []int) error {
	_, err := repo.db.Exec(queryBulkUpdateScheduleStatusPending, pq.Array(scheduleIDs))
	if err != nil {
		return err
	}

	return nil
}

func (repo *Repository) DeleteScheduleByID(ctx context.Context, scheduleID int) error {
	result, err := repo.db.ExecContext(ctx, queryDeleteScheduleByID, scheduleID)
	if err != nil {
		return err
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return constants.ErrorScheduleNotFound
	}

	return nil
}

func (repo *Repository) GetScheduleByID(ctx context.Context, id int) (ActionSchedule, error) {
	var result ActionSchedule
	if err := repo.db.GetContext(ctx, &result, queryGetScheduleByID, id); err != nil {
		return ActionSchedule{}, err
	}

	return result, nil
}

func (repo *Repository) GetScheduledAction() ([]ActionSchedule, error) {
	var result []ActionSchedule
	if err := repo.db.Select(&result, queryGetScheduledAction); err != nil {
		return nil, err
	}

	return result, nil
}

func (repo *Repository) GetUpcomingSchedules() ([]ActionSchedule, error) {
	var result []ActionSchedule
	if err := repo.db.Select(&result, queryGetUpcomingSchedules); err != nil {
		return nil, err
	}

	return result, nil
}

func (repo *Repository) InsertActionSchedule(schedule ActionSchedule) error {
	query, args, err := sqlx.Named(queryInsertActionSchedule, schedule)
	if err != nil {
		return err
	}

	_, err = repo.db.Exec(repo.Rebind(query), args...)
	if err != nil {
		return err
	}

	return nil
}

func (repo *Repository) UpdateActionSchedule(ctx context.Context, action ActionSchedule) error {
	query, args, err := sqlx.Named(queryUpdateActionSchedule, action)
	if err != nil {
		return err
	}

	if _, err := repo.db.ExecContext(ctx, repo.Rebind(query), args...); err != nil {
		return err
	}
	return nil
}
