package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	gomail "gopkg.in/gomail.v2"
)

type BorrowEvent struct {
	ID         string `json:"id"`
	UserID     string `json:"user_id"`
	BookID     string `json:"book_id"`
	BorrowDate string `json:"borrow_date"`
	DueDate    string `json:"due_date"`
	Status     string `json:"status"`
}

func main() {
	_ = godotenv.Load()

	natsURL := getEnv("NATS_URL", nats.DefaultURL)

	nc, err := nats.Connect(
		natsURL,
		nats.Name("notification-service"),
		nats.MaxReconnects(10),
		nats.ReconnectWait(2*time.Second),
		nats.Timeout(5*time.Second),
	)
	if err != nil {
		log.Fatal("failed to connect to NATS:", err)
	}
	defer nc.Close()

	log.Println("notification-service connected to NATS")

	subscribe(nc, "borrow.created", "Borrow Created")
	subscribe(nc, "borrow.returned", "Borrow Returned")

	select {}
}

func subscribe(nc *nats.Conn, subject string, title string) {
	_, err := nc.Subscribe(subject, func(msg *nats.Msg) {
		var event BorrowEvent

		if err := json.Unmarshal(msg.Data, &event); err != nil {
			log.Println("failed to parse event:", err)
			return
		}

		body := buildEmailBody(title, event)

		if err := sendEmail(title, body); err != nil {
			log.Println("email send failed:", err)
			return
		}

		log.Println("notification handled:", subject, "borrow_id:", event.ID)
	})

	if err != nil {
		log.Fatal("failed to subscribe:", subject, err)
	}

	log.Println("subscribed to:", subject)
}

func buildEmailBody(title string, event BorrowEvent) string {
	return "E-Library Notification\n\n" +
		"Event: " + title + "\n" +
		"Borrow ID: " + event.ID + "\n" +
		"User ID: " + event.UserID + "\n" +
		"Book ID: " + event.BookID + "\n" +
		"Status: " + event.Status + "\n"
}

func sendEmail(subject string, body string) error {
	if getEnv("SMTP_ENABLED", "false") != "true" {
		log.Println("[SMTP disabled] Email subject:", subject)
		log.Println("[SMTP disabled] Email body:", body)
		return nil
	}

	port, err := strconv.Atoi(getEnv("SMTP_PORT", "587"))
	if err != nil {
		return err
	}

	message := gomail.NewMessage()
	message.SetHeader("From", os.Getenv("SMTP_FROM"))
	message.SetHeader("To", os.Getenv("SMTP_TO"))
	message.SetHeader("Subject", subject)
	message.SetBody("text/plain", body)

	dialer := gomail.NewDialer(
		os.Getenv("SMTP_HOST"),
		port,
		os.Getenv("SMTP_USER"),
		os.Getenv("SMTP_PASSWORD"),
	)

	return dialer.DialAndSend(message)
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
