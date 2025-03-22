package bot

import (
	"context"
	tgbotapi "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const (
	levelCallbackLength    = 2
	languageCallbackLength = 5
)

// processCallbackQuery routes callback to the handler functions
func (b *Bot) processCallbackQuery(ctx context.Context, update *models.Update) {
	b.logger.Infow("Callback Query Received", "from", update.Message.From.Username, "callback data", update.CallbackQuery.Data)
	switch {
	//callback to buy premium
	case update.CallbackQuery.Data == premiumCallback:
		b.processPremiumCallback(ctx, update)
	//If the length of the callback data is equal to the level length (e.g. A1) - process level callback
	case len(update.CallbackQuery.Data) == levelCallbackLength:
		b.processLevelCallback(ctx, update)
	//If the length of the callback data is equal to the language code length (e.g. es-ES) - process language callback
	case len(update.CallbackQuery.Data) == languageCallbackLength:
		b.processLanguageCallback(ctx, update)
	}
}

// processLanguageCallback handles callback choosing language
func (b *Bot) processLanguageCallback(ctx context.Context, update *models.Update) {
	//Update user sentence language
	if err := b.store.SetUserSentenceLanguage(ctx, update.CallbackQuery.From.ID, update.CallbackQuery.Data); err != nil {
		b.logger.Errorw("failed to set user sentence language", "err", err)
		return
	}

	//Delete message with inline keyboard
	if _, err := b.b.DeleteMessage(ctx, &tgbotapi.DeleteMessageParams{ChatID: update.CallbackQuery.From.ID, MessageID: update.CallbackQuery.Message.Message.ID}); err != nil {
		b.logger.Errorw("failed to delete message", "err", err)
		return
	}

	//Prompt user to choose language level (e.g. A1)
	if _, err := b.b.SendMessage(ctx, &tgbotapi.SendMessageParams{ChatID: update.CallbackQuery.From.ID, Text: b.messages.Level[language(&update.CallbackQuery.From)], ReplyMarkup: levelsMarkup()}); err != nil {
		b.logger.Errorw("failed to send level message", "err", err)
	}
}

// processLevelCallback handles callback choosing language level
func (b *Bot) processLevelCallback(ctx context.Context, update *models.Update) {
	//Update user language level
	if err := b.store.SetUserLevel(ctx, update.CallbackQuery.From.ID, update.CallbackQuery.Data); err != nil {
		b.logger.Errorw("failed to set user language level", "err", err)
		return
	}
	//Delete message with inline keyboard
	if _, err := b.b.DeleteMessage(ctx, &tgbotapi.DeleteMessageParams{ChatID: update.CallbackQuery.From.ID, MessageID: update.CallbackQuery.Message.Message.ID}); err != nil {
		b.logger.Errorw("failed to delete message", "err", err)
		return
	}

	//Send message telling user that everything is set
	if _, err := b.b.SendMessage(ctx, &tgbotapi.SendMessageParams{ChatID: update.CallbackQuery.From.ID, Text: b.messages.PreferencesSet[language(&update.CallbackQuery.From)]}); err != nil {
		b.logger.Errorw("failed to send message", "err", err)
	}
}

// processPremiumCallback sends an invoice to buy premium to the user
func (b *Bot) processPremiumCallback(ctx context.Context, update *models.Update) {
	lang := language(&update.CallbackQuery.From)

	//Send the invoice
	if err := b.sendInvoice(ctx, &update.CallbackQuery.From, b.messages.PremiumTitle[lang], b.messages.PremiumDescription[lang]); err != nil {
		b.logger.Errorw("failed to send invoice", "err", err)
		return
	}

	//Delete message with inline keyboard
	_, err := b.b.DeleteMessage(ctx, &tgbotapi.DeleteMessageParams{ChatID: update.CallbackQuery.From.ID, MessageID: update.CallbackQuery.Message.Message.ID})
	if err != nil {
		b.logger.Errorw("failed to delete message", "err", err)
	}
}
