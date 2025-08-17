package app

import (
	"context"
	"errors"
	"strings"

	"github.com/go-telegram/bot"
)

type Config struct {
	Values map[string]string
}

func (a App) loadConfig(ctx context.Context, chatID int64) (*Config, error) {
	chat, err := a.bot.GetChat(ctx, &bot.GetChatParams{
		ChatID: chatID,
	})
	if err != nil {
		a.writeError(ctx, "не удалось получить чат", chatID, err)
		return nil, err
	}

	if chat.PinnedMessage == nil {
		a.writeError(ctx, "нет закрепленного сообщения", chatID, errors.New("pinned not found"))
		return nil, err
	}
	return a.parseConfig(chat.PinnedMessage.Text), nil
}

func (a App) parseConfig(text string) *Config {
	cfg := &Config{Values: make(map[string]string)}
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			// убрать комментарии
			if idx := strings.Index(val, "#"); idx != -1 {
				val = strings.TrimSpace(val[:idx])
			}
			cfg.Values[key] = val
		}
	}
	return cfg
}

func (a App) getSheetID(ctx context.Context, id int64) string {
	var err error
	a.chats.mu.RLock()
	cfg, ok := a.chats.chats[id]

	if !ok {
		a.chats.mu.RUnlock()
		a.chats.mu.Lock()
		cfg, err = a.loadConfig(ctx, id)
		if err != nil {
			return ""
		}
		a.chats.chats[id] = cfg
		a.chats.mu.Unlock()
	}

	return cfg.Values["sheet_id"]
}
