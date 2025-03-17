package bot

import (
	"context"
	"strconv"

	"github.com/dafraer/sentence-gen-tg-bot/db"
	tgbotapi "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (b *Bot) processMessage(ctx context.Context, update *models.Update) error {
	user, err := b.store.GetUser(ctx, strconv.Itoa(int(update.Message.Chat.ID)))
	if err != nil {
		return err
	}
	switch user.State {
	case waitingForLanguage:
		b.processLanguage(ctx, update, user)
	case waitingForWord:
		b.processWord(ctx, update)
	default:
		b.processUnknownCommand(ctx, update)
	}
	return nil
}

func (b *Bot) processLanguage(ctx context.Context, update *models.Update, user *db.User) error {
	user.SentenceLanguage = update.Message.Text
	user.State = waitingForWord
	if err := b.store.SaveUser(ctx, user); err != nil {
		return err
	}
	b.b.SendMessage(ctx, &tgbotapi.SendMessageParams{ChatID: update.Message.Chat.ID, ReplyMarkup: levelsMarkup()})
	return nil
}

func (b *Bot) processWord(ctx context.Context, update *models.Update) error {
	return nil
}
