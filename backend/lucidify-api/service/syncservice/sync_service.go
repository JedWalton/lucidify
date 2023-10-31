package syncservice

import (
	"log"
	"lucidify-api/data/store/postgresqlclient"
)

// ServerResponse is the structure that defines the standard response from the server.
type ServerResponse struct {
	Success bool        `json:"success"`           // Indicates if the operation was successful
	Data    interface{} `json:"data,omitempty"`    // Holds the actual data, if any
	Message string      `json:"message,omitempty"` // Descriptive message, especially useful in case of errors
}

type SyncService interface {
	HandleSet(userID, key, value string) ServerResponse
	HandleGet(userID, key string) ServerResponse
	HandleRemove(userID string) ServerResponse
}

type SyncServiceImpl struct {
	postgresqlDB *postgresqlclient.PostgreSQL
}

func NewSyncService() (SyncService, error) {
	postgresqlDB, err := postgresqlclient.NewPostgreSQL()
	if err != nil {
		return nil, err
	}
	return &SyncServiceImpl{postgresqlDB: postgresqlDB}, nil
}

var store = make(map[string]string)

func (s *SyncServiceImpl) HandleSet(userID, key, value string) ServerResponse {
	log.Println("Setting data for key:", key)
	switch key {
	case "conversationHistory":
		err := s.postgresqlDB.SetData(userID, "conversationHistory", value)
		if err != nil {
			return ServerResponse{Success: false, Message: "Error setting data for key: " + key}
		}
	case "prompts":
		err := s.postgresqlDB.SetData(userID, "prompts", value)
		if err != nil {
			return ServerResponse{Success: false, Message: "Error setting data for key: " + key}
		}
	case "folders":
		err := s.postgresqlDB.SetData(userID, "folders", value)
		if err != nil {
			return ServerResponse{Success: false, Message: "Error setting data for key: " + key}
		}
	default:
		return ServerResponse{Success: false, Message: "Invalid key"}
	}

	return ServerResponse{Success: true, Message: "Data set successfully for key: " + key}
}

func (s *SyncServiceImpl) HandleGet(userID, key string) ServerResponse {
	log.Println("Getting data for key:", key)
	switch key {
	case "conversationHistory":
		data, err := s.postgresqlDB.GetData(userID, "conversationHistory")
		if err != nil {
			return ServerResponse{Success: false, Message: "Error getting data for key: " + key}
		}
		return ServerResponse{Success: true, Data: data, Message: "Data fetched successfully"}
	case "prompts":
		data, err := s.postgresqlDB.GetData(userID, "prompts")
		if err != nil {
			return ServerResponse{Success: false, Message: "Error getting data for key: " + key}
		}
		return ServerResponse{Success: true, Data: data, Message: "Data fetched successfully"}
	case "folders":
		data, err := s.postgresqlDB.GetData(userID, "folders")
		if err != nil {
			return ServerResponse{Success: false, Message: "Error getting data for key: " + key}
		}
		return ServerResponse{Success: true, Data: data, Message: "Data fetched successfully"}
	default:
		return ServerResponse{Success: false, Message: "Invalid key"}
	}
}

func (s *SyncServiceImpl) HandleRemove(userID string) ServerResponse {
	err := s.postgresqlDB.ClearConversations(userID)
	if err != nil {
		return ServerResponse{Success: true, Message: "Conversations cleared successfully"}
	}
	return ServerResponse{Success: false, Message: "Something went wrong with clear conversations: " + err.Error()}
}
