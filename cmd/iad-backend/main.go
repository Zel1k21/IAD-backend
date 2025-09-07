package main

import (
	"iad-backend/internal/api"

	"log"
)

func main() {
	log.Println("Application started")
	api.StartServer()
	log.Println("Application terminated")
}
