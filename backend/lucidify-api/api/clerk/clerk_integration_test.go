// // go:build integration
// // +build integration
package clerk

// func TestIntegration_usercreatedevent(t *testing.T) {
// 	testconfig := config.NewTestServerConfig()
// 	db := testconfig.TestStore
//
// 	clerkSecretKey := testconfig.ClerkSecretKey
// 	userID, err := store.CreateUserInClerk(clerkSecretKey)
// 	if err != nil {
// 		t.Fatalf("Failed to create user in Clerk: %v", err)
// 	}
// 	log.Printf("Created user in Clerk with userID: %s\n", userID)
//
// 	err = checkUserInDB(db, userID, 5) // Try 5 times
// 	if err != nil {
// 		t.Fatalf("Should have fetched user, user with ID, userID: %v", err)
// 	}
//
// 	err = deleteUserInClerk(clerkSecretKey, userID)
// }
//
// func TestIntegration_usercreatedevent_unauthenticated(t *testing.T) {
// 	testconfig := config.NewTestServerConfig()
// 	db := testconfig.TestStore
//
// 	// clerkSecretKey := testconfig.ClerkSecretKey
// 	// userID, err := createUserInClerk(clerkSecretKey)
// 	// if err != nil {
// 	// 	t.Fatalf("Failed to create user in Clerk: %v", err)
// 	// }
// 	// log.Printf("Created user in Clerk with userID: %s\n", userID)
//
// 	MakeCurlRequest := func() (int, string, error) {
// 		cmd := exec.Command("curl", "-s", "-o", "/dev/null", "-w", "%{http_code}", "-X", "POST", "http://localhost:8080/clerk/webhook", "-H", "Content-Type: application/json", "-d", "@test/example_user_created_event.txt")
// 		out, err := cmd.Output()
// 		if err != nil {
// 			return 0, "", err
// 		}
// 		statusCode, _ := strconv.Atoi(string(out))
// 		return statusCode, string(out), nil
// 	}
//
// 	statusCode, response, err := MakeCurlRequest()
// 	if err != nil {
// 		t.Fatalf("Error making curl request: %v", err)
// 	}
//
// 	// Check if the status code indicates success (e.g., 200 OK)
// 	if statusCode >= 200 && statusCode < 300 {
// 		t.Fatalf("Expected the request to fail, but got a %d status code.", statusCode)
// 	}
//
// 	log.Printf("Response: %s\n", response)
//
// 	content, err := ioutil.ReadFile("test/example_user_created_event.txt")
// 	if err != nil {
// 		fmt.Println("Error reading file:", err)
// 		return
// 	}
//
// 	// Parse the content
// 	var event ClerkEvent
// 	err = json.Unmarshal(content, &event)
// 	if err != nil {
// 		fmt.Println("Error parsing JSON:", err)
// 		return
// 	}
//
// 	// Now you can use the event variable
// 	UserID := event.Data["id"].(string)
//
// 	err = checkUserInDB(db, UserID, 3) // Try 5 times
// 	if err == nil {
// 		db.DeleteUser(UserID)
// 		t.Fatalf("Should have failed to fetch user, unauthenticated user not in db with UserID: %v", err)
// 	}
// 	db.DeleteUser(UserID)
// }
//
// func TestIntegration_UpdateUser(t *testing.T) {
// 	testconfig := config.NewTestServerConfig()
// 	db := testconfig.TestStore
//
// 	clerkSecretKey := testconfig.ClerkSecretKey
// 	userID, err := createUserInClerk(clerkSecretKey)
// 	if err != nil {
// 		t.Fatalf("Failed to create user in Clerk: %v", err)
// 	}
// 	log.Printf("Created user in Clerk with userID: %s\n", userID)
//
// 	err = checkUserInDB(db, userID, 5) // Try 5 times
// 	if err != nil {
// 		t.Fatalf("Should have fetched user, user with ID, userID: %v", err)
// 	}
//
// 	// Update the user in Clerk
// 	err = updateUserInClerk(clerkSecretKey, userID, "UpdatedFirstName", "UpdatedLastName")
// 	if err != nil {
// 		t.Fatalf("Failed to update user in Clerk: %v", err)
// 	}
//
// 	err = checkUpdatedUserInDB(db, userID, 5, "UpdatedFirstName", "UpdatedLastName")
// 	if err != nil {
// 		t.Fatalf("Failed to fetch updated user from local database: %v", err)
// 	}
//
// 	// Check if the user was updated in the local database
// 	updatedUser, err := db.GetUser(userID)
// 	if err != nil {
// 		t.Fatalf("Failed to fetch updated user from local database: %v", err)
// 	}
//
// 	if updatedUser.FirstName != "UpdatedFirstName" || updatedUser.LastName != "UpdatedLastName" {
// 		t.Fatalf("User update in local database failed. Expected first name: UpdatedFirstName, got: %v. Expected last name: UpdatedLastName, got: %v", updatedUser.FirstName, updatedUser.LastName)
// 	}
//
// 	// Cleanup: Delete the user from Clerk
// 	err = deleteUserInClerk(clerkSecretKey, userID)
// 	if err != nil {
// 		t.Fatalf("Failed to delete user from Clerk: %v", err)
// 	}
// }
//
// func TestIntegration_DeleteUser(t *testing.T) {
// 	testconfig := config.NewTestServerConfig()
// 	db := testconfig.TestStore
//
// 	clerkSecretKey := testconfig.ClerkSecretKey
// 	userID, err := createUserInClerk(clerkSecretKey)
// 	if err != nil {
// 		t.Fatalf("Failed to create user in Clerk: %v", err)
// 	}
// 	log.Printf("Created user in Clerk with userID: %s\n", userID)
//
// 	err = checkUserInDB(db, userID, 5) // Try 5 times
// 	if err != nil {
// 		t.Fatalf("Should have fetched user, user with ID, userID: %v", err)
// 	}
//
// 	// Cleanup: Delete the user from Clerk
// 	err = deleteUserInClerk(clerkSecretKey, userID)
// 	if err != nil {
// 		t.Fatalf("Failed to delete user from Clerk: %v", err)
// 	}
//
// 	err = checkUserDeletedFromDB(db, userID, 5)
// 	if err != nil {
// 		t.Fatalf("Failed to delete user from local database: %v", err)
// 	}
// }
