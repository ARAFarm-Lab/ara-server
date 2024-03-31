package db

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

func (repo *Repository) GetUserByEmail(ctx context.Context, email string) (User, error) {
	var result User
	if err := repo.db.GetContext(ctx, &result, queryGetUserByEmail, email); err != nil && err != sql.ErrNoRows {
		return User{}, err
	}

	return result, nil
}

func (repo *Repository) GetUserByID(ctx context.Context, id int) (User, error) {
	var result User
	if err := repo.db.GetContext(ctx, &result, queryGetUserByID, id); err != nil && err != sql.ErrNoRows {
		return User{}, err
	}

	return result, nil
}

func (repo *Repository) GetUserInfoByUserID(ctx context.Context, userID int) (UserInfo, error) {
	var result UserInfo
	if err := repo.db.GetContext(ctx, &result, queryGetUserInfoByUserID, userID); err != nil && err != sql.ErrNoRows {
		return UserInfo{}, err
	}

	return result, nil
}

func (repo *Repository) GetUserInfoList(ctx context.Context) ([]UserInfo, error) {
	var result []UserInfo
	if err := repo.db.SelectContext(ctx, &result, queryGetUserInfoList); err != nil {
		return nil, err
	}

	return result, nil
}

func (repo *Repository) GetUsers(ctx context.Context) ([]User, error) {
	var result []User
	if err := repo.db.SelectContext(ctx, &result, queryGetUsers); err != nil {
		return nil, err
	}

	return result, nil
}

func (repo *Repository) InsertUser(ctx context.Context, tx *sql.Tx, param InsertUserParam) (int, error) {
	query, args, err := sqlx.Named(queryInsertUser, param)
	if err != nil {
		return 0, err
	}

	var id int
	if err := tx.QueryRowContext(ctx, repo.Rebind(query), args...).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (repo *Repository) InsertProfile(ctx context.Context, tx *sql.Tx, param InsertProfileParam) error {
	query, args, err := sqlx.Named(queryInsertProfile, param)
	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, repo.Rebind(query), args...); err != nil {
		return err
	}

	return nil
}

func (repo *Repository) UpdateUser(ctx context.Context, tx *sql.Tx, user User) error {
	query, args, err := sqlx.Named(queryUpdateUser, user)
	if err != nil {
		return err
	}

	_, err = repo.db.ExecContext(ctx, repo.Rebind(query), args...)
	if err != nil {
		return err
	}

	return nil
}

func (repo *Repository) UpdateUserProfile(ctx context.Context, tx *sql.Tx, profile UserInfo) error {
	query, args, err := sqlx.Named(queryUpdateUserProfile, profile)
	if err != nil {
		return err
	}

	_, err = repo.db.ExecContext(ctx, repo.Rebind(query), args...)
	if err != nil {
		return err
	}

	return nil
}
