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

func (b *Bot) processMessage(ctx context.Context, update *models.Update) error {
	//Check if message is of appropriate length
	if len(update.Message.Text) > maxMessageLen {
		if err := b.processMessageTooLong(ctx, update); err != nil {
			return err
		}
	}

	//Get user from the database to check if their preferences are set
	user, err := b.store.GetUser(ctx, update.Message.Chat.ID)
	if err != nil {
		return err
	}

	if user.PreferencesSet {
		if err := b.processWord(ctx, update); err != nil {
			return err
		}
		return nil
	}

	//If user hasn't chosen their preferences yet prompt them to do that
	if err := b.processPreferencesNotSet(ctx, update); err != nil {
		return err
	}
	return nil
}

func (b *Bot) processWord(ctx context.Context, update *models.Update) error {
	//Check if message text is empty
	if update.Message.Text == "" {
		return nil
	}

	//Get user from the database
	user, err := b.store.GetUser(ctx, update.Message.Chat.ID)
	if err != nil {
		return err
	}

	//Add 50 free sentences if user has not used the bot today
	if user.FreeSentences <= 0 && time.Unix(user.LastUsed, 0).Before(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)) {
		user.FreeSentences = 50
	}

	//Check if user can generate sentences
	if !premium(user) && user.FreeSentences <= 0 {
		if _, err := b.b.SendMessage(ctx, &tgbotapi.SendMessageParams{
			ChatID:      update.Message.Chat.ID,
			Text:        b.messages.LimitReached[language(update.Message.From)],
			ReplyMarkup: &models.InlineKeyboardMarkup{InlineKeyboard: [][]models.InlineKeyboardButton{{{Text: b.messages.PremiumTitle[language(update.Message.From)], CallbackData: premiumCallback}}}}}); err != nil {
			return err
		}
		return nil
	}

	//Use better model if the user has premium
	geminiVersion := geminiFlash
	if premium(user) {
		geminiVersion = geminiPro
	}

	//Request sentences from gemini
	res, err := b.geminiClient.Request(ctx, gemini.FormatRequestString(user.Level, user.SentenceLanguage, update.Message.Text, update.Message.From.LanguageCode), geminiVersion)
	if err != nil {
		return err
	}
	b.logger.Debugw("Response from gemini:", "response", res)
	//Parse gemini response into 2 sentences
	sentence1, sentence2, err := parseSentences(res)
	if err != nil {
		if _, err := b.b.SendMessage(ctx, &tgbotapi.SendMessageParams{ChatID: update.Message.Chat.ID, Text: b.messages.BadRequest[language(update.Message.From)]}); err != nil {
			b.logger.Errorw("error sending bad request", "error", err)
			return err
		}
		return nil
	}

	//Generate mp3 audio
	audio, err := b.tts.Generate(ctx, sentence1, user.SentenceLanguage)
	if err != nil {
		b.logger.Errorw("error generating audio", "error", err)
		return err
	}
	if _, err := b.b.SendMessage(ctx, &tgbotapi.SendMessageParams{ChatID: update.Message.Chat.ID, Text: fmt.Sprintf(b.messages.ResponseMsg[language(update.Message.From)], sentence1, sentence2), ParseMode: models.ParseModeMarkdown}); err != nil {
		b.logger.Errorw("error sending message", "error", err)
		return err
	}

	//Send audio
	params := &tgbotapi.SendDocumentParams{
		ChatID:   update.Message.Chat.ID,
		Document: &models.InputFileUpload{Filename: "audio.mp3", Data: bytes.NewReader(audio.AudioContent)},
	}
	if _, err := b.b.SendDocument(ctx, params); err != nil {
		return err
	}

	//Update user
	user.LastUsed = time.Now().Unix()
	if time.Unix(user.PremiumUntil, 0).Before(time.Now()) {
		user.FreeSentences--
	}
	if err := b.store.UpdateUser(ctx, user); err != nil {
		return err
	}
	return nil
}

func (b *Bot) processMessageTooLong(ctx context.Context, update *models.Update) error {
	if _, err := b.b.SendMessage(ctx, &tgbotapi.SendMessageParams{ChatID: update.Message.Chat.ID, Text: b.messages.TooLong[language(update.Message.From)]}); err != nil {
		return err
	}
	return nil
}

func (b *Bot) processPreferencesNotSet(ctx context.Context, update *models.Update) error {
	if _, err := b.b.SendMessage(ctx, &tgbotapi.SendMessageParams{ChatID: update.Message.Chat.ID, Text: b.messages.PreferencesNotSet[language(update.Message.From)]}); err != nil {
		return err
	}
	return nil
}
