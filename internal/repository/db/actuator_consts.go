package db

const (
	queryGetActiveActuators = `
		SELECT *
		FROM device_actuators
		WHERE device_id = $1 AND is_active = true
		ORDER BY terminal_number ASC;
	`

	queryGetActuatorByID = `
		SELECT *
		FROM device_actuators
		WHERE id = $1;
	`

	queryInsertActuator = `
		INSERT INTO device_actuators (
			device_id,
			pin_number,
			action_type,
			name,
			icon
		) VALUES (
			:device_id,
			:pin_number,
			:action_type,
			:name,
			:icon
		);
	`

	queryUpdateActuatorByID = `
		UPDATE device_actuators
		SET
			device_id = :device_id,
			pin_number = :pin_number,
			action_type = :action_type,
			name = :name,
			icon = :icon
		WHERE id = :id;
	`
)
