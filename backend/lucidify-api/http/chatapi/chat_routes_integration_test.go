// //go:build integration
// // +build integration
package chatapi

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"lucidify-api/data/store/postgresqlclient"
	"lucidify-api/data/store/storemodels"
	"lucidify-api/data/store/weaviateclient"
	"lucidify-api/server/config"
	"lucidify-api/service/chatservice"
	"lucidify-api/service/clerkservice"
	"lucidify-api/service/documentservice"
	"lucidify-api/service/userservice"

	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/sashabaranov/go-openai"
)

func createTestUserInDb(cfg *config.ServerConfig, db *postgresqlclient.PostgreSQL) error {
	// the user id registered by the jwt token must exist in the local database
	user := storemodels.User{
		UserID:           cfg.TestUserID,
		ExternalID:       "TestChatAPIUserInUsersTableExternalIDDocuments",
		Username:         "TestChatAPIUsersTableUsernameDocuments",
		PasswordEnabled:  true,
		Email:            "TestChatAPIUserUsersTableDocuments@example.com",
		FirstName:        "TestUsersTableCreateTest",
		LastName:         "TestUsersTableUser",
		ImageURL:         "https://TestInUsersTable.com/image.jpg",
		ProfileImageURL:  "https://TestInUsersTable.com/profile.jpg",
		TwoFactorEnabled: false,
		CreatedAt:        1654012591514,
		UpdatedAt:        1654012591514,
	}

	weaviate, err := weaviateclient.NewWeaviateClientTest()
	if err != nil {
		log.Fatalf("Failed to create WeaviateClient: %v", err)
	}
	userService, err := userservice.NewUserService(db, weaviate)
	if err != nil {
		log.Fatalf("Failed to create UserService: %v", err)
	}

	err = userService.DeleteUser(user.UserID)
	if err != nil {
		log.Fatalf("Failed to delete user: %v", err)
	}
	if !userService.HasUserBeenDeleted(user.UserID, 3) {
		log.Fatalf("Failed to delete user: %v", err)
	}

	err = userService.CreateUser(user)
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}

	// Check if the user exists
	_, err = userService.GetUserWithRetries(user.UserID, 3)
	if err != nil {
		log.Fatalf("User not found after creation: %v", err)
	}

	return nil
}

type TestSetup struct {
	Config        *config.ServerConfig
	PostgresqlDB  *postgresqlclient.PostgreSQL
	ClerkInstance clerk.Client
	Weaviate      weaviateclient.WeaviateClient
	DocService    documentservice.DocumentService
}

func SetupTestEnvironment(t *testing.T) *TestSetup {
	cfg := config.NewServerConfig()

	postgresqlDB, err := postgresqlclient.NewPostgreSQL()
	if err != nil {
		t.Fatalf("Failed to create test postgresqlclient: %v", err)
	}

	clerkInstance, err := clerkservice.NewClerkClient()
	if err != nil {
		t.Fatalf("Failed to create Clerk client: %v", err)
	}

	weaviate, err := weaviateclient.NewWeaviateClientTest()
	if err != nil {
		t.Fatalf("Failed to create WeaviateClient: %v", err)
	}

	docService := documentservice.NewDocumentService(postgresqlDB, weaviate)

	err = createTestUserInDb(cfg, postgresqlDB)
	if err != nil {
		t.Fatalf("Failed to create test user in db: %v", err)
	}

	return &TestSetup{
		Config:        cfg,
		PostgresqlDB:  postgresqlDB,
		ClerkInstance: clerkInstance,
		Weaviate:      weaviate,
		DocService:    docService,
	}
}

func TestChatHandlerIntegration(t *testing.T) {
	setup := SetupTestEnvironment(t)
	cfg := setup.Config
	clerkInstance := setup.ClerkInstance
	openaiClient := openai.NewClient(cfg.OPENAI_API_KEY)
	documentService := setup.DocService
	chatVectorService := chatservice.NewChatVectorService(setup.Weaviate, openaiClient, documentService)

	// Create a test server
	mux := http.NewServeMux()
	SetupRoutes(cfg, mux, chatVectorService, clerkInstance)
	server := httptest.NewServer(mux)
	defer server.Close()

	catDoc, err := setup.DocService.UploadDocument(cfg.TestUserID, "Cat Knowledge",
		`Cats, with their graceful movements and independent nature,
		have been revered and adored by many civilizations throughout history. In
		ancient Egypt, they were considered sacred and were even associated with
		the goddess Bastet, who was depicted as a lioness or a woman with the head
		of a lioness. Cats were believed to have protective qualities, and harming
		one was considered a grave offense. Their sleek and mysterious demeanor
		earned them a special place in the hearts of the Egyptians, a sentiment
		that has persisted to modern times.

		One of the most captivating features of a cat is its eyes. With large,
		round pupils that can expand and contract based on the amount of light, a
		cat's eyes are a marvel of evolution. This feature allows them to have
		excellent night vision, making them adept hunters even in low-light
		conditions. The reflective layer behind their retinas, known as the tapetum
		lucidum, gives their eyes a distinctive glow in the dark and further
		enhances their ability to see at night.

		Cats are known for their grooming habits. They spend a significant amount
		of time each day cleaning their fur with their rough tongues. This not only
		keeps them clean but also helps regulate their body temperature. The act of
		grooming also has a calming effect on cats, and it's not uncommon to see
		them grooming themselves or other cats as a sign of affection and bonding.
		This meticulous cleaning ritual also aids in reducing scent, making them
		stealthier hunters.

		The purring of a cat is a sound that many find soothing and comforting.
		While it's commonly associated with contentment, cats also purr when they
		are in pain, anxious, or even when they're near death. The exact mechanism
		and purpose of purring remain a subject of research and speculation. Some
		theories suggest that purring has healing properties, as the vibrations can
		stimulate the production of certain growth factors that aid in wound
		healing.

		Domestic cats, despite being pampered pets in many households, still retain
		many of their wild instincts. Their tendency to "hunt" toys, pounce on
		moving objects, or even their habit of bringing back prey to their owners
		are all remnants of their wild ancestry. These behaviors are deeply
		ingrained and serve as a reminder that beneath their cuddly exterior lies a
		skilled predator, honed by millions of years of evolution.`)

	dogDoc, err := setup.DocService.UploadDocument(cfg.TestUserID, "Dog Knowledge",
		`Introduction to Dogs: Dogs, often referred to as "man's best friend," have been
		companions to humans for thousands of years. Originating from wild wolves,
		these loyal creatures have been domesticated and bred for various roles
		throughout history, from hunting and herding to companionship. Their keen
		senses, especially their sense of smell, combined with their innate
		intelligence, make them invaluable partners in numerous tasks. Today, dogs are
		found in countless households worldwide, providing joy, comfort, and sometimes
		even protection to their human families.

		Diverse Breeds: The world of dogs is incredibly diverse, with over 340
		recognized breeds, each with its unique characteristics, temperament, and
		appearance. From the tiny Chihuahua to the majestic Great Dane, dogs come in
		all shapes and sizes. Some breeds, like the Border Collie, are known for their
		intelligence and agility, while others, such as the Saint Bernard, are
		celebrated for their strength and gentle nature. This vast diversity ensures
		that there's a perfect dog breed for almost every individual and lifestyle.

		Roles and Responsibilities: Beyond being mere pets, dogs play various roles in
		human societies. Service dogs assist individuals with disabilities, guiding the
		visually impaired or alerting those with hearing loss. Therapy dogs provide
		emotional support in hospitals, schools, and nursing homes, offering comfort to
		those in need. Working dogs, like police K9 units or search and rescue teams,
		perform critical tasks that save lives. However, with these roles comes the
		responsibility for owners to provide proper training, care, and attention to
		their canine companions.

		Health and Care: Just like humans, dogs have specific health and care needs
		that owners must address. Regular veterinary check-ups, vaccinations, and a
		balanced diet are essential for a dog's well-being. Grooming, depending on the
		breed, can range from daily brushing to occasional baths. Exercise is crucial
		for a dog's physical and mental health, with daily walks and playtime being
		beneficial. Additionally, training and socialization from a young age ensure
		that dogs are well-behaved and can interact positively with other animals and
		people.

		The Bond Between Humans and Dogs: The relationship between humans and dogs is
		profound and multifaceted. Dogs offer unconditional love, loyalty, and
		companionship, often becoming integral members of the family. In return, humans
		provide care, shelter, and affection. Numerous studies have shown that owning a
		dog can reduce stress, increase physical activity, and bring joy to their
		owners. This symbiotic relationship, built on mutual trust and respect,
		showcases the incredible bond that has existed between our two species for
		millennia.`)

	if _, err := setup.DocService.GetDocument(cfg.TestUserID, "Dog Knowledge"); err != nil {
		t.Fatalf("Failed to get dog document: %v", err)
	}
	if _, err := setup.DocService.GetDocumentByID(cfg.TestUserID, dogDoc.DocumentUUID); err != nil {
		t.Fatalf("Failed to get dog document: %v", err)
	}
	cat, err := setup.DocService.GetDocumentByID(cfg.TestUserID, catDoc.DocumentUUID)
	if cat == nil || err != nil {
		t.Fatalf("Failed to get cat document: %v", err)
	}

	// Obtain a JWT token from Clerk
	jwtToken := cfg.TestJWTSessionToken

	// Construct a message
	messages := []Message{
		{Role: RoleUser, Content: "Hello, what are dogs?"},
	}

	// Send a POST request to the server with the JWT token and message
	body, _ := json.Marshal(map[string][]Message{"messages": messages})

	// Authenticated request
	req, _ := http.NewRequest(
		http.MethodPost,
		server.URL+"/api/chat/vector-search",
		bytes.NewBuffer(body))

	req.Header.Set("Authorization", "Bearer "+jwtToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}

	var serverResp ChatResponse
	err = json.Unmarshal(respBody, &serverResp)
	if err != nil {
		t.Errorf("Failed to unmarshal response body: %v", err)
		return
	}

	// Now you can use serverResp to make assertions or further logic
	if serverResp.Status == "fail" {
		t.Errorf("Server responded with failure: %s", serverResp.Message)
	}
	if serverResp.Status != "success" {
		t.Errorf("Expected successful response, got %+v", serverResp.Status)
	}
	if serverResp.Message == "" {
		t.Errorf("Expected message in response, got %+v", serverResp.Message)
	}

	// t.Fatalf("serverResp.Data: %+v", serverResp.Data)
	data, ok := serverResp.Data.(string)
	if !ok {
		t.Fatalf("Expected serverResp.Data to be a string, got %T", serverResp.Data)
	}

	if !strings.Contains(data, `Dogs, often referred to as "man's best friend,"`) {
		t.Errorf("Response data does not contain expected dog knowledge: %v", data)
	}

	if strings.Contains(data, "The purring of a cat is a sound that many find soothing") {
		t.Errorf("Response data should not contain cat knowledge")
	}

	t.Cleanup(func() {
		userService, err := userservice.NewUserService(setup.PostgresqlDB, setup.Weaviate)
		if err != nil {
			log.Fatalf("Failed to create UserService: %v", err)
		}
		userService.DeleteUser(cfg.TestUserID)
	})
}
