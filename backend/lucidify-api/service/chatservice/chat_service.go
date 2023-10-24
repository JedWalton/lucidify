package chatservice

import "log"

// ChatService defines the interface for chat operations.
type ChatService interface {
	StartChat(userID, initialMessage, model, temperature, prompt, folderID string) error
	SendMessage(userID, chatID, role, content string) error
	// ... other chat operations methods ...
}

// ChatServiceImpl is a concrete implementation of ChatService.
type ChatServiceImpl struct {
	cvs ChatVectorService
}

// NewChatThreadService creates a new ChatService with its dependencies.
func NewChatService(cvs ChatVectorService) ChatService {
	return &ChatServiceImpl{
		cvs: cvs,
	}
}

// StartChat initiates a new chat session.
func (cs *ChatServiceImpl) StartChat(userID, initialMessage, model, temperature, prompt, folderID string) error {
	// Start a new chat, which may involve creating a new chat history
	return cs.SendMessage("useridhere", "user", initialMessage, userID)
}

// SendMessage sends a message to the chat identified by chatID.
func (cs *ChatServiceImpl) SendMessage(userID, chatID, role, content string) error {
	// Logic to send message
	// If role is "user", communicate with OpenAI API and get a response
	if role == "user" {
		response, err := cs.cvs.GetAnswerFromFiles(content, userID)
		if err != nil {
			return err
		}
		log.Printf("Do something with response %v", response)
	}
	// If the role is not "user", we might save the message differently, handle this according to your logic.
	return nil
}
