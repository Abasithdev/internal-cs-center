package main

import (
	"log"

	_ "abasithdev.github.io/internal-cs-center-backend/cmd/server/docs"
	"abasithdev.github.io/internal-cs-center-backend/internal/router"
)

func main() {
	log.Println("Starting server :8080")

	r := router.NewRouter()
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
