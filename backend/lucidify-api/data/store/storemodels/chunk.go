package storemodels

import (
	"github.com/google/uuid"
)

type Chunk struct {
	ChunkID      uuid.UUID `db:"chunk_id"`
	UserID       string    `db:"user_id"`
	DocumentID   uuid.UUID `db:"document_id"`
	ChunkContent string    `db:"chunk_content"`
	ChunkIndex   int       `db:"chunk_index"`
}
