package bot

import (
	"context"

	"github.com/dafraer/sentence-gen-tg-bot/db"
	tgbotapi "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (b *Bot) processCommand(ctx context.Context, update *models.Update) error {
	switch update.Message.Text {
	case "/start":
		if err := b.processStartCommand(ctx, update); err != nil {
			return err
		}
	case "/help":
		if err := b.processHelpCommand(ctx, update); err != nil {
			return err
		}
	case "/premium":
		if err := b.processPremiumCommand(ctx, update); err != nil {
			return err
		}
	case "/preferences":
		if err := b.processPreferencesCommand(ctx, update); err != nil {
			return err
		}
	default:
		if err := b.processUnknownCommand(ctx, update); err != nil {
			return err
		}
	}
	return nil
}

func (b *Bot) processStartCommand(ctx context.Context, update *models.Update) error {
	b.store.SaveUser(ctx, &db.User{ChatId: update.Message.Chat.ID, UserName: update.Message.From.Username})
	if language(update.Message.From) == russian {
		if _, err := b.b.SendMessage(ctx, &tgbotapi.SendMessageParams{ChatID: update.Message.Chat.ID, Text: startMsgRu}); err != nil {
			return err
		}
	} else {
		if _, err := b.b.SendMessage(ctx, &tgbotapi.SendMessageParams{ChatID: update.Message.Chat.ID, Text: startMsgEn}); err != nil {
			return err
		}
	}
	return nil
}

func (b *Bot) processPreferencesCommand(ctx context.Context, update *models.Update) error {
	if language(update.Message.From) == russian {
		if _, err := b.b.SendMessage(ctx, &tgbotapi.SendMessageParams{Text: chooseLangRu}); err != nil {
			return err
		}
	} else {
		if _, err := b.b.SendMessage(ctx, &tgbotapi.SendMessageParams{Text: chooseLangEn}); err != nil {
			return err
		}
	}
	return nil
}

func (b *Bot) processHelpCommand(ctx context.Context, update *models.Update) error {
	if language(update.Message.From) == russian {
		if _, err := b.b.SendMessage(ctx, &tgbotapi.SendMessageParams{ChatID: update.Message.Chat.ID, Text: helpMsgRu}); err != nil {
			return err
		}
	} else {
		if _, err := b.b.SendMessage(ctx, &tgbotapi.SendMessageParams{ChatID: update.Message.Chat.ID, Text: helpMsgEn}); err != nil {
			return err
		}
	}
	return nil
}

func (b *Bot) processPremiumCommand(ctx context.Context, update *models.Update) error {
	return nil
}

func (b *Bot) processUnknownCommand(ctx context.Context, update *models.Update) error {
	if language(update.Message.From) == russian {
		if _, err := b.b.SendMessage(ctx, &tgbotapi.SendMessageParams{ChatID: update.Message.Chat.ID, Text: unknownCommandRu}); err != nil {
			return err
		}
	} else {
		if _, err := b.b.SendMessage(ctx, &tgbotapi.SendMessageParams{ChatID: update.Message.Chat.ID, Text: unknownCommandEn}); err != nil {
			return err
		}
	}
	return nil
}
