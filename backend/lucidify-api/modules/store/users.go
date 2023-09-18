// store/users.go

package store

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
	Deleted          bool
}

func (s *Store) CreateUser(user User) error {
	query := `INSERT INTO users (user_id, external_id, username, password_enabled, email, first_name, last_name, image_url, profile_image_url, two_factor_enabled, created_at, updated_at, deleted) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`
	_, err := s.db.Exec(query, user.UserID, user.ExternalID, user.Username, user.PasswordEnabled, user.Email, user.FirstName, user.LastName, user.ImageURL, user.ProfileImageURL, user.TwoFactorEnabled, user.CreatedAt, user.UpdatedAt, user.Deleted)
	return err
}

func (s *Store) UpdateUser(user User) error {
	query := `UPDATE users SET external_id = $2, username = $3, password_enabled = $4, email = $5, first_name = $6, last_name = $7, image_url = $8, profile_image_url = $9, two_factor_enabled = $10, updated_at = $11 WHERE user_id = $1`
	_, err := s.db.Exec(query, user.UserID, user.ExternalID, user.Username, user.PasswordEnabled, user.Email, user.FirstName, user.LastName, user.ImageURL, user.ProfileImageURL, user.TwoFactorEnabled, user.UpdatedAt)
	return err
}

func (s *Store) DeleteUser(userID string) error {
	query := `UPDATE users SET deleted = TRUE WHERE user_id = $1`
	_, err := s.db.Exec(query, userID)
	return err
}
