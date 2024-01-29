package db

import (
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

func (repo *Repository) UpdateActionSchedule(action ActionSchedule) error {
	query, args, err := sqlx.Named(queryUpdateActionSchedule, action)
	if err != nil {
		return err
	}

	_, err = repo.db.Exec(repo.Rebind(query), args...)
	if err != nil {
		return err
	}

	return nil
}
