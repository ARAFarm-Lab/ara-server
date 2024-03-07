package db

const (
	queryGetActionHistory = `
		SELECT *
		FROM action_histories
		WHERE device_id = $1
		ORDER BY action_at DESC
		LIMIT 10
	`

	queryGetLastAction = `
		SELECT DISTINCT ON (action_type) * 
		FROM action_histories 
		WHERE %s
		ORDER BY action_type, action_at DESC
	`

	queryInsertActionLog = `
		INSERT INTO action_histories (device_id, action_type, value, action_by, action_at) 
		VALUES ($1, $2, $3, $4, $5)
	`
)
