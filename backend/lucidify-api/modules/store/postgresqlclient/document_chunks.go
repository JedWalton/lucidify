package postgresqlclient

import (
	"lucidify-api/modules/store/storemodels"

	"github.com/google/uuid"
)

// UploadChunk uploads a chunk to the database.
// DeleteAllChunksByDocumentID deletes all chunks for a document.
// DeleteAllChunksByUserID deletes all chunks for a user.

func (s *PostgreSQL) UploadChunks(chunks []storemodels.Chunk) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO document_chunks (document_id, chunk_content, chunk_index)
	          VALUES ($1, $2, $3)`

	for _, chunk := range chunks {
		_, err := tx.Exec(query, chunk.DocumentID, chunk.ChunkContent, chunk.ChunkIndex)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *PostgreSQL) DeleteAllChunksByDocumentID(documentID uuid.UUID) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `DELETE FROM document_chunks WHERE document_id = $1`
	_, err = tx.Exec(query, documentID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
