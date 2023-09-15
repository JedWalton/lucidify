package store

func (store *Store) UploadDocument(dataSchema map[string]string) (int64, error) {
	title := dataSchema["title"]
	content := dataSchema["content"]
	result, err := store.db.Exec("INSERT INTO documents(title, content) VALUES ($1, $2)", title, content)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
