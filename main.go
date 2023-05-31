package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"uptime-monitor/src/models"
	"uptime-monitor/src/records"
	"uptime-monitor/src/utils"

	"github.com/joho/godotenv"
)

func main() {
	loadEnvVariables()
	models.Setup()
	go records.Monitor()
	setupRoutes()
}

func loadEnvVariables() {
	if os.Getenv("GO_ENV") != "production" {
		godotenv.Load()
	}
}

func setupRoutes() {
	http.HandleFunc("/", home)
	http.HandleFunc("/records", records.Create)

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"success": true,
		"message": "Hello from Uptime Monitor!",
	}

	utils.ReturnJsonResponse(w, http.StatusOK, response)
}
