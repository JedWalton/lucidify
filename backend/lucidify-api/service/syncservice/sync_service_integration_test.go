// //go:build integration
// // +build integration
package syncservice

import (
	"testing"
)

func resetStorage() {
	Storage = LocalStorage{}
}

func TestIntegration(t *testing.T) {
	// Reset storage before starting tests
	resetStorage()

	// Testing HandleSet
	t.Run("HandleSet", func(t *testing.T) {
		resp := HandleSet("apiKey", "testKey123")
		if !resp.Success || resp.Message != "Data synced successfully" {
			t.Fatalf("Expected success setting APIKey but got: %v", resp)
		}

		resp = HandleSet("conversationHistory", []Conversation{})
		if !resp.Success || resp.Message != "Data synced successfully" {
			t.Fatalf("Expected success setting conversationHistory but got: %v", resp)
		}

		resp = HandleSet("selectedConversation", Conversation{})
		if !resp.Success || resp.Message != "Data synced successfully" {
			t.Fatalf("Expected success setting selectedConversation but got: %v", resp)
		}

		resp = HandleSet("theme", "dark")
		if !resp.Success || resp.Message != "Data synced successfully" {
			t.Fatalf("Expected success setting theme but got: %v", resp)
		}

		resp = HandleSet("folders", []FolderInterface{})
		if !resp.Success || resp.Message != "Data synced successfully" {
			t.Fatalf("Expected success setting folders but got: %v", resp)
		}

		resp = HandleSet("prompts", []Prompt{})
		if !resp.Success || resp.Message != "Data synced successfully" {
			t.Fatalf("Expected success setting prompts but got: %v", resp)
		}

		resp = HandleSet("showChatbar", true)
		if !resp.Success || resp.Message != "Data synced successfully" {
			t.Fatalf("Expected success setting showChatbar but got: %v", resp)
		}

		resp = HandleSet("showPromptbar", true)
		if !resp.Success || resp.Message != "Data synced successfully" {
			t.Fatalf("Expected success setting showPromptbar but got: %v", resp)
		}

		resp = HandleSet("pluginKeys", []PluginKey{})
		if !resp.Success || resp.Message != "Data synced successfully" {
			t.Fatalf("Expected success setting pluginKeys but got: %v", resp)
		}

		resp = HandleSet("settings", Settings{})
		if !resp.Success || resp.Message != "Data synced successfully" {
			t.Fatalf("Expected success setting settings but got: %v", resp)
		}
	})

	// Testing HandleGet
	t.Run("HandleGet", func(t *testing.T) {
		data, resp := HandleGet("apiKey")
		if !resp.Success || resp.Message != "Data fetched successfully" || data != "testKey123" {
			t.Fatalf("Expected to fetch APIKey correctly but got: %v, data: %v", resp, data)
		}

		data, resp = HandleGet("theme")
		if !resp.Success || resp.Message != "Data fetched successfully" || data != "dark" {
			t.Fatalf("Expected to fetch theme correctly but got: %v, data: %v", resp, data)
		}

		// Insert more HandleGet tests for other keys as needed...
	})

	// Testing HandleRemove
	t.Run("HandleRemove", func(t *testing.T) {
		resp := HandleRemove("apiKey")
		if !resp.Success || resp.Message != "Data deleted successfully" {
			t.Fatalf("Expected success deleting APIKey but got: %v", resp)
		}

		_, resp = HandleGet("apiKey")
		if resp.Success || resp.Message != "No data found for key: apiKey" {
			t.Fatalf("Expected no data for APIKey after removal but got: %v", resp)
		}

		resp = HandleRemove("theme")
		if !resp.Success || resp.Message != "Data deleted successfully" {
			t.Fatalf("Expected success deleting theme but got: %v", resp)
		}

		_, resp = HandleGet("theme")
		if resp.Success || resp.Message != "No data found for key: theme" {
			t.Fatalf("Expected no data for theme after removal but got: %v", resp)
		}

		// Insert more HandleRemove tests for other keys as needed...
	})
}
