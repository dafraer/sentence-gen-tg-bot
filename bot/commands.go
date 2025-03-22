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
func (b *Bot) processCommand(ctx context.Context, update *models.Update) {
	b.logger.Infow("Command Received", "from", update.Message.From.Username, "command", update.Message.Text)

	switch update.Message.Text {
	case "/start":
		b.processStartCommand(ctx, update)
	case "/help":
		b.processHelpCommand(ctx, update)
	case "/premium":
		b.processPremiumCommand(ctx, update)
	case "/preferences":
		b.processPreferencesCommand(ctx, update)
	default:
		b.processUnknownCommand(ctx, update)
	}
}

// processStartCommand creates user in the database if user does not exist and sends starting message to the user
func (b *Bot) processStartCommand(ctx context.Context, update *models.Update) {
	//Create user document if it does not exist
	if _, err := b.store.GetUser(ctx, update.Message.Chat.ID); err != nil {
		if err := b.store.CreateUser(ctx, &db.User{ChatId: update.Message.Chat.ID, UserName: update.Message.From.Username, FreeSentences: freeSentencesAmount}); err != nil {
			b.logger.Errorw("error creating user int the database", "error", err)
			return
		}
	}

	//Send starting message
	if _, err := b.b.SendMessage(ctx, &tgbotapi.SendMessageParams{ChatID: update.Message.Chat.ID, Text: b.messages.Start[language(update.Message.From)]}); err != nil {
		b.logger.Errorw("error sending message", "error", err)
		return
	}
}

// processPreferencesCommand sends settings message to the user
func (b *Bot) processPreferencesCommand(ctx context.Context, update *models.Update) {
	if _, err := b.b.SendMessage(ctx, &tgbotapi.SendMessageParams{ChatID: update.Message.Chat.ID, Text: b.messages.Lang[language(update.Message.From)], ReplyMarkup: b.messages.LanguageMarkup[language(update.Message.From)]}); err != nil {
		b.logger.Errorw("error sending message", "error", err)
	}
}

// processHelpCommand sends user the list of the available commands
func (b *Bot) processHelpCommand(ctx context.Context, update *models.Update) {
	if _, err := b.b.SendMessage(ctx, &tgbotapi.SendMessageParams{ChatID: update.Message.Chat.ID, Text: b.messages.Help[language(update.Message.From)]}); err != nil {
		b.logger.Errorw("error sending message", "error", err)
	}
}

// processPremiumCommand sends user a message with an inline keyboard prompting them to buy premium.
// If the user is already premium it notifies user about it
func (b *Bot) processPremiumCommand(ctx context.Context, update *models.Update) {
	//Get user from the db to check if they already have premium
	user, err := b.store.GetUser(ctx, update.Message.Chat.ID)
	if err != nil {
		b.logger.Errorw("error getting user from the database", "error", err)
	}

	//Check if the user is already premium
	if premium(user) {
		//daysLeft stores amount of days of premium left rounded upwards.
		daysLeft := (int(time.Unix(user.PremiumUntil, 0).Sub(time.Now()).Hours()) + 23) / 24

		//Tell user that they  already have premium
		_, err := b.b.SendMessage(ctx, &tgbotapi.SendMessageParams{ChatID: update.Message.Chat.ID, Text: fmt.Sprintf(b.messages.AlreadyPremium[language(update.Message.From)], daysLeft)})
		if err != nil {
			b.logger.Errorw("error sending message", "error", err)
		}
	}

	//Send message with an inline keyboard prompting user to buy premium
	_, err = b.b.SendMessage(ctx, &tgbotapi.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        b.messages.Premium[language(update.Message.From)],
		ReplyMarkup: &models.InlineKeyboardMarkup{InlineKeyboard: [][]models.InlineKeyboardButton{{{Text: b.messages.PremiumTitle[language(update.Message.From)], CallbackData: premiumCallback}}}},
	})
	if err != nil {
		b.logger.Errorw("error sending message", "error", err)
	}
}

// processUnknownCommand sends user the message stating that the bot does not know this command
func (b *Bot) processUnknownCommand(ctx context.Context, update *models.Update) {
	if _, err := b.b.SendMessage(ctx, &tgbotapi.SendMessageParams{ChatID: update.Message.Chat.ID, Text: b.messages.UnknownCommand[language(update.Message.From)]}); err != nil {
		b.logger.Errorw("error sending message", "error", err)
	}
}
