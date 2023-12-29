package db

const (
	queryGetLastAction = `
		SELECT * 
		FROM action_histories 
		WHERE 
			device_id = $1 AND action_type = $2 
		ORDER BY action_at DESC 
		LIMIT 1
	`

	queryInsertActionLog = `
		INSERT INTO action_histories (device_id, action_type, value, action_by, action_at) 
		VALUES ($1, $2, $3, $4, $5)
	`
)
