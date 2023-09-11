package store

type Storer interface {
	Save(text, source string) (int64, error)
	DeleteByID(id int64) error
}

func (store *Store) Save(text, source string) (int64, error) {
	result, err := store.db.Exec("INSERT INTO table_name (text, source) VALUES ($1, $2)", text, source)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (store *Store) DeleteByID(id int64) error {
	_, err := store.db.Exec("DELETE FROM table_name WHERE id = $1", id)
	return err
}
