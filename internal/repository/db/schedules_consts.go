package db

const (
	queryBulkUpdateScheduleStatusPending = `
		UPDATE
			schedules
		SET
			last_run_status = 1
		WHERE
			id = ANY($1)
	`

	queryGetScheduledAction = `
		SELECT
			id,
			name,
			description,
			actions,
			schedule,
			is_active,
			next_run_at,
			last_lock_at,
			last_run_at,
			last_run_status,
			created_at
		FROM
			schedules
		WHERE
			is_active = true
			AND next_run_at <= NOW()
			AND last_lock_at IS NULL
		ORDER BY
			next_run_at ASC
	`

	queryUpdateActionSchedule = `
		UPDATE
			schedules
		SET
			name = :name,
			description = :description,
			actions = :actions,
			schedule = :schedule,
			is_active = :is_active,
			next_run_at = :next_run_at,
			last_lock_at = :last_lock_at,
			last_run_at = :last_run_at,
			last_run_status = :last_run_status
		WHERE
			id = :id
	`
)
