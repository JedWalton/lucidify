package postgresqlclient

import (
	"errors"
	"lucidify-api/modules/store/storemodels"

	"github.com/google/uuid"
)

// Extract interface here for these functions that has a PostgreSQL.
// Then make these extent the documents interface rather than PostgreSQL interface.

// type Document struct {
// 	storemodels.Document
// }

func (s *PostgreSQL) UploadDocument(userID string, name, content string) (*storemodels.Document, error) {
	doc := &storemodels.Document{}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `INSERT INTO documents (user_id, document_name, content) 
	          VALUES ($1, $2, $3) 
	          RETURNING document_id, user_id, document_name, content, created_at, updated_at`
	err = tx.QueryRow(query, userID, name, content).Scan(
		&doc.DocumentUUID, &doc.UserID, &doc.DocumentName, &doc.Content, &doc.CreatedAt, &doc.UpdatedAt)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func (s *PostgreSQL) GetDocument(userID string, name string) (*storemodels.Document, error) {
	doc := &storemodels.Document{}
	query := `SELECT document_id, user_id, document_name, content, created_at, updated_at
	          FROM documents
	          WHERE user_id = $1 AND document_name = $2`
	err := s.db.QueryRow(query, userID, name).Scan(
		&doc.DocumentUUID, &doc.UserID, &doc.DocumentName, &doc.Content, &doc.CreatedAt, &doc.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func (s *PostgreSQL) GetDocumentByUUID(documentUUID uuid.UUID) (*storemodels.Document, error) {
	doc := &storemodels.Document{}
	query := `SELECT document_id, user_id, document_name, content, created_at, updated_at
	          FROM documents
	          WHERE document_id = $1`
	err := s.db.QueryRow(query, documentUUID).Scan(
		&doc.DocumentUUID, &doc.UserID, &doc.DocumentName, &doc.Content, &doc.CreatedAt, &doc.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

//	func (s *PostgreSQL) GetAllDocuments(userID string) ([]storemodels.Document, error) {
//		var documents []storemodels.Document
//		query := `SELECT document_id, user_id, document_name, content, created_at, updated_at
//		          FROM documents WHERE user_id = $1`
//		rows, err := s.db.Query(query, userID)
//		if err != nil {
//			return nil, err
//		}
//		defer rows.Close()
//		for rows.Next() {
//			var doc storemodels.Document
//			err := rows.Scan(&doc.DocumentUUID, &doc.UserID, &doc.DocumentName, &doc.Content, &doc.CreatedAt, &doc.UpdatedAt)
//			if err != nil {
//				return nil, err
//			}
//			documents = append(documents, doc)
//		}
//		return documents, nil
//	}
func (s *PostgreSQL) GetAllDocuments(userID string) ([]storemodels.Document, error) {
	if s.db == nil {
		return nil, errors.New("database connection is nil")
	}

	var documents []storemodels.Document
	query := `SELECT document_id, user_id, document_name, content, created_at, updated_at 
	          FROM documents WHERE user_id = $1`
	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var doc storemodels.Document
		err := rows.Scan(&doc.DocumentUUID, &doc.UserID, &doc.DocumentName, &doc.Content, &doc.CreatedAt, &doc.UpdatedAt)
		if err != nil {
			return nil, err
		}
		documents = append(documents, doc)
	}

	// Check for errors from the rows.Next() loop.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Optionally handle the case where there are no documents.
	if len(documents) == 0 {
		// You can return a custom error if you want to handle this case specifically.
		// return nil, errors.New("no documents found for the given userID")
		// Or just return the empty slice without error.
	}

	return documents, nil
}

func (s *PostgreSQL) DeleteDocument(userID string, name string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `DELETE FROM documents WHERE user_id = $1 AND document_name = $2`
	_, err = tx.Exec(query, userID, name)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *PostgreSQL) DeleteDocumentByUUID(documentUUID uuid.UUID) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `DELETE FROM documents WHERE document_id = $1`
	_, err = tx.Exec(query, documentUUID.String())
	if err != nil {
		return err
	}

	return tx.Commit()
}

// func (s *PostgreSQL) UpdateDocument(userID string, name, content string) error {
// 	tx, err := s.db.Begin()
// 	if err != nil {
// 		return err
// 	}
// 	defer tx.Rollback()
//
// 	query := `UPDATE documents SET content = $1, updated_at = CURRENT_TIMESTAMP WHERE user_id = $2 AND document_name = $3`
// 	_, err = tx.Exec(query, content, userID, name)
// 	if err != nil {
// 		return err
// 	}
//
// 	return tx.Commit()
// }

func (s *PostgreSQL) UpdateDocumentContent(documentID uuid.UUID, newContent string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update the content using the document ID (UUID) in the WHERE clause
	query := `UPDATE documents SET content = $1, updated_at = CURRENT_TIMESTAMP WHERE document_id = $2`
	_, err = tx.Exec(query, newContent, documentID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *PostgreSQL) UpdateDocumentName(documentID uuid.UUID, newDocumentName string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update the document_name using the document ID (UUID) in the WHERE clause
	query := `UPDATE documents SET document_name = $1, updated_at = CURRENT_TIMESTAMP WHERE document_id = $2`
	_, err = tx.Exec(query, newDocumentName, documentID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

//
// func (s *PostgreSQL) UpdateDocumentContent(documentID uuid.UUID, newContent string) error {
// 	tx, err := s.db.Begin()
// 	if err != nil {
// 		return err
// 	}
// 	defer tx.Rollback()
//
// 	// Update the content using the document ID (UUID) in the WHERE clause
// 	query := `UPDATE documents SET content = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
// 	_, err = tx.Exec(query, newContent, documentID)
// 	if err != nil {
// 		return err
// 	}
//
// 	return tx.Commit()
// }

// Impl Delete all documents by user ID
