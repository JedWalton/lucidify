package syncservice

type Role string

const (
	Assistant Role = "assistant"
	User      Role = "user"
)

type OpenAIModel struct {
	// Include fields from the 'OpenAIModel' type
}

type Message struct {
	Role    Role
	Content string
}

type Conversation struct {
	ID          string
	Name        string
	Messages    []Message
	Model       OpenAIModel
	Prompt      string
	Temperature float64
	FolderID    *string
}

type FolderType string

const (
	Chat       FolderType = "chat"
	PromptType FolderType = "prompt"
)

type FolderInterface struct {
	ID   string
	Name string
	Type FolderType
}

type KeyValuePair struct {
	Key   string
	Value string
}

type PluginID string

const (
	GoogleSearch PluginID = "google-search"
)

type PluginName string

const (
	GoogleSearchName PluginName = "Google Search"
)

type Plugin struct {
	ID           PluginID
	Name         PluginName
	RequiredKeys []KeyValuePair
}

type PluginKey struct {
	PluginID     PluginID
	RequiredKeys []KeyValuePair
}

var Plugins = map[PluginID]Plugin{
	GoogleSearch: {
		ID:   GoogleSearch,
		Name: GoogleSearchName,
		RequiredKeys: []KeyValuePair{
			{Key: "GOOGLE_API_KEY", Value: ""},
			{Key: "GOOGLE_CSE_ID", Value: ""},
		},
	},
}

type PluginList []Plugin // No direct equivalent in Go but this can be useful

type Prompt struct {
	ID          string
	Name        string
	Description string
	Content     string
	Model       OpenAIModel
	FolderID    *string
}

type Settings struct {
	Theme string // Consider using a custom type for strict values like in TypeScript
}

type LocalStorage struct {
	APIKey               string
	ConversationHistory  []Conversation
	SelectedConversation Conversation
	Theme                string
	Folders              []FolderInterface
	Prompts              []Prompt
	ShowChatbar          bool
	ShowPromptbar        bool
	PluginKeys           []PluginKey
	Settings             Settings
}
