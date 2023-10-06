package postgresqlclient

import (
	"lucidify-api/modules/store/storemodels"

	"github.com/google/uuid"
)

// UploadChunk uploads a chunk to the database.
// DeleteAllChunksByDocumentID deletes all chunks for a document.
// DeleteAllChunksByUserID deletes all chunks for a user.

//	func (s *PostgreSQL) UploadChunks(chunks []storemodels.Chunk) error {
//		tx, err := s.db.Begin()
//		if err != nil {
//			return err
//		}
//		defer tx.Rollback()
//
//		query := `INSERT INTO document_chunks (document_id, chunk_content, chunk_index)
//		          VALUES ($1, $2, $3)`
//
//		for _, chunk := range chunks {
//			_, err := tx.Exec(query, chunk.DocumentID, chunk.ChunkContent, chunk.ChunkIndex)
//			if err != nil {
//				return err
//			}
//		}
//
//		return tx.Commit()
//	}
func (s *PostgreSQL) UploadChunks(chunks []storemodels.Chunk) ([]storemodels.Chunk, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `INSERT INTO document_chunks (document_id, chunk_content, chunk_index)
	          VALUES ($1, $2, $3) RETURNING chunk_id`

	var chunksWithIDs []storemodels.Chunk

	for _, chunk := range chunks {
		var id uuid.UUID
		err := tx.QueryRow(query, chunk.DocumentID, chunk.ChunkContent, chunk.ChunkIndex).Scan(&id)
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

//	func (s *PostgreSQL) GetChunksByDocumentID(userId string, documentID uuid.UUID) ([]storemodels.Chunk, error) {
//		query := `SELECT document_id, chunk_content, chunk_index FROM document_chunks WHERE document_id = $1`
//		rows, err := s.db.Query(query, documentID)
//		if err != nil {
//			return nil, err
//		}
//		defer rows.Close()
//
//		var chunks []storemodels.Chunk
//		for rows.Next() {
//			var chunk storemodels.Chunk
//			err = rows.Scan(&chunk.DocumentID, &chunk.ChunkContent, &chunk.ChunkIndex)
//			if err != nil {
//				return nil, err
//			}
//			chunks = append(chunks, chunk)
//		}
//
//		return chunks, nil
//	}
func (s *PostgreSQL) GetChunksOfDocument(document *storemodels.Document) ([]storemodels.Chunk, error) {
	// Include chunk_id in the SELECT statement
	query := `SELECT chunk_id, document_id, chunk_content, chunk_index FROM document_chunks WHERE document_id = $1`
	rows, err := s.db.Query(query, document.DocumentUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chunks []storemodels.Chunk
	for rows.Next() {
		var chunk storemodels.Chunk
		// Scan the chunk_id into the ChunkID field
		err = rows.Scan(&chunk.ChunkID, &chunk.DocumentID, &chunk.ChunkContent, &chunk.ChunkIndex)
		if err != nil {
			return nil, err
		}
		chunk.UserID = document.UserID
		chunks = append(chunks, chunk)
	}

	return chunks, nil
}
