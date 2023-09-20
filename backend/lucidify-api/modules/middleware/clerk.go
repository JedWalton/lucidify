package middleware

// var client clerk.Client
//
// func init() {
// 	apiKey := os.Getenv("CLERK_API_KEY")
// 	var err error
// 	client, err = clerk.NewClient(apiKey)
// 	if err != nil {
// 		log.Fatalf("Failed to create Clerk client: %v", err)
// 	}
// }
//
// func ClerkAuthenticationMiddleware(next http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		// Wrap the next handler with Clerk's middleware
// 		protectedHandler := clerk.WithSessionV2(client)(next)
// 		protectedHandler.ServeHTTP(w, r)
//
// 		// Retrieve the authenticated session's claims
// 		session, ok := clerk.SessionFromContext(r.Context())
// 		if !ok {
// 			http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 			return
// 		}
//
// 		// Access the "Subject" field for user ID from the jwt.Claims
// 		userID := session.Claims.Subject
// 		if userID == "" {
// 			http.Error(w, "User ID not found", http.StatusUnauthorized)
// 			return
// 		}
//
// 		ctx := context.WithValue(r.Context(), "user_id", userID)
// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	}
// }
