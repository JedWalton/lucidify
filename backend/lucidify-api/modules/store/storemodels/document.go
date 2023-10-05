package storemodels

import (
	"time"

	"github.com/google/uuid"
)

type Document struct {
	DocumentUUID uuid.UUID `db:"id"`
	UserID       string    `db:"user_id"`
	DocumentName string    `db:"document_name"`
	Content      string    `db:"content"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
