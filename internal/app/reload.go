package app

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (a App) reloadHandler(ctx context.Context, _ *bot.Bot, update *models.Update) {
	a.chats.mu.Lock()
	defer a.chats.mu.Unlock()

	cfg, err := a.loadConfig(ctx, update.Message.Chat.ID)
	if err != nil {
		return
	}

	a.chats.chats[update.Message.Chat.ID] = cfg

	return
}
