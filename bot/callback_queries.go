package bot

import (
	"context"
	tgbotapi "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"strconv"
)

func (b *Bot) processCallbackQuery(ctx context.Context, update *models.Update) error {
	switch len(update.CallbackQuery.Data) {
	case 2:
		if err := b.processLevelCallback(ctx, update); err != nil {
			return err
		}
	case 5:
		if err := b.processLanguageCallback(ctx, update); err != nil {
			return err
		}
	}
	return nil
}

func (b *Bot) processLanguageCallback(ctx context.Context, update *models.Update) error {
	user, err := b.store.GetUser(ctx, strconv.Itoa(int(update.CallbackQuery.From.ID)))
	if err != nil {
		return err
	}
	user.SentenceLanguage = update.CallbackQuery.Data
	user.State = waitingForWord
	if err := b.store.SaveUser(ctx, user); err != nil {
		return err
	}
	//Delete message with inline keyboard
	if _, err := b.b.DeleteMessage(ctx, &tgbotapi.DeleteMessageParams{ChatID: update.CallbackQuery.From.ID, MessageID: update.CallbackQuery.Message.Message.ID}); err != nil {
		return err
	}
	if _, err := b.b.SendMessage(ctx, &tgbotapi.SendMessageParams{ChatID: update.CallbackQuery.From.ID, Text: b.messages.Level[language(&update.CallbackQuery.From)], ReplyMarkup: levelsMarkup()}); err != nil {
		return err
	}
	return nil
}

func (b *Bot) processLevelCallback(ctx context.Context, update *models.Update) error {
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
	if _, err := b.b.SendMessage(ctx, &tgbotapi.SendMessageParams{ChatID: update.CallbackQuery.From.ID, Text: b.messages.PreferencesSet[language(&update.CallbackQuery.From)]}); err != nil {
		return err
	}
	return nil
}
