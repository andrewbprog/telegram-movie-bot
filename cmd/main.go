package main

import (
	"context"
	"log"
	"tlgbs/internal/app"
)

func main() {
	// запуск приложения
	ctx := context.Background()

	if err := app.Start(ctx); err != nil {
		log.Fatalf("start app error: %v", err)
	}

}
