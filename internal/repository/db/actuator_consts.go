package db

const (
	queryGetActuatorsByFilter = `
		SELECT *
		FROM device_actuators
		%s
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
			icon = :icon,
			is_active = :is_active
		WHERE id = :id;
	`
)
