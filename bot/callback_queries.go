package bot

import (
	"context"
	tgbotapi "github.com/go-telegram/bot"
	"strconv"

	"github.com/go-telegram/bot/models"
)

func (b *Bot) processCallbackQuery(ctx context.Context, update *models.Update) error {
	//Check that data is correct
	if len(update.CallbackQuery.Data) == 2 {
		//Update user level
		if err := b.store.UpdateUserLevel(ctx, strconv.Itoa(int(update.CallbackQuery.From.ID)), update.CallbackQuery.Data); err != nil {
			return err
		}
		//Update user state
		if err := b.store.UpdateUserState(ctx, strconv.Itoa(int(update.CallbackQuery.From.ID)), waitingForWord); err != nil {
			return err
		}
		//Delete message with inline keyboard
		if _, err := b.b.DeleteMessage(ctx, &tgbotapi.DeleteMessageParams{ChatID: update.CallbackQuery.From.ID, MessageID: update.CallbackQuery.Message.Message.ID}); err != nil {
			return err
		}
		if language(&update.CallbackQuery.From) == russian {
			if _, err := b.b.SendMessage(ctx, &tgbotapi.SendMessageParams{ChatID: update.CallbackQuery.From.ID, Text: preferencesSetRu}); err != nil {
				return err
			}
		} else {
			if _, err := b.b.SendMessage(ctx, &tgbotapi.SendMessageParams{ChatID: update.CallbackQuery.From.ID, Text: preferencesSetEn}); err != nil {
				return err
			}
		}
	}
	return nil
}
