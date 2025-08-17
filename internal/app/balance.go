package app

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (a App) balanceHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	vals, err := a.sheetSvc.Spreadsheets.Values.Get(a.getSheetID(ctx, update.Message.Chat.ID), "A:E").Do()
	if err != nil {
		a.writeError(ctx, "чтение таблицы", update.Message.Chat.ID, err)
		return
	}
	total := 0
	for _, row := range vals.Values {
		if len(row) < 3 {
			continue
		}
		amountStr, ok := row[2].(string)
		if !ok {
			continue
		}
		amt, err := strconv.Atoi(amountStr)
		if err == nil {
			total += amt
		}
	}
	msg := fmt.Sprintf("Текущий баланс: %d", total)
	_, err = a.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   msg,
	})
	if err != nil {
		slog.Error("ошибка записи баланса", "err", err)
	}
}
