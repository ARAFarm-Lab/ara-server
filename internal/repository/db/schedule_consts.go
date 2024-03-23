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

	queryDeleteScheduleByID = `
		DELETE FROM
			schedules
		WHERE
			id = $1
	`

	queryGetScheduleByID = `
		SELECT * FROM schedules WHERE id = $1
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
			cleanup_time,
			duration_minute,
			last_lock_at,
			last_run_at,
			last_run_status,
			created_at
		FROM
			schedules
		WHERE 
			is_active = true
			AND (next_run_at <= NOW() OR cleanup_time <= NOW())
			AND last_lock_at IS NULL
	`

	queryGetUpcomingSchedules = `
		SELECT
			id,
			name,
			description,
			actions,
			schedule,
			duration_minute,
			cleanup_time,
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
			AND next_run_at >= NOW() OR cleanup_time >= NOW()
		ORDER BY
			COALESCE(next_run_at, '0001-12-31'::TIMESTAMP) ASC, COALESCE(cleanup_time, '0001-12-31'::TIMESTAMP) ASC
	`

	queryInsertActionSchedule = `
		INSERT INTO
			schedules(
				name, 
				description, 
				actions, 
				schedule,
				duration_minute,
				is_active, 
				next_run_at, 
				last_lock_at, 
				last_run_at, 
				last_run_status
			)
		VALUES(
			:name,
			:description,
			:actions,
			:schedule,
			:duration_minute,
			:is_active,
			:next_run_at,
			:last_lock_at,
			:last_run_at,
			:last_run_status
		)
	`

	queryUpdateActionSchedule = `
		UPDATE
			schedules
		SET
			name = :name,
			description = :description,
			actions = :actions,
			duration_minute = :duration_minute,
			schedule = :schedule,
			is_active = :is_active,
			next_run_at = :next_run_at,
			cleanup_time = :cleanup_time,
			last_lock_at = :last_lock_at,
			last_run_at = :last_run_at,
			last_run_status = :last_run_status
		WHERE
			id = :id
	`
)
