package main

import (
	"flag"
	"fmt"
	"github.com/Tus1688/library-management-api/storage"
	"log"
	"os"
)

// main is the entry point of the CLI application.
// It parses command-line flags for initializing an admin user and connects to the Postgres database.
// If the 'init-admin' subcommand is provided, it initializes an admin user with the given username and password.
func main() {
	// Define the 'init-admin' subcommand and its flags
	initAdmin := flag.NewFlagSet("init-admin", flag.ExitOnError)
	username := initAdmin.String("username", "", "Admin username")
	password := initAdmin.String("password", "", "Admin password")

	// Check if a subcommand is provided
	if len(os.Args) < 2 {
		fmt.Println("expected 'init-admin' subcommand")
		os.Exit(1)
	}

	// Parse the 'init-admin' subcommand
	switch os.Args[1] {
	case "init-admin":
		err := initAdmin.Parse(os.Args[2:])
		if err != nil {
			log.Fatal("unable to parse flags: ", err)
			return
		}
	default:
		fmt.Println("expected 'init-admin' subcommand")
		os.Exit(1)
	}

	// Validate the parsed flags and initialize the admin user
	if initAdmin.Parsed() {
		if *username == "" || *password == "" {
			initAdmin.PrintDefaults()
			os.Exit(1)
		}

		// Connect to the Postgres database
		postgres, err := storage.NewPostgresStore()
		if err != nil {
			log.Fatal("unable to connect to postgres: ", err)
		}

		// Initialize the admin user with the provided username and password
		err = postgres.InitAdmin(username, password)
		if err != nil {
			log.Fatal("unable to initialize admin: ", err)
		}

		log.Print("admin user initialized successfully")
	}
}
