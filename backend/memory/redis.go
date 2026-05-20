package memory

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

func Connect(host, port, password string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       0,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}

	log.Println("redis connected")
	return client
}
