package main

import (
	"context"
	"fin-tg/internal/app"
	"log"
	"log/slog"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func main() {
	ctx := context.Background()

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	b, err := bot.New(os.Getenv("TG_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	sheetSvc, err := sheets.NewService(ctx,
		option.WithCredentialsFile(os.Getenv("GOOGLE_CREDENTIALS_JSON")))
	if err != nil {
		log.Fatal(err)
	}

	newApp := app.NewApp(b, sheetSvc)

	slog.Info("starting telegram bot")
	newApp.Run(ctx)
}
