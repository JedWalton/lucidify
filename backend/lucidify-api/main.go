package main

import (
	"lucidify-api/server"
)

func main() {
	// This stays as production. This just means the server will use the .env file
	// in the root directory of the project.
	server.StartServer()

	// Use server.StartServer("development") for running the test suites locally.
}
