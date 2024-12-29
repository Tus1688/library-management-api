package main

import (
	"flag"
	"fmt"
	"github.com/Tus1688/library-management-api/storage"
	"log"
	"os"
)

func main() {
	initAdmin := flag.NewFlagSet("init-admin", flag.ExitOnError)
	username := initAdmin.String("username", "", "Admin username")
	password := initAdmin.String("password", "", "Admin password")

	if len(os.Args) < 2 {
		fmt.Println("expected 'init-admin' subcommand")
		os.Exit(1)
	}

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

	if initAdmin.Parsed() {
		if *username == "" || *password == "" {
			initAdmin.PrintDefaults()
			os.Exit(1)
		}

		postgres, err := storage.NewPostgresStore()
		if err != nil {
			log.Fatal("unable to connect to postgres: ", err)
		}

		err = postgres.InitAdmin(username, password)
		if err != nil {
			log.Fatal("unable to initialize admin: ", err)
		}

		log.Print("admin user initialized successfully")
	}
}
