package storemodels

import "time"

type Role string

const (
	RoleAssistant Role = "assistant"
	RoleUser      Role = "user"
)

type OpenAIModel string // Define this based on the possible values for an OpenAI model

type Message struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
}

type Conversation struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Messages    []Message   `json:"messages"`
	Model       OpenAIModel `json:"model"`
	Prompt      string      `json:"prompt"`
	Temperature float64     `json:"temperature"`
	FolderID    string      `json:"folderId,omitempty"` // omitempty because it can be null
	CreatedAt   time.Time   `json:"createdAt"`
	UpdatedAt   time.Time   `json:"updatedAt"`
}

type FolderType string

const (
	FolderTypeChat   FolderType = "chat"
	FolderTypePrompt FolderType = "prompt"
)

type Folder struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Type      FolderType `json:"type"`
	CreatedAt time.Time  `json:"createdAt"`
}

type Prompt struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Content     string      `json:"content"`
	Model       OpenAIModel `json:"model"`
	FolderID    string      `json:"folderId,omitempty"` // omitempty because it can be null
	CreatedAt   time.Time   `json:"createdAt"`
}

type KeyValuePair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type PluginID string // Define this based on the possible values for a PluginID

type PluginName string // Define this based on the possible values for a PluginName

type Plugin struct {
	ID           PluginID       `json:"id"`
	Name         PluginName     `json:"name"`
	RequiredKeys []KeyValuePair `json:"requiredKeys"`
}

type PluginKey struct {
	PluginID     PluginID       `json:"pluginId"`
	RequiredKeys []KeyValuePair `json:"requiredKeys"`
}

type Settings struct {
	Theme string `json:"theme"` // This matches the TypeScript 'light' | 'dark'
}

// LocalStorageData represents a user's local storage data on the server.
type LocalStorageData struct {
	UserID               string         `json:"userId"` // It's essential to associate each LocalStorage with a user
	ApiKey               string         `json:"apiKey"`
	ConversationHistory  []Conversation `json:"conversationHistory"`
	SelectedConversation Conversation   `json:"selectedConversation"`
	Theme                string         `json:"theme"`
	Folders              []Folder       `json:"folders"`
	Prompts              []Prompt       `json:"prompts"`
	ShowChatbar          bool           `json:"showChatbar"`
	ShowPromptbar        bool           `json:"showPromptbar"`
	PluginKeys           []PluginKey    `json:"pluginKeys"`
	Settings             Settings       `json:"settings"`
}
