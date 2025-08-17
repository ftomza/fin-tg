package app

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/go-telegram/bot"
	"google.golang.org/api/sheets/v4"
)

type App struct {
	sheetSvc *sheets.Service
	chats    *ChatConfig
	bot      *bot.Bot
}

type ChatConfig struct {
	mu    sync.RWMutex
	chats map[int64]*Config
}

func (a App) Run(ctx context.Context) {

	bot.WithDefaultHandler(a.messageHandler)(a.bot)

	a.bot.RegisterHandler(bot.HandlerTypeMessageText, "balance@fin_bal_bot", bot.MatchTypeCommand, a.balanceHandler)
	a.bot.RegisterHandler(bot.HandlerTypeMessageText, "reload@fin_bal_bot", bot.MatchTypeCommand, a.reloadHandler)
	a.bot.RegisterHandler(bot.HandlerTypeMessageText, "balance", bot.MatchTypeCommand, a.balanceHandler)
	a.bot.RegisterHandler(bot.HandlerTypeMessageText, "reload", bot.MatchTypeCommand, a.reloadHandler)

	a.bot.Start(ctx)
}

func (a App) writeError(ctx context.Context, s string, chatID int64, err error) {
	slog.Error(s, "err", err)
	_, err = a.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   fmt.Sprintf("ОШИБКА: %s", s),
	})

	if err != nil {
		slog.Error("запись ошибки в чат", "err", err)
	}
}

func NewApp(bot *bot.Bot, sheetSvc *sheets.Service) *App {
	return &App{
		bot:      bot,
		sheetSvc: sheetSvc,
		chats: &ChatConfig{
			mu:    sync.RWMutex{},
			chats: make(map[int64]*Config),
		},
	}
}
