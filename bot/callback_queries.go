package bot

import (
	"context"
	"strconv"

	"github.com/go-telegram/bot/models"
)

func (b *Bot) processCallbackQuery(ctx context.Context, update *models.Update) error {
	//Check that data is correct
	if len(update.CallbackQuery.Data) == 2 {
		b.store.UpdateUserLevel(ctx, strconv.Itoa(int(update.CallbackQuery.From.ID)), update.CallbackQuery.Data)

	}
	return nil
}
