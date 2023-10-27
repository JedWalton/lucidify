package syncservice

type Role string

const (
	Assistant Role = "assistant"
	User      Role = "user"
)

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

// OpenAIModel is a struct that defines the properties of an OpenAI model.
type OpenAIModel struct {
	ID         OpenAIModelID `json:"id"`
	Name       string        `json:"name"`
	MaxLength  int           `json:"maxLength"`  // Maximum length of a message
	TokenLimit int           `json:"tokenLimit"` // Token limit for the model
}

// OpenAIModelID is an enum that defines the identifiers for various OpenAI models.
type OpenAIModelID string

const (
	GPT_3_5    OpenAIModelID = "gpt-3.5-turbo"
	GPT_3_5_AZ OpenAIModelID = "gpt-35-turbo"
	GPT_4      OpenAIModelID = "gpt-4"
	GPT_4_32K  OpenAIModelID = "gpt-4-32k"
)

// FallbackModelID is a constant that provides a default OpenAIModelID in case the `DEFAULT_MODEL` environment variable is not set or set to an unsupported model.
const FallbackModelID = GPT_3_5

// OpenAIModels is a map that associates each OpenAIModelID with its corresponding OpenAIModel details.
var OpenAIModels = map[OpenAIModelID]OpenAIModel{
	GPT_3_5: {
		ID:         GPT_3_5,
		Name:       "GPT-3.5",
		MaxLength:  12000,
		TokenLimit: 4000,
	},
	GPT_3_5_AZ: {
		ID:         GPT_3_5_AZ,
		Name:       "GPT-3.5",
		MaxLength:  12000,
		TokenLimit: 4000,
	},
	GPT_4: {
		ID:         GPT_4,
		Name:       "GPT-4",
		MaxLength:  24000,
		TokenLimit: 8000,
	},
	GPT_4_32K: {
		ID:         GPT_4_32K,
		Name:       "GPT-4-32K",
		MaxLength:  96000,
		TokenLimit: 32000,
	},
}
