package db

const (
	queryGetActiveActuators = `
		SELECT *
		FROM device_actuators
		WHERE device_id = $1 AND is_active = true;
	`

	queryGetActuatorByID = `
		SELECT *
		FROM device_actuators
		WHERE id = $1
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
)
