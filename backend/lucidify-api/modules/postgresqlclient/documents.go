package postgresqlclient

import (
	"time"
)

type Document struct {
	UserID       string
	DocumentName string
	Content      string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (s *PostgreSQL) UploadDocument(userID string, name, content string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO documents (user_id, document_name, content) VALUES ($1, $2, $3)`
	_, err = tx.Exec(query, userID, name, content)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *PostgreSQL) GetDocument(userID string, name string) (*Document, error) {
	doc := &Document{}
	query := `SELECT user_id, document_name, content, created_at, updated_at FROM documents WHERE user_id = $1 AND document_name = $2`
	err := s.db.QueryRow(query, userID, name).Scan(&doc.UserID, &doc.DocumentName, &doc.Content, &doc.CreatedAt, &doc.UpdatedAt)
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
