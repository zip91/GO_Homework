package main

import (
	"log"

	"github.com/nats-io/nats.go"
)

func main() {

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Ошибка подключения к NATS: %v", err)
	}

	defer nc.Close()

	subject := "greetings"

	message := "Hello, NATS!"

	err = nc.Publish(subject, []byte(message))
	if err != nil {
		log.Fatalf("Ошибка отправки сообщения: %v", err)
	}

	err = nc.Flush()
	if err != nil {
		log.Fatalf("Ошибка при flush: %v", err)
	}

	log.Printf("Сообщение отправлено в канал [%s]: %s\n", subject, message)
}
