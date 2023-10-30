package syncservice

import "lucidify-api/data/store/postgresqlclient"

// ServerResponse is the structure that defines the standard response from the server.
type ServerResponse struct {
	Success bool        `json:"success"`           // Indicates if the operation was successful
	Data    interface{} `json:"data,omitempty"`    // Holds the actual data, if any
	Message string      `json:"message,omitempty"` // Descriptive message, especially useful in case of errors
}

type SyncService interface {
	HandleSet(key, value string) ServerResponse
	HandleGet(key string) ServerResponse
	HandleRemove(key string) ServerResponse
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

func (s *SyncServiceImpl) HandleSet(key string, value string) ServerResponse {
	// Set in postgres for user_id
	store[key] = value
	return ServerResponse{Success: true, Message: "Data set successfully for key: " + key}
}

func (s *SyncServiceImpl) HandleGet(key string) ServerResponse {
	// Get in postgres for user_id
	if data, ok := store[key]; ok {
		return ServerResponse{Success: true, Data: data, Message: "Data fetched successfully"}
	}
	return ServerResponse{Success: false, Message: "Data not found for key: " + key}
}

func (s *SyncServiceImpl) HandleRemove(key string) ServerResponse {
	// Remove in postgres for user_id
	if _, ok := store[key]; ok {
		delete(store, key)
		return ServerResponse{Success: true, Message: "Data deleted successfully for key: " + key}
	}
	return ServerResponse{Success: false, Message: "No data to delete for key: " + key}
}
