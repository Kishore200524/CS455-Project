package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/Kishore200524/CS455-Project/backend/internal/auth"
	"github.com/Kishore200524/CS455-Project/backend/internal/config"
	"github.com/Kishore200524/CS455-Project/backend/internal/database"
	"golang.org/x/term"
)

func main() {
	if !term.IsTerminal(int(syscall.Stdin)) {
		log.Fatal("run this command from an interactive terminal so the password can be entered without echo")
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Administrator email: ")
	email, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("read administrator email: %v", err)
	}
	email = strings.TrimSpace(email)

	password, err := readPassword("Administrator password: ")
	if err != nil {
		log.Fatalf("read administrator password: %v", err)
	}
	confirmation, err := readPassword("Confirm password: ")
	if err != nil {
		log.Fatalf("read password confirmation: %v", err)
	}
	if password != confirmation {
		log.Fatal("passwords do not match")
	}

	cfg := config.Load()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	store, err := database.Connect(ctx, cfg.MongoURI, cfg.MongoDatabase)
	if err != nil {
		log.Fatalf("connect to MongoDB: %v", err)
	}
	defer func() {
		disconnectCtx, cancelDisconnect := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelDisconnect()
		if err := store.Disconnect(disconnectCtx); err != nil {
			log.Printf("disconnect from MongoDB: %v", err)
		}
	}()

	user, err := auth.NewService(store).CreateAdministrator(ctx, email, password)
	if err != nil {
		log.Fatalf("create administrator: %v", err)
	}
	fmt.Printf("Created administrator account for %s\n", user.Email)
}

func readPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	password, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return "", err
	}
	return string(password), nil
}
