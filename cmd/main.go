package main

import (
	"context"
	"log"
	"telegram-movie-bot/internal/app"
)

func main() {
	// запуск приложения
	ctx := context.Background()

	if err := app.Start(ctx); err != nil {
		log.Fatalf("start app error: %v", err)
	}

}
