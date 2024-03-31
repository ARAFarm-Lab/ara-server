package db

const (
	queryGetActionHistory = `
		SELECT 
			actuator.id AS actuator_id,
			actuator.action_type,
			actuator.name,
			actuator.icon,
			history.value,
			history.action_by,
			history.action_at
		FROM action_histories history
		INNER JOIN device_actuators actuator ON history.actuator_id = actuator.id
		WHERE actuator.device_id = $1
		ORDER BY action_at DESC
		LIMIT 10
	`

	queryGetLastAction = `
		SELECT DISTINCT ON (actuator.id)
			actuator.id AS actuator_id,
			actuator.action_type,
			history.value,
			history.action_by,
			history.action_at
		FROM action_histories history
		INNER JOIN device_actuators actuator ON history.actuator_id = actuator.id
		WHERE %s
		ORDER BY actuator.id, action_at DESC
	`

	queryInsertActionLog = `
		INSERT INTO action_histories (actuator_id, value, action_by, action_at) 
		VALUES ($1, $2, $3, $4)
	`
)
