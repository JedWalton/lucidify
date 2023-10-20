package chatservice

// ChatService defines the interface for chat operations.
type ChatService interface {
	StartChat(userID, initialMessage, model, temperature, prompt, folderID string) error
	SendMessage(userID, chatID, role, content string) error
	// ... other chat operations methods ...
}

// ChatServiceImpl is a concrete implementation of ChatService.
type ChatServiceImpl struct {
	chs ChatHistoryService
	cvs ChatVectorService
}

// NewChatThreadService creates a new ChatService with its dependencies.
func NewChatService(chs ChatHistoryService, cvs ChatVectorService) ChatService {
	return &ChatServiceImpl{
		chs: chs,
		cvs: cvs,
	}
}

// StartChat initiates a new chat session.
func (cs *ChatServiceImpl) StartChat(userID, initialMessage, model, temperature, prompt, folderID string) error {
	// Start a new chat, which may involve creating a new chat history
	chatID, err := cs.chs.CreateNewChatHistory(userID, model, temperature, prompt, folderID)
	if err != nil {
		return err
	}
	return cs.SendMessage(chatID, "user", initialMessage, userID)
}

// SendMessage sends a message to the chat identified by chatID.
func (cs *ChatServiceImpl) SendMessage(userID, chatID, role, content string) error {
	// Logic to send message
	// If role is "user", communicate with OpenAI API and get a response
	if role == "user" {
		// Save user's message
		err := cs.chs.AddMessageToHistory(chatID, role, content)
		if err != nil {
			return err
		}

		response, err := cs.cvs.GetAnswerFromFiles(content, userID)
		if err != nil {
			return err
		}

		err = cs.chs.AddMessageToHistory(chatID, "assistant", response)
		if err != nil {
			return err
		}
	}
	// If the role is not "user", we might save the message differently, handle this according to your logic.
	return nil
}
