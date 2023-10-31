package postgresqlclient

import (
	"errors"
	"log"
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
		log.Println("Error fetching data:", err)
		return "", err
	}
	log.Println("Data fetched successfully for key:", key)
	log.Println("Data:", data)
	return data, nil
}

func (s *PostgreSQL) ClearConversations(userID string) error {
	query := `DELETE FROM conversation_history WHERE user_id = $1`
	_, errDelConversationsHistory := s.db.Exec(query, userID)
	if errDelConversationsHistory != nil {
		return errDelConversationsHistory
	}
	query2 := `DELETE FROM folders WHERE user_id = $1`
	_, errDelFolders := s.db.Exec(query2, userID)
	if errDelFolders != nil {
		return errDelFolders
	}
	return nil
}
