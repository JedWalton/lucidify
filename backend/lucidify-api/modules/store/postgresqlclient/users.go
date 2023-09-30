// postgresqlclient/users.go

package postgresqlclient

import (
	"fmt"
	"time"
)

type User struct {
	UserID           string
	ExternalID       string
	Username         string
	PasswordEnabled  bool
	Email            string
	FirstName        string
	LastName         string
	ImageURL         string
	ProfileImageURL  string
	TwoFactorEnabled bool
	CreatedAt        int64
	UpdatedAt        int64
}

// Add all these to a userstable interface.

func (s *PostgreSQL) CreateUserInUsersTable(user User) error {
	query := `INSERT INTO users (user_id, external_id, username, password_enabled, email, first_name, last_name, image_url, profile_image_url, two_factor_enabled, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	_, err := s.db.Exec(query, user.UserID, user.ExternalID, user.Username, user.PasswordEnabled, user.Email, user.FirstName, user.LastName, user.ImageURL, user.ProfileImageURL, user.TwoFactorEnabled, user.CreatedAt, user.UpdatedAt)
	return err
}

func (s *PostgreSQL) UpdateUserInUsersTable(user User) error {
	query := `UPDATE users SET external_id = $2, username = $3, password_enabled = $4, email = $5, first_name = $6, last_name = $7, image_url = $8, profile_image_url = $9, two_factor_enabled = $10, updated_at = $11 WHERE user_id = $1`
	_, err := s.db.Exec(query, user.UserID, user.ExternalID, user.Username, user.PasswordEnabled, user.Email, user.FirstName, user.LastName, user.ImageURL, user.ProfileImageURL, user.TwoFactorEnabled, user.UpdatedAt)
	return err
}

func (s *PostgreSQL) GetUserInUsersTable(userID string) (*User, error) {
	query := `SELECT user_id, external_id, username, password_enabled, email, first_name, last_name, image_url, profile_image_url, two_factor_enabled, created_at, updated_at FROM users WHERE user_id = $1`
	row := s.db.QueryRow(query, userID)
	var user User
	err := row.Scan(&user.UserID, &user.ExternalID, &user.Username, &user.PasswordEnabled, &user.Email, &user.FirstName, &user.LastName, &user.ImageURL, &user.ProfileImageURL, &user.TwoFactorEnabled, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *PostgreSQL) DeleteUserInUsersTable(userID string) error {
	query := `DELETE FROM users WHERE user_id = $1`
	_, err := s.db.Exec(query, userID)
	return err
}

func (s *PostgreSQL) CheckIfUserInUsersTable(userID string, retries int) error {
	for i := 0; i < retries; i++ {
		_, err := s.GetUserInUsersTable(userID)
		if err == nil {
			return nil
		}
		time.Sleep(time.Second) // Wait for 1 second before retrying
	}
	return fmt.Errorf("User not found after %d retries", retries)
}

func (s *PostgreSQL) CheckUserDeletedInUsersTable(userID string, retries int) error {
	for i := 0; i < retries; i++ {
		_, err := s.GetUserInUsersTable(userID)
		if err != nil {
			// If the user is not found, it means the user has been deleted
			return nil
		}
		time.Sleep(time.Second) // Wait for 1 second before retrying
	}
	return fmt.Errorf("User still exists in the database after %d retries", retries)
}

func (s *PostgreSQL) CheckUserHasExpectedFirstNameAndLastNameInUsersTable(userID string, retries int, expectedFirstName string, expectedLastName string) error {
	for i := 0; i < retries; i++ {
		user, err := s.GetUserInUsersTable(userID)
		if err == nil && user.FirstName == expectedFirstName && user.LastName == expectedLastName {
			return nil
		}
		time.Sleep(time.Second) // Wait for 1 second before retrying
	}
	return fmt.Errorf("User not updated correctly after %d retries", retries)
}
