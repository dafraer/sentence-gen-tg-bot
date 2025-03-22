package bot

import (
	"context"
	"fmt"
	"github.com/dafraer/sentence-gen-tg-bot/db"
	tgbotapi "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"time"
)

// processCommand routes command to the method that handles it
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

// processStartCommand creates user in the database if user does not exist and sends starting message to the user
func (b *Bot) processStartCommand(ctx context.Context, update *models.Update) error {
	//Create user document if it does not exist
	if _, err := b.store.GetUser(ctx, update.Message.Chat.ID); err != nil {
		if err := b.store.CreateUser(ctx, &db.User{ChatId: update.Message.Chat.ID, UserName: update.Message.From.Username, FreeSentences: freeSentencesAmount}); err != nil {
			return err
		}
	}

	//Send starting message
	if _, err := b.b.SendMessage(ctx, &tgbotapi.SendMessageParams{ChatID: update.Message.Chat.ID, Text: b.messages.Start[language(update.Message.From)]}); err != nil {
		return err
	}
	return nil
}

// processPreferencesCommand sends settings message to the user
func (b *Bot) processPreferencesCommand(ctx context.Context, update *models.Update) error {
	if _, err := b.b.SendMessage(ctx, &tgbotapi.SendMessageParams{ChatID: update.Message.Chat.ID, Text: b.messages.Lang[language(update.Message.From)], ReplyMarkup: b.messages.LanguageMarkup[language(update.Message.From)]}); err != nil {
		return err
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
	//Get user to check if they already have premium
	user, err := b.store.GetUser(ctx, update.Message.Chat.ID)
	if err != nil {
		return err
	}
	if premium(user) {
		daysLeft := (int(time.Unix(user.PremiumUntil, 0).Sub(time.Now()).Hours()) + 23) / 24
		_, err := b.b.SendMessage(ctx, &tgbotapi.SendMessageParams{ChatID: update.Message.Chat.ID, Text: fmt.Sprintf(b.messages.AlreadyPremium[language(update.Message.From)], daysLeft)})
		return err
	}
	_, err = b.b.SendMessage(ctx, &tgbotapi.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        b.messages.Premium[language(update.Message.From)],
		ReplyMarkup: &models.InlineKeyboardMarkup{InlineKeyboard: [][]models.InlineKeyboardButton{{{Text: b.messages.PremiumTitle[language(update.Message.From)], CallbackData: premiumCallback}}}},
	})
	return err
}

func (b *Bot) processUnknownCommand(ctx context.Context, update *models.Update) error {
	if _, err := b.b.SendMessage(ctx, &tgbotapi.SendMessageParams{ChatID: update.Message.Chat.ID, Text: b.messages.UnknownCommand[language(update.Message.From)]}); err != nil {
		return err
	}
	return nil
}
