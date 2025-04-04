package bot

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/dafraer/sentence-gen-tg-bot/gemini"

	tgbotapi "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// processMessage routes message to the handler functions
func (b *Bot) processMessage(ctx context.Context, update *models.Update) {
	b.logger.Infow("Message Received", "from", update.Message.From.Username, "message", update.Message.Text)

	//Check if message is of appropriate length
	if len(update.Message.Text) > maxMessageLen {
		b.processMessageTooLong(ctx, update)
		return
	}

	//Get user from the database to check if their preferences are set
	user, err := b.store.GetUser(ctx, update.Message.Chat.ID)
	if err != nil {
		b.logger.Errorw("error getting user from the database", "error", err)
		return
	}

	//If user has set their preferences, process the word
	if user.PreferencesSet {
		b.processWord(ctx, update)
		return
	}

	//If user hasn't chosen their preferences yet prompt them to do that
	b.processPreferencesNotSet(ctx, update)
}

// processWord generates two sentences and audio for the provided word
func (b *Bot) processWord(ctx context.Context, update *models.Update) {
	//Check if message text is empty
	if update.Message.Text == "" {
		return
	}

	//Get user from the database
	user, err := b.store.GetUser(ctx, update.Message.Chat.ID)
	if err != nil {
		b.logger.Errorw("error getting user from the db", "error", err)
		return
	}

	//Add 50 free sentences if user has not used the bot today
	if user.FreeSentences <= 0 && time.Unix(user.LastUsed, 0).Before(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)) {
		user.FreeSentences = 50
	}

	//Check if user can generate sentences
	if !premium(user) && user.FreeSentences <= 0 {
		//If they can't send them message notifying them that free sentence limit has been reached
		if _, err := b.b.SendMessage(ctx, &tgbotapi.SendMessageParams{
			ChatID:      update.Message.Chat.ID,
			Text:        b.messages.LimitReached[language(update.Message.From)],
			ReplyMarkup: &models.InlineKeyboardMarkup{InlineKeyboard: [][]models.InlineKeyboardButton{{{Text: b.messages.PremiumTitle[language(update.Message.From)], CallbackData: premiumCallback}}}}}); err != nil {
			b.logger.Errorw("error sending message", "error", err)
		}
		return
	}

	//Request sentences from gemini
	res, err := b.geminiClient.Request(ctx, gemini.FormatRequestString(user.Level, user.SentenceLanguage, update.Message.Text, update.Message.From.LanguageCode), geminiProModel)
	if err != nil {
		b.logger.Errorw("error getting response from gemini", "error", err)
		return
	}
	b.logger.Debugw("Response from gemini:", "response", res)

	//Parse gemini response into 2 sentences
	sentence1, sentence2, err := parseSentences(res)
	if err != nil {
		if _, err := b.b.SendMessage(ctx, &tgbotapi.SendMessageParams{ChatID: update.Message.Chat.ID, Text: b.messages.BadRequest[language(update.Message.From)]}); err != nil {
			b.logger.Errorw("error sending message", "error", err)
		}
		return
	}

	//Generate mp3 audio
	audio, err := b.tts.Generate(ctx, sentence1, user.SentenceLanguage)
	if err != nil {
		b.logger.Errorw("error generating audio", "error", err)
		return
	}

	//Send sentences
	if _, err := b.b.SendMessage(ctx, &tgbotapi.SendMessageParams{ChatID: update.Message.Chat.ID, Text: fmt.Sprintf(b.messages.ResponseMsg[language(update.Message.From)], sentence1, sentence2), ParseMode: models.ParseModeMarkdown}); err != nil {
		b.logger.Errorw("error sending message", "error", err)
		return
	}

	//Send audio
	params := &tgbotapi.SendDocumentParams{
		ChatID:   update.Message.Chat.ID,
		Document: &models.InputFileUpload{Filename: update.Message.Text + ".mp3", Data: bytes.NewReader(audio)},
	}
	if _, err := b.b.SendDocument(ctx, params); err != nil {
		b.logger.Errorw("error sending document", "error", err)
		return
	}

	//Update user data
	user.LastUsed = time.Now().Unix()
	if !premium(user) {
		user.FreeSentences--
	}
	if err := b.store.UpdateUser(ctx, user); err != nil {
		b.logger.Errorw("error updating user", "error", err)
	}
}

// processMessageTooLong notifies user that their message is too long
func (b *Bot) processMessageTooLong(ctx context.Context, update *models.Update) {
	if _, err := b.b.SendMessage(ctx, &tgbotapi.SendMessageParams{ChatID: update.Message.Chat.ID, Text: b.messages.TooLong[language(update.Message.From)]}); err != nil {
		b.logger.Errorw("error sending message", "error", err)
	}
}

// processPreferencesNotSet notifies user that their preferences are not set
func (b *Bot) processPreferencesNotSet(ctx context.Context, update *models.Update) {
	if _, err := b.b.SendMessage(ctx, &tgbotapi.SendMessageParams{ChatID: update.Message.Chat.ID, Text: b.messages.PreferencesNotSet[language(update.Message.From)]}); err != nil {
		b.logger.Errorw("error sending message", "error", err)
	}
}
