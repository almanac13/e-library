package events

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

type Publisher struct {
	conn *nats.Conn
}

func NewPublisher() *Publisher {
	url := os.Getenv("NATS_URL")
	if url == "" {
		url = nats.DefaultURL
	}

	nc, err := nats.Connect(
		url,
		nats.Name("borrow-service"),
		nats.MaxReconnects(10),
		nats.ReconnectWait(2*time.Second),
		nats.Timeout(5*time.Second),
	)

	if err != nil {
		log.Println("failed to connect to NATS:", err)
		return nil
	}

	log.Println("borrow-service connected to NATS")

	return &Publisher{
		conn: nc,
	}
}

func (p *Publisher) Publish(subject string, data any) {
	if p == nil || p.conn == nil {
		return
	}

	payload, err := json.Marshal(data)
	if err != nil {
		log.Println("failed to marshal event:", err)
		return
	}

	err = p.conn.Publish(subject, payload)
	if err != nil {
		log.Println("publish error:", err)
		return
	}

	_ = p.conn.Flush()
	log.Println("published event:", subject)
}

func (p *Publisher) Close() {
	if p != nil && p.conn != nil {
		p.conn.Close()
	}
}
