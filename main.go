package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"gopkg.in/telebot.v3"
)

// InfluxDBNotification represents the JSON payload from InfluxDB
type InfluxDBNotification struct {
	CheckID      string `json:"_check_id"`
	CheckName    string `json:"_check_name"`
	Level        string `json:"_level"`
	Message      string `json:"_message"`
	Measurement  string `json:"_measurement"`
	Time         string `json:"_time"`
	Type         string `json:"_type"`
	Status       string `json:"_status"`
	Notification struct {
		EndpointID   string `json:"endpointID"`
		EndpointName string `json:"endpointName"`
	} `json:"_notification"`
}

func main() {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	channelID := os.Getenv("TELEGRAM_CHANNEL_ID")

	if botToken == "" || channelID == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN and TELEGRAM_CHANNEL_ID must be set")
	}

	http.HandleFunc("/notify", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var notification InfluxDBNotification
		if err := json.NewDecoder(r.Body).Decode(&notification); err != nil {
			http.Error(w, "Error decoding request body", http.StatusBadRequest)
			return
		}

		message := fmt.Sprintf(
			"InfluxDB Alert:\n"+
				"Check: %s (%s)\n"+
				"Level: %s\n"+
				"Message: %s\n"+
				"Measurement: %s\n"+
				"Time: %s",
			notification.CheckName, notification.CheckID,
			notification.Level,
			notification.Message,
			notification.Measurement,
			notification.Time,
		)

		log.Println("Received notification:", message)

		// Send message to Telegram
		pref := telebot.Settings{
			Token:  botToken,
			Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
		}

		b, err := telebot.NewBot(pref)
		if err != nil {
			log.Println("Error creating bot:", err)
			return
		}

		var recipient telebot.Recipient
		if chatID, err := strconv.ParseInt(channelID, 10, 64); err == nil {
			recipient = telebot.ChatID(chatID)
		} else {
			recipient = &telebot.Chat{Username: channelID}
		}

		_, err = b.Send(recipient, message)
		if err != nil {
			log.Println("Error sending message:", err)
			return
		}

		fmt.Fprintln(w, "Notification received and sent to Telegram")
	})

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
