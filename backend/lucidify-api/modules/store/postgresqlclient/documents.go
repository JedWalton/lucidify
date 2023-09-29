package postgresqlclient

import (
	"time"

	"github.com/google/uuid"
)

type Document struct {
	DocumentUUID uuid.UUID `db:"id"` // db tag is used to map the struct field to the SQL column name
	UserID       string    `db:"user_id"`
	DocumentName string    `db:"document_name"`
	Content      string    `db:"content"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

//	func (s *PostgreSQL) UploadDocument(userID string, name, content string) error {
//		documentID := uuid.New()
//
//		tx, err := s.db.Begin()
//		if err != nil {
//			return err
//		}
//		defer tx.Rollback()
//
//		query := `INSERT INTO documents (id, user_id, document_name, content) VALUES ($1, $2, $3, $4)`
//		_, err = tx.Exec(query, documentID, userID, name, content)
//		if err != nil {
//			return err
//		}
//
//		return tx.Commit()
//	}
// func (s *PostgreSQL) UploadDocument(userID string, name, content string) (uuid.UUID, error) {
// 	var documentID uuid.UUID
//
// 	tx, err := s.db.Begin()
// 	if err != nil {
// 		return documentID, err
// 	}
// 	defer tx.Rollback()
//
// 	// Omit the id from the INSERT statement, and use the RETURNING clause to return the generated id
// 	query := `INSERT INTO documents (user_id, document_name, content) VALUES ($1, $2, $3) RETURNING id`
// 	err = tx.QueryRow(query, userID, name, content).Scan(&documentID)
// 	if err != nil {
// 		return documentID, err
// 	}
//
// 	err = tx.Commit()
// 	if err != nil {
// 		return documentID, err
// 	}
//
// 	return documentID, nil
// }
//
// func (s *PostgreSQL) GetDocument(userID string, name string) (*Document, error) {
// 	doc := &Document{}
// 	query := `SELECT user_id, document_name, content, created_at, updated_at FROM documents WHERE user_id = $1 AND document_name = $2`
// 	err := s.db.QueryRow(query, userID, name).Scan(&doc.UserID, &doc.DocumentName, &doc.Content, &doc.CreatedAt, &doc.UpdatedAt)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return doc, nil
// }

func (s *PostgreSQL) UploadDocument(userID string, name, content string) (*Document, error) {
	doc := &Document{}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Use the RETURNING clause to return all fields of the inserted row
	query := `INSERT INTO documents (user_id, document_name, content) 
	          VALUES ($1, $2, $3) 
	          RETURNING id, user_id, document_name, content, created_at, updated_at`
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

func (s *PostgreSQL) GetDocument(userID string, name string) (*Document, error) {
	doc := &Document{}
	query := `SELECT id, user_id, document_name, content, created_at, updated_at 
	          FROM documents 
	          WHERE user_id = $1 AND document_name = $2`
	err := s.db.QueryRow(query, userID, name).Scan(
		&doc.DocumentUUID, &doc.UserID, &doc.DocumentName, &doc.Content, &doc.CreatedAt, &doc.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func (s *PostgreSQL) GetAllDocuments(userID string) ([]Document, error) {
	var documents []Document
	query := `SELECT user_id, document_name, content, created_at, updated_at FROM documents WHERE user_id = $1`
	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var doc Document
		err := rows.Scan(&doc.UserID, &doc.DocumentName, &doc.Content, &doc.CreatedAt, &doc.UpdatedAt)
		if err != nil {
			return nil, err
		}
		documents = append(documents, doc)
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

func (s *PostgreSQL) UpdateDocument(userID string, name, content string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `UPDATE documents SET content = $1, updated_at = CURRENT_TIMESTAMP WHERE user_id = $2 AND document_name = $3`
	_, err = tx.Exec(query, content, userID, name)
	if err != nil {
		return err
	}

	return tx.Commit()
}
