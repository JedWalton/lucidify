package store

type DocumentStore interface {
	UploadDocument(title, content string) (int64, error)
}

func (s *Store) UploadDocument(name, content string) error {
	query := `INSERT INTO documents (document_name, content) VALUES ($1, $2)`
	_, err := s.db.Exec(query, name, content)
	return err
}

func (s *Store) GetDocument(name string) (string, error) {
	var content string
	query := `SELECT content FROM documents WHERE document_name = $1`
	err := s.db.QueryRow(query, name).Scan(&content)
	if err != nil {
		return "", err
	}
	return content, nil
}

func (s *Store) DeleteDocument(name string) error {
	query := `DELETE FROM documents WHERE document_name = $1`
	_, err := s.db.Exec(query, name)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) UpdateDocument(name, content string) error {
	query := `UPDATE documents SET content = $1 WHERE document_name = $2`
	_, err := s.db.Exec(query, content, name)
	if err != nil {
		return err
	}
	return nil
}
