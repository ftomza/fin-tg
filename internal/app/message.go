package app

import (
	"context"
	"regexp"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"google.golang.org/api/sheets/v4"
)

var reMoney = regexp.MustCompile(`^([+-]\d+)(?:\s+(\S+))?(?:\s+(.*))?$`)

func (a App) messageHandler(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	if !reMoney.MatchString(update.Message.Text) {
		return
	}

	m := reMoney.FindStringSubmatch(update.Message.Text)
	amount := m[1]
	category := m[2]
	comment := m[3]

	row := []interface{}{
		time.Now().Format("2006-01-02 15:04:05"),
		update.Message.From.Username,
		amount,
		category,
		comment,
	}

	_, err := a.sheetSvc.Spreadsheets.Values.Append(a.getSheetID(ctx, update.Message.Chat.ID), "A:E",
		&sheets.ValueRange{Values: [][]interface{}{row}}).
		ValueInputOption("RAW").Do()
	if err != nil {
		a.writeError(ctx, "запись в таблицу", update.Message.Chat.ID, err)
	}
}
