package bot

import (
	"context"
	tgbotapi "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (b *Bot) processCallbackQuery(ctx context.Context, update *models.Update) error {
	switch {
	case update.CallbackQuery.Data == premiumCallback:
		if err := b.processPremiumCallback(ctx, update); err != nil {
			return err
		}
	case len(update.CallbackQuery.Data) == 2:
		if err := b.processLevelCallback(ctx, update); err != nil {
			return err
		}
	case len(update.CallbackQuery.Data) == 5:
		if err := b.processLanguageCallback(ctx, update); err != nil {
			return err
		}
	}
	return nil
}

func (b *Bot) processLanguageCallback(ctx context.Context, update *models.Update) error {
	//Update user sentence language
	if err := b.store.SetUserSentenceLanguage(ctx, update.CallbackQuery.From.ID, update.CallbackQuery.Data); err != nil {
		return err
	}

	//Delete message with inline keyboard
	if _, err := b.b.DeleteMessage(ctx, &tgbotapi.DeleteMessageParams{ChatID: update.CallbackQuery.From.ID, MessageID: update.CallbackQuery.Message.Message.ID}); err != nil {
		return err
	}

	//Prompt user to choose language level (e.g. A1)
	if _, err := b.b.SendMessage(ctx, &tgbotapi.SendMessageParams{ChatID: update.CallbackQuery.From.ID, Text: b.messages.Level[language(&update.CallbackQuery.From)], ReplyMarkup: levelsMarkup()}); err != nil {
		return err
	}
	return nil
}

func (b *Bot) processLevelCallback(ctx context.Context, update *models.Update) error {
	//Update user language level
	if err := b.store.SetUserLevel(ctx, update.CallbackQuery.From.ID, update.CallbackQuery.Data); err != nil {
		return err
	}
	//Delete message with inline keyboard
	if _, err := b.b.DeleteMessage(ctx, &tgbotapi.DeleteMessageParams{ChatID: update.CallbackQuery.From.ID, MessageID: update.CallbackQuery.Message.Message.ID}); err != nil {
		return err
	}

	//Send message telling user that everything is set
	if _, err := b.b.SendMessage(ctx, &tgbotapi.SendMessageParams{ChatID: update.CallbackQuery.From.ID, Text: b.messages.PreferencesSet[language(&update.CallbackQuery.From)]}); err != nil {
		return err
	}
	return nil
}

func (b *Bot) processPremiumCallback(ctx context.Context, update *models.Update) error {
	lang := language(&update.CallbackQuery.From)
	if err := b.sendInvoice(ctx, &update.CallbackQuery.From, b.messages.PremiumTitle[lang], b.messages.PremiumDescription[lang]); err != nil {
		return err
	}
	//Delete message with inline keyboard
	_, err := b.b.DeleteMessage(ctx, &tgbotapi.DeleteMessageParams{ChatID: update.CallbackQuery.From.ID, MessageID: update.CallbackQuery.Message.Message.ID})
	return err
}
