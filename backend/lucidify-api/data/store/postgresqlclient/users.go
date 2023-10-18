// postgresqlclient/users.go

package postgresqlclient

import (
	"lucidify-api/data/store/storemodels"
)

func (s *PostgreSQL) CreateUserInUsersTable(user storemodels.User) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() // Make sure it's idempotent

	query := `INSERT INTO users (user_id, external_id, username, password_enabled, email, first_name, last_name, image_url, profile_image_url, two_factor_enabled, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	_, err = tx.Exec(query, user.UserID, user.ExternalID, user.Username, user.PasswordEnabled, user.Email, user.FirstName, user.LastName, user.ImageURL, user.ProfileImageURL, user.TwoFactorEnabled, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *PostgreSQL) UpdateUserInUsersTable(user storemodels.User) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `UPDATE users SET external_id = $2, username = $3, password_enabled = $4, email = $5, first_name = $6, last_name = $7, image_url = $8, profile_image_url = $9, two_factor_enabled = $10, updated_at = $11 WHERE user_id = $1`
	_, err = tx.Exec(query, user.UserID, user.ExternalID, user.Username, user.PasswordEnabled, user.Email, user.FirstName, user.LastName, user.ImageURL, user.ProfileImageURL, user.TwoFactorEnabled, user.UpdatedAt)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *PostgreSQL) GetUserInUsersTable(userID string) (*storemodels.User, error) {
	query := `SELECT user_id, external_id, username, password_enabled, email, first_name, last_name, image_url, profile_image_url, two_factor_enabled, created_at, updated_at FROM users WHERE user_id = $1`
	row := s.db.QueryRow(query, userID)
	var user storemodels.User
	err := row.Scan(&user.UserID, &user.ExternalID, &user.Username, &user.PasswordEnabled, &user.Email, &user.FirstName, &user.LastName, &user.ImageURL, &user.ProfileImageURL, &user.TwoFactorEnabled, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *PostgreSQL) DeleteUserInUsersTable(userID string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `DELETE FROM users WHERE user_id = $1`
	_, err = tx.Exec(query, userID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
