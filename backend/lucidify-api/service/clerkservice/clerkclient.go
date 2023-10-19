package clerkservice

import (
	"context"
	"fmt"
	"lucidify-api/server/config"

	"github.com/clerkinc/clerk-sdk-go/clerk"
)

type ClerkClient interface {
	GetUserIDFromSession(ctx context.Context) (string, error)
	GetClerkClient() clerk.Client
}

type ClerkClientImpl struct {
	clerkClient clerk.Client
}

func NewClerkClient() (ClerkClient, error) {
	cfg := config.NewServerConfig()
	client, err := clerk.NewClient(cfg.ClerkSecretKey)
	if err != nil {
		return nil, err
	}
	return &ClerkClientImpl{clerkClient: client}, nil

}

func (c *ClerkClientImpl) GetUserIDFromSession(ctx context.Context) (string, error) {
	sessClaims, ok := ctx.Value(clerk.ActiveSessionClaims).(*clerk.SessionClaims)
	if !ok {
		// w.WriteHeader(http.StatusUnauthorized)
		// w.Write([]byte("Unauthorized"))
		return "unauthorized", fmt.Errorf("unauthorized")
	}

	user, err := c.clerkClient.Users().Read(sessClaims.Claims.Subject)
	if err != nil {
		panic(err)
	}
	return user.ID, nil
}

func (c *ClerkClientImpl) GetClerkClient() clerk.Client {
	return c.clerkClient
}
