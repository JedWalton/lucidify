package postgresqlclient

import (
	"errors"
	"lucidify-api/modules/store/storemodels"

	"github.com/google/uuid"
)

func (s *PostgreSQL) UploadChunks(chunks []storemodels.Chunk) ([]storemodels.Chunk, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Include user_id in the INSERT statement
	query := `INSERT INTO document_chunks (user_id, document_id, chunk_content, chunk_index)
	          VALUES ($1, $2, $3, $4) RETURNING chunk_id`

	var chunksWithIDs []storemodels.Chunk

	for _, chunk := range chunks {
		var id uuid.UUID
		// Include chunk.UserID in the QueryRow function
		err := tx.QueryRow(query, chunk.UserID, chunk.DocumentID, chunk.ChunkContent, chunk.ChunkIndex).Scan(&id)
		if err != nil {
			return nil, err
		}
		chunk.ChunkID = id
		chunksWithIDs = append(chunksWithIDs, chunk)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return chunksWithIDs, nil
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

func (s *PostgreSQL) GetChunksOfDocument(document *storemodels.Document) ([]storemodels.Chunk, error) {
	if document == nil {
		return nil, errors.New("provided document is nil")
	}
	// Include user_id in the SELECT statement
	query := `SELECT chunk_id, user_id, document_id, chunk_content, chunk_index FROM document_chunks WHERE user_id = $1 AND document_id = $2`
	rows, err := s.db.Query(query, document.UserID, document.DocumentUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chunks []storemodels.Chunk
	for rows.Next() {
		var chunk storemodels.Chunk
		// Scan the user_id into the UserID field
		err = rows.Scan(&chunk.ChunkID, &chunk.UserID, &chunk.DocumentID, &chunk.ChunkContent, &chunk.ChunkIndex)
		if err != nil {
			return nil, err
		}
		chunks = append(chunks, chunk)
	}

	return chunks, nil
}
