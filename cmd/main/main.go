package main

import (
	"github.com/abdullahelwalid/tradelog-go/pkg/utils"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"time"
	"github.com/abdullahelwalid/tradelog-go/pkg/routes"
)


func main() {
	server := http.Server{
		Addr:           ":8000",
		Handler:        routes.Mux(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	//init DB
	utils.InitDB()
	log.Printf("Server running on port 8000")
	log.Fatal(server.ListenAndServe())
}
