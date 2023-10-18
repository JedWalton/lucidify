package clerkservice

import (
	"lucidify-api/server/config"

	"github.com/clerkinc/clerk-sdk-go/clerk"
)

func NewClerkClient() (clerk.Client, error) {
	cfg := config.NewServerConfig()
	client, err := clerk.NewClient(cfg.ClerkSecretKey)
	if err != nil {
		return nil, err
	}
	return client, nil
}
