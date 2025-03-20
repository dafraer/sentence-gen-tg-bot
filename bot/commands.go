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
	if err := b.store.SaveUser(ctx, &db.User{ChatId: update.Message.Chat.ID, UserName: update.Message.From.Username, FreeSentences: freeSentencesAmount}); err != nil {
		return err
	}
	if _, err := b.b.SendMessage(ctx, &tgbotapi.SendMessageParams{ChatID: update.Message.Chat.ID, Text: b.messages.Start[language(update.Message.From)]}); err != nil {
		return err
	}
	return nil
}

func (b *Bot) processPreferencesCommand(ctx context.Context, update *models.Update) error {
	if language(update.Message.From) == russian {
		if _, err := b.b.SendMessage(ctx, &tgbotapi.SendMessageParams{ChatID: update.Message.Chat.ID, Text: chooseLangRu, ReplyMarkup: languageMarkupRu()}); err != nil {
			return err
		}
	} else {
		if _, err := b.b.SendMessage(ctx, &tgbotapi.SendMessageParams{ChatID: update.Message.Chat.ID, Text: chooseLangEn, ReplyMarkup: languageMarkupEn()}); err != nil {
			return err
		}
	}
	return nil
}

func (b *Bot) processHelpCommand(ctx context.Context, update *models.Update) error {
	if _, err := b.b.SendMessage(ctx, &tgbotapi.SendMessageParams{ChatID: update.Message.Chat.ID, Text: b.messages.Help[language(update.Message.From)]}); err != nil {
		return err
	}
	return nil
}

func (b *Bot) processPremiumCommand(ctx context.Context, update *models.Update) error {
	return b.sendInvoice(ctx, update, b.messages.PremiumTitle[language(update.Message.From)], b.messages.Premium[language(update.Message.From)])
}

func (b *Bot) processUnknownCommand(ctx context.Context, update *models.Update) error {
	if _, err := b.b.SendMessage(ctx, &tgbotapi.SendMessageParams{ChatID: update.Message.Chat.ID, Text: b.messages.UnknownCommand[language(update.Message.From)]}); err != nil {
		return err
	}
	return nil
}
