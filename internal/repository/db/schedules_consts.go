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
			duration_minute,
			last_lock_at,
			last_run_at,
			last_run_status,
			created_at
		FROM
			schedules
		WHERE
			(
				schedule IS NULL
				AND is_active = true
				AND (next_run_at <= NOW() OR next_run_at + INTERVAL '1 minute' * duration_minute <= NOW())
				AND (last_run_status IS NULL OR last_run_status != 3)
				AND last_lock_at IS NULL
			)
			OR
			(
				schedule IS NOT NULL
				AND is_active = true
				AND (next_run_at <= NOW() OR next_run_at + INTERVAL '1 minute' * duration_minute <= NOW())
				AND last_lock_at IS NULL
			)
		ORDER BY
			next_run_at ASC
	`

	queryGetUpcomingSchedules = `
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
			AND next_run_at >= NOW()
		ORDER BY
			next_run_at ASC
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
