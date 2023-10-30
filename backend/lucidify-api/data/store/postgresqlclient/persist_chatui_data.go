package postgresqlclient

import (
	"errors"
	"strings"
)

func determineTableFromKey(key string) (string, error) {
	switch {
	case strings.HasPrefix(key, "conversationHistory"):
		return "conversation_history", nil
	case strings.HasPrefix(key, "folders"):
		return "folders", nil
	case strings.HasPrefix(key, "prompts"):
		return "prompts", nil
	default:
		return "", errors.New("invalid key prefix")
	}
}

func (s *PostgreSQL) SetData(userID, key, value string) error {
	table, err := determineTableFromKey(key)
	if err != nil {
		return err
	}

	// Note: For the upsert to work, there must be a unique constraint (or primary key) on user_id
	query := `
		INSERT INTO ` + table + ` (user_id, data) 
		VALUES ($1, $2) 
		ON CONFLICT (user_id) 
		DO UPDATE SET data = EXCLUDED.data
	`

	_, err = s.db.Exec(query, userID, value)
	return err
}

func (s *PostgreSQL) GetData(userID, key string) (string, error) {
	table, err := determineTableFromKey(key)
	if err != nil {
		return "", err
	}

	var data string
	query := `SELECT data FROM ` + table + ` WHERE user_id = $1`
	err = s.db.QueryRow(query, userID).Scan(&data)
	if err != nil {
		return "", err
	}
	return data, nil
}

func (s *PostgreSQL) RemoveData(userID, key string) error {
	table, err := determineTableFromKey(key)
	if err != nil {
		return err
	}

	query := `DELETE FROM ` + table + ` WHERE id = $1 AND user_id = $2`
	_, err = s.db.Exec(query, key, userID)
	return err
}
