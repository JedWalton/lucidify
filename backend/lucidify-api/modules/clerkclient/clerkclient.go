package clerkclient

import (
	"github.com/clerkinc/clerk-sdk-go/clerk"
)

func NewClerkClient(CLERK_API_KEY string) (clerk.Client, error) {
	client, err := clerk.NewClient(CLERK_API_KEY)
	if err != nil {
		return nil, err
	}
	return client, nil
}
