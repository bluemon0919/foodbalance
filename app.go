package main

import (
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
)

func main() {
	projectID := os.Getenv("PROJECT_ID")
	regist := NewRegistration(projectID, "RegistrationData")

	botClient, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}
	h := NewHandler(botClient, regist)

	// Setup HTTP Server for receiving requests from Linebot
	http.HandleFunc("/callback", h.WebhookHandler)
	http.HandleFunc("/input", h.InputformHandler)
	http.HandleFunc("/post", h.PostHandler)

	log.Println("http://localhost:8080 で起動中...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
