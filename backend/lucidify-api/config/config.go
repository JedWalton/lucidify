package config

import (
	"lucidify-api/store"
)

type ServerConfig struct {
	OPENAI_API_KEY string
	AllowedOrigins []string
	Port           string
	Store          *store.Store
}
