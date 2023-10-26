package syncapi

import (
	"encoding/json"
	"lucidify-api/service/syncservice"
	"net/http"
)

type ChangeLog struct {
	ChangeID  *int                      `json:"changeId,omitempty"`
	Key       LocalStorageKey           `json:"key"`
	Operation string                    `json:"operation"`
	OldValue  *syncservice.LocalStorage `json:"oldValue,omitempty"`
	NewValue  *syncservice.LocalStorage `json:"newValue,omitempty"`
	Timestamp int64                     `json:"timestamp"`
}

func ChangeLogHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodPost {
			MethodNotAllowed(w)
			return
		}

		var changelog []ChangeLog
		if err := json.NewDecoder(r.Body).Decode(&changelog); err != nil {
			sendJSONResponse(w, http.StatusBadRequest, syncservice.ServerResponse{Success: false, Message: "Invalid request payload"})
			return
		}

		// Store changelog in persistence layer
		// Placeholder, based on your actual implementation of the persistence layer
		// e.g., storeChangelogInDB(changelog)

		sendJSONResponse(w, http.StatusOK, syncservice.ServerResponse{Success: true, Message: "Changelog stored successfully"})
	}
}
