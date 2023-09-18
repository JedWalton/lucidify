package store

type Document struct {
	user_id       int
	document_name string
	content       string
}

func (s *Store) UploadDocument(userID int, name, content string) error {
	query := `INSERT INTO documents (user_id, document_name, content) VALUES ($1, $2, $3)`
	_, err := s.db.Exec(query, userID, name, content)
	return err
}

func (s *Store) GetDocument(userID int, name string) (string, error) {
	var content string
	query := `SELECT content FROM documents WHERE user_id = $1 AND document_name = $2`
	err := s.db.QueryRow(query, userID, name).Scan(&content)
	if err != nil {
		return "", err
	}
	return content, nil
}

func (s *Store) GetAllDocuments(userID int) ([]Document, error) {
	var documents []Document
	query := `SELECT user_id, document_name, content FROM documents WHERE user_id = $1`
	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var doc Document
		err := rows.Scan(&doc.user_id, &doc.document_name, &doc.content)
		if err != nil {
			return nil, err
		}
		documents = append(documents, doc)
	}
	return documents, nil
}

func (s *Store) DeleteDocument(userID int, name string) error {
	query := `DELETE FROM documents WHERE user_id = $1 AND document_name = $2`
	_, err := s.db.Exec(query, userID, name)
	return err
}

func (s *Store) UpdateDocument(userID int, name, content string) error {
	query := `UPDATE documents SET content = $1 WHERE user_id = $2 AND document_name = $3`
	_, err := s.db.Exec(query, content, userID, name)
	return err
}
