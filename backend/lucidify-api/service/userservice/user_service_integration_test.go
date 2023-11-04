package userservice

import (
	"lucidify-api/data/store/postgresqlclient"
	"lucidify-api/data/store/storemodels"
	"lucidify-api/data/store/weaviateclient"
	"lucidify-api/service/documentservice"
	"testing"
	"time"
)

// func setupTests() (UserService, storemodels.User, error, *postgresqlclient.PostgreSQL) {
// 	user := storemodels.User{
// 		UserID:           "TestCreateUserTableUserServiceUserID",
// 		ExternalID:       "TestCreateUserTableExternalID",
// 		Username:         "TestCreateUserableUsername",
// 		PasswordEnabled:  true,
// 		Email:            "TestCreateUser@example.com",
// 		FirstName:        "TestCreateUserCreateTest",
// 		LastName:         "TestCreateUserUser",
// 		ImageURL:         "https://TestCreateUser.com/image.jpg",
// 		ProfileImageURL:  "https://TestCreateUser.com/profile.jpg",
// 		TwoFactorEnabled: false,
// 		CreatedAt:        1654012591514,
// 		UpdatedAt:        1654012591514,
// 	}
//
// 	db, err := postgresqlclient.NewPostgreSQL()
// 	if err != nil {
// 		return nil, user, err, db
// 	}
//
// 	weaviate, err := weaviateclient.NewWeaviateClientTest()
// 	if err != nil {
// 		log.Fatalf("Failed to create WeaviateClient: %v", err)
// 	}
// 	userService, err := NewUserService(weaviate)
// 	if err != nil {
// 		return nil, user, err, db
// 	}
// 	return userService, user, nil, db
// }
//
// func cleanupTests(user storemodels.User, db *postgresqlclient.PostgreSQL) error {
// 	return db.DeleteUserInUsersTable(user.UserID)
// }

//	func TestCreateUser(t *testing.T) {
//		userService, user, err, db := setupTests()
//		if err != nil {
//			t.Error(err)
//		}
//		err = userService.CreateUser(user)
//		if err != nil {
//			t.Error(err)
//		}
//		userFromDb, err := db.GetUserInUsersTable(user.UserID)
//		if err != nil {
//			t.Error(err)
//		}
//		if userFromDb.UserID != user.UserID {
//			t.Errorf("Expected %s, got %s", user.UserID, userFromDb.UserID)
//		}
//
//		t.Cleanup(func() {
//			cleanupTests(user, db)
//		})
//	}
//
//	func TestUpdateUser(t *testing.T) {
//		userService, user, err, db := setupTests()
//		if err != nil {
//			t.Error(err)
//		}
//		err = userService.CreateUser(user)
//		if err != nil {
//			t.Error(err)
//		}
//		userUpdated := storemodels.User{
//			UserID:           "TestCreateUserTableUserServiceUserID",
//			ExternalID:       "TestCreateUserTableExternalIDUpdated",
//			Username:         "TestCreateUserableUsernameUpdated",
//			PasswordEnabled:  true,
//			Email:            "TestCreateUserUpdated@example.com",
//			FirstName:        "TestCreateUserCreateTestUpdated",
//			LastName:         "TestCreateUserUserUpdated",
//			ImageURL:         "https://TestCreateUserUpdated.com/image.jpg",
//			ProfileImageURL:  "https://TestCreateUserUpdated.com/profile.jpg",
//			TwoFactorEnabled: false,
//			CreatedAt:        1654012591514,
//			UpdatedAt:        1654012591514,
//		}
//		err = userService.UpdateUser(userUpdated)
//		if err != nil {
//			t.Error(err)
//		}
//		userAfterUpdate, err := db.GetUserInUsersTable(user.UserID)
//		if err != nil {
//			t.Error(err)
//		}
//		if userAfterUpdate.UserID != user.UserID {
//			t.Errorf("Expected %s, got %s", user.UserID, userAfterUpdate.UserID)
//		}
//		if userAfterUpdate.ExternalID != userUpdated.ExternalID {
//			t.Errorf("Expected %s, got %s", userUpdated.ExternalID, userAfterUpdate.ExternalID)
//		}
//		if userAfterUpdate.Username != userUpdated.Username {
//			t.Errorf("Expected %s, got %s", userUpdated.Username, userAfterUpdate.Username)
//		}
//		if userAfterUpdate.PasswordEnabled != userUpdated.PasswordEnabled {
//			t.Errorf("Expected %t, got %t", userUpdated.PasswordEnabled, userAfterUpdate.PasswordEnabled)
//		}
//		if userAfterUpdate.Email != userUpdated.Email {
//			t.Errorf("Expected %s, got %s", userUpdated.Email, userAfterUpdate.Email)
//		}
//		if userAfterUpdate.FirstName != userUpdated.FirstName {
//			t.Errorf("Expected %s, got %s", userUpdated.FirstName, userAfterUpdate.FirstName)
//		}
//		if userAfterUpdate.LastName != userUpdated.LastName {
//			t.Errorf("Expected %s, got %s", userUpdated.LastName, userAfterUpdate.LastName)
//		}
//		if userAfterUpdate.ImageURL != userUpdated.ImageURL {
//			t.Errorf("Expected %s, got %s", userUpdated.ImageURL, userAfterUpdate.ImageURL)
//		}
//		if userAfterUpdate.ProfileImageURL != userUpdated.ProfileImageURL {
//			t.Errorf("Expected %s, got %s", userUpdated.ProfileImageURL, userAfterUpdate.ProfileImageURL)
//		}
//		if userAfterUpdate.TwoFactorEnabled != userUpdated.TwoFactorEnabled {
//			t.Errorf("Expected %t, got %t", userUpdated.TwoFactorEnabled, userAfterUpdate.TwoFactorEnabled)
//		}
//		if userAfterUpdate.CreatedAt != userUpdated.CreatedAt {
//			t.Errorf("Expected %d, got %d", userUpdated.CreatedAt, userAfterUpdate.CreatedAt)
//		}
//		if userAfterUpdate.UpdatedAt != userUpdated.UpdatedAt {
//			t.Errorf("Expected %d, got %d", userUpdated.UpdatedAt, userAfterUpdate.UpdatedAt)
//		}
//
//		t.Cleanup(func() {
//			cleanupTests(user, db)
//		})
//	}
func TestDeleteUser(t *testing.T) {
	user := storemodels.User{
		UserID:           "TestDeleteUserAndAssociatedDocumentsUserID",
		ExternalID:       "TestCreateUserTableExternalID",
		Username:         "TestCreateUserableUsername",
		PasswordEnabled:  true,
		Email:            "TestCreateUser@example.com",
		FirstName:        "TestCreateUserCreateTest",
		LastName:         "TestCreateUserUser",
		ImageURL:         "https://TestCreateUser.com/image.jpg",
		ProfileImageURL:  "https://TestCreateUser.com/profile.jpg",
		TwoFactorEnabled: false,
		CreatedAt:        1654012591514,
		UpdatedAt:        1654012591514,
	}
	// userService, _, err := setupTests()
	// if err != nil {
	// 	t.Error(err)
	// }
	weaviateClient, err := weaviateclient.NewWeaviateClientTest()
	if err != nil {
		t.Fatalf("failed to create weaviate client: %v", err)
	}

	postgre, err := postgresqlclient.NewPostgreSQL()
	if err != nil {
		t.Errorf("Failed to create test postgresqlclient: %v", err)
	}
	// userService, err := userservice.NewUserService(weaviateClient)
	userService, err := NewUserService(postgre, weaviateClient)
	if err != nil {
		t.Errorf("Failed to create test userservice: %v", err)
	}

	err = userService.CreateUser(user)
	if err != nil {
		t.Error(err)
	}
	_, err = userService.GetUserWithRetries(user.UserID, 5)
	if err != nil {
		t.Errorf("User not found after creation: %v", err)
	}
	// _, err = weaviateclient.NewWeaviateClientTest()
	postgres, err := postgresqlclient.NewPostgreSQL()
	if err != nil {
		t.Errorf("Failed to create test postgresqlclient: %v", err)
	}
	// // Test data
	// name := "test-document-name"
	// content := "This is a test document content. Placeholder content"

	documentService := documentservice.NewDocumentService(postgres, weaviateClient)

	// // 2. Call the function
	// document, err := documentService.UploadDocument(user.UserID, name, content)
	////////////////////////// THIS LINE IS THE PROBLEM //////////////////////////
	// _, err = documentService.UploadDocument(user.UserID, name, content)
	document, err := documentService.UploadDocument("TestDeleteUserAndAssociatedDocumentsUserID", "Dog Knowledge",
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
	if err != nil {
		t.Fatalf("Failed to upload document: %v", err)
	}
	t.Logf("document: %v", document)

	time.Sleep(2 * time.Second)

	err = userService.DeleteUser(user.UserID)
	if err != nil {
		t.Error(err)
	}

	// chunks, err := db.GetChunksOfDocument(document)
	// if err != nil || len(chunks) == 0 {
	// 	t.Error("Chunks were not uploaded to PostgreSQL")
	// }
	// chunksWeaviate, err := weaviateClient.GetChunks(chunks)
	// if err != nil || len(chunksWeaviate) == 0 {
	// 	t.Error("Chunks were not uploaded to Weaviate")
	// }
	// err = weaviateClient.UploadChunks(chunks)
	// if err != nil {
	// 	t.Errorf("UploadChunks failed: %v", err)
	// }
	// err = weaviateClient.UploadChunks(chunks)
	// if err == nil {
	// 	t.Errorf("UploadChunks should have failed due to duplication: %v", err)
	// }
	// chunks, err = weaviateClient.GetChunks(chunks)
	// if err != nil || len(chunks) != 3 {
	// 	t.Errorf("GetChunks should have succeeded: %v", err)
	// }

	// time.Sleep(10 * time.Second)

	// chunksReturned, err := weaviateClient.GetChunks(chunks)
	// if err == nil || len(chunksReturned) != 1 {
	// 	t.Errorf("Get Chunks should return 1 chunk. returned chunks: %v", len(chunks))
	// }

	// if !userService.HasUserBeenDeleted(user.UserID, 5) {
	// 	t.Errorf("User not deleted after deletion: %v", err)
	// }

	// t.Cleanup(func() {
	// 	cleanupTests(user, db)
	// })
}

// func TestGetUser(t *testing.T) {
// 	userService, user, err, db := setupTests()
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	err = userService.CreateUser(user)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	_, err = userService.GetUserWithRetries(user.UserID, 5)
// 	if err != nil {
// 		t.Errorf("User not found after creation: %v", err)
// 	}
// 	userFromDb, err := userService.GetUser(user.UserID)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if userFromDb.Email != user.Email {
// 		t.Errorf("Expected %s, got %s", user.Email, userFromDb.Email)
// 	}
//
// 	t.Cleanup(func() {
// 		cleanupTests(user, db)
// 	})
// }
//
// func TestGetUserWithRetries(t *testing.T) {
// 	userService, user, err, db := setupTests()
// 	if err != nil {
// 		t.Error(err)
// 	}
//
// 	err = userService.CreateUser(user)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	// Case 1: User is found without needing retries
// 	foundUser, err := userService.GetUserWithRetries(user.UserID, 5)
// 	if err != nil || foundUser == nil {
// 		t.Error("GetUserWithRetries failed to retrieve the user, but the user should have been found")
// 	}
// 	// Case 2: User is not found even after retries
// 	nonExistentUserID := "non-existent-user-id"
// 	// Delay the start of the next case to avoid immediate execution after user creation
// 	time.Sleep(2 * time.Second)
//
// 	_, err = userService.GetUserWithRetries(nonExistentUserID, 5)
// 	if err == nil {
// 		t.Error("GetUserWithRetries did not return an error, but it should have since the user does not exist")
// 	}
//
// 	t.Cleanup(func() {
// 		cleanupTests(user, db)
// 	})
// }
//
// func TestHasUserBeenDeleted(t *testing.T) {
// 	userService, user, err, db := setupTests()
// 	if err != nil {
// 		t.Error(err)
// 	}
//
// 	err = userService.CreateUser(user)
// 	if err != nil {
// 		t.Error(err)
// 	}
//
// 	// Case 1: User has not been deleted
// 	userNotDeleted := userService.HasUserBeenDeleted(user.UserID, 5)
// 	if userNotDeleted {
// 		t.Error("HasUserBeenDeleted returns true, but the user was not deleted")
// 	}
//
// 	// Assume we have a function to delete a user for test
// 	userService.DeleteUser(user.UserID)
//
// 	// Case 2: User has been deleted
// 	userHasBeenDeleted := userService.HasUserBeenDeleted(user.UserID, 5)
// 	if !userHasBeenDeleted {
// 		t.Error("HasUserBeenDeleted returns false, but the user was deleted")
// 	}
//
// 	t.Cleanup(func() {
// 		cleanupTests(user, db)
// 	})
// }
