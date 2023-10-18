// //go:build integration
// // +build integration
package chatservice

import (
	"log"
	"lucidify-api/data/store/postgresqlclient"
	"lucidify-api/data/store/storemodels"
	"lucidify-api/data/store/weaviateclient"
	"lucidify-api/server/config"
	"lucidify-api/service/documentservice"
	"lucidify-api/service/userservice"

	"github.com/sashabaranov/go-openai"
)

func createTestUserInDb() string {
	db, err := postgresqlclient.NewPostgreSQL()

	// the user id registered by the jwt token must exist in the local database
	user := storemodels.User{
		UserID:           "TestChatServiceIntegrationTestUUID",
		ExternalID:       "TestChatServiceIntegrationTestExternalID",
		Username:         "TestChatServiceIntegrationTestUsername",
		PasswordEnabled:  true,
		Email:            "TestChatServiceIntTest@gmail.com",
		FirstName:        "TestChatServiceIntegrationTestFirstName",
		LastName:         "TestChatServiceIntegrationTestLastName",
		ImageURL:         "https://TestChatServiceIntegrationTestURL.com/image.jpg",
		ProfileImageURL:  "https://TestChatServiceTestProfileURL.com/profile.jpg",
		TwoFactorEnabled: false,
		CreatedAt:        1654012591514,
		UpdatedAt:        1654012591514,
	}

	userService, err := userservice.NewUserService()
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

	err = db.CreateUserInUsersTable(user)
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}

	_, err = userService.GetUserWithRetries(user.UserID, 3)
	if err != nil {
		log.Fatalf("User not found after creation: %v", err)
	}
	return user.UserID
}

func setupTestChatService() ChatService {
	// Initialize PostgreSQL for tests
	postgresqlDB, err := postgresqlclient.NewPostgreSQL() // Adjust this to match your actual constructor
	if err != nil {
		log.Fatalf("Failed to initialize PostgreSQL: %v", err)
	}

	// Initialize Weaviate for tests
	weaviateDB, err := weaviateclient.NewWeaviateClientTest() // Adjust this to match your actual constructor
	if err != nil {
		log.Fatalf("Failed to create Weaviate client: %v", err)
	}

	cfg := config.NewServerConfig()
	openaiClient := openai.NewClient(cfg.OPENAI_API_KEY)

	documentService := documentservice.NewDocumentService(postgresqlDB, weaviateDB)

	// Create instance of ChatService
	chatService := NewChatService(postgresqlDB, weaviateDB, openaiClient, documentService)

	createTestUserInDb()

	documentService.UploadDocument("TestChatServiceIntegrationTestUUID", "Cat Knowledge",
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

	documentService.UploadDocument("TestChatServiceIntegrationTestUUID", "Dog Knowledge",
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

	return chatService
}

// func TestChatCompletion(t *testing.T) {
// 	chatService := setupTestChatService()
//
// 	response, err := chatService.ChatCompletion("TestChatServiceIntegrationTestUUID")
//
// 	if err != nil {
// 		t.Errorf("Error was not expected while processing current thread: %v", err)
// 	}
//
// 	expectedResponse := "PLACEHOLDER RESPONSE" // Adjust "EXPECTED RESPONSE" to match what you're actually expecting.
// 	if response != expectedResponse {
// 		t.Errorf("Unexpected response: got %v want %v", response, expectedResponse)
// 	}
//
// 	// Optionally, you might want to query your databases here to assert that the expected
// 	// updates have been made as a result of calling the method.
//
// 	// Cleanup after test
// 	// Here you would clean up your database from any records you created for your test.
// }

// func TestGetAnswerFromFiles(t *testing.T) {
// 	chatService := setupTestChatService()
//
// 	response, err := chatService.GetAnswerFromFiles("Tell me about dogs", "TestChatServiceIntegrationTestUUID")
//
// 	if err != nil {
// 		t.Errorf("Error was not expected while processing current thread: %v", err)
// 	}
//
// 	expectedResponse := "PLACEHOLDER RESPONSE" // Adjust "EXPECTED RESPONSE" to match what you're actually expecting.
// 	if response != expectedResponse {
// 		t.Errorf("Unexpected response: got %v want %v", response, expectedResponse)
// 	}
//
// 	// Optionally, you might want to query your databases here to assert that the expected
// 	// updates have been made as a result of calling the method.
//
// 	// Cleanup after test
// 	// Here you would clean up your database from any records you created for your test.
// }
