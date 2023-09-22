// // go:build integration
// // +build integration
package clerk

// single responsibility principle

// func TestIntegration_usercreatedevent(t *testing.T) {
// 	testconfig := config.NewTestServerConfig()
// 	clerkSecretKey := testconfig.ClerkSecretKey
// 	userID, err := store.CreateUserInClerk(clerkSecretKey, "clerk_handler_uce_int_firstname", "clerk_handler_uce_int_firstname", "clerk_handler_uce_integration@example.com", "$sswordoatnsu28348ckj")
// 	if err != nil {
// 		t.Fatalf("Failed to create user in Clerk: %v", err)
// 	}
//
// 	// Cleanup
// 	t.Cleanup(func() {
// 		if userID != "" {
// 			err = store.DeleteUserInClerk(clerkSecretKey, userID)
// 			if err != nil {
// 				log.Printf("Did not delete test user in clerk: %v\n", err)
// 			}
// 		}
// 	})
// }

// func TestIntegration_userupdateevent(t *testing.T) {
// 	testconfig := config.NewTestServerConfig()
// 	clerkSecretKey := testconfig.ClerkSecretKey
// 	firstName := "TestFirstName"
// 	lastName := "TestLastName"
// 	email := "testUpdateUserInClerk@example.com"
// 	password := "soaenuth4yg8fdbioea"
//
// 	userID, err := store.CreateUserInClerk(clerkSecretKey, firstName, lastName, email, password)
// 	if err != nil {
// 		t.Fatalf("Failed to create user in Clerk: %v", err)
// 	}
//
// 	newFirstName := "UpdatedFirstName"
// 	newLastName := "UpdatedLastName"
// 	err = store.UpdateUserInClerk(clerkSecretKey, userID, newFirstName, newLastName)
// 	if err != nil {
// 		t.Fatalf("Failed to update user in Clerk: %v", err)
// 	}
//
// 	storeInstance, err := store.NewStore(testconfig.PostgresqlURL)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	// Cleanup
// 	t.Cleanup(func() {
// 		if userID != "" {
// 			err = store.DeleteUserInClerk(clerkSecretKey, userID)
// 			if err != nil {
// 				log.Printf("Did not delete test user in clerk: %v\n", err)
// 			}
// 			err = store.DeleteUserInClerk(clerkSecretKey, userID)
// 			if err != nil {
// 				log.Printf("Did not delete test user in clerk: %v\n", err)
// 			}
// 			err = storeInstance.DeleteUserInUsersTable(userID)
// 			if err != nil {
// 				log.Printf("Did not delete test user in users table: %v\n", err)
// 			}
// 			err = storeInstance.CheckUserDeletedInUsersTable(userID, 5)
// 		}
// 	})
// }
//
// func TestIntegration_userdeletedevent(t *testing.T) {
// 	testconfig := config.NewTestServerConfig()
// 	clerkSecretKey := testconfig.ClerkSecretKey
// 	userID, err := store.CreateUserInClerk(clerkSecretKey, "userdeleteevent_firstname", "userdeleteevent_lastname", "userdeleteevent@example.com", "$sswordoatnsu28348ckj")
// 	if err != nil {
// 		t.Fatalf("Failed to create user in Clerk: %v", err)
// 	}
// 	storeInstance, err := store.NewStore(testconfig.PostgresqlURL)
// 	if err != nil {
// 		t.Fatalf("Failed to create storeInstance: %v", err)
// 	}
// 	err = storeInstance.CheckIfUserInUsersTable(userID, 5)
// 	if err == nil {
// 		t.Fatalf("successfully deleted %v", err)
// 	}
//
// 	store.DeleteUserInClerk(clerkSecretKey, userID)
//
// 	// Cleanup
// 	t.Cleanup(func() {
// 		if userID != "" {
// 			err = store.DeleteUserInClerk(clerkSecretKey, userID)
// 			if err != nil {
// 				log.Printf("Did not delete test user in clerk: %v\n", err)
// 			}
// 			err = store.DeleteUserInClerk(clerkSecretKey, userID)
// 			if err != nil {
// 				log.Printf("Did not delete test user in clerk: %v\n", err)
// 			}
// 			err = storeInstance.DeleteUserInUsersTable(userID)
// 			if err != nil {
// 				log.Printf("Did not delete test user in users table: %v\n", err)
// 			}
// 			err = storeInstance.CheckUserDeletedInUsersTable(userID, 5)
// 		}
// 	})
// }
