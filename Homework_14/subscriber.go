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

	_, err = nc.Subscribe(subject, func(msg *nats.Msg) {
		log.Printf("Получено сообщение из [%s]: %s\n", msg.Subject, string(msg.Data))
	})
	if err != nil {
		log.Fatalf("Ошибка подписки: %v", err)
	}

	log.Printf("Подписчик слушает канал [%s]...\n", subject)

	select {}
}
