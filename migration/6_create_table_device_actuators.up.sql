CREATE TABLE IF NOT EXISTS device_actuators (
    id SERIAL PRIMARY KEY,
    device_id INTEGER NOT NULL,
    pin_number INTEGER NOT NULL,
    action_type SMALLINT NOT NULL,
    name VARCHAR(50),
    icon VARCHAR(20),
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);