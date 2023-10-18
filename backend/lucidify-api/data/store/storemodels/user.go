package storemodels

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
