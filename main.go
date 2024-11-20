package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	// Print the current working directory for debugging
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Current working directory:", cwd)

	// Start the backend service (running main.go in backend/api)
	fmt.Println("Starting backend...")
	cmdBackend := exec.Command("go", "run", "backend/api/main.go")
	cmdBackend.Stdout = os.Stdout
	cmdBackend.Stderr = os.Stderr

	// Run the backend
	if err := cmdBackend.Run(); err != nil {
		log.Fatalf("Error running backend: %v", err)
	}
}
