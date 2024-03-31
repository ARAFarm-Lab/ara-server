package db

const (
	queryInsertUser = `
		INSERT INTO users (email, password)
		VALUES (:email, :password)
		RETURNING id;
	`

	queryInsertProfile = `
		INSERT INTO user_profiles (user_id, name, role)
		VALUES (:user_id, :name, :role);
	`

	queryGetUsers = `
		SELECT * FROM users
	`

	queryGetUserByEmail = `
		SELECT * FROM users
		WHERE LOWER(email) = LOWER($1)
	`

	queryGetUserByID = `
		SELECT * FROM users
		WHERE id = $1
	`

	queryUpdateUser = `
		UPDATE users SET 
			is_active = :is_active,
			updated_at = NOW(),
			password = :password
		WHERE id = :id
	`

	queryGetUserInfoByUserID = `
		SELECT profile.name, profile.role, "user".is_active FROM user_profiles profile
		INNER JOIN users "user" 
		ON "user".id = profile.user_id
		WHERE "user".id = $1
	`

	queryGetUserInfoList = `
		SELECT "user".id as user_id, "user".is_active, profile.name, profile.role FROM user_profiles profile
		INNER JOIN users "user" 
		ON "user".id = profile.user_id
	`

	queryUpdateUserProfile = `
		UPDATE user_profiles
		SET
			name = :name,
			role = :role,
			updated_at = NOW()
		WHERE user_id = :user_id;
	`
)
