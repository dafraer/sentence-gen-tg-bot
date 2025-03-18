package bot

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/dafraer/sentence-gen-tg-bot/gemini"

	tgbotapi "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (b *Bot) processMessage(ctx context.Context, update *models.Update) error {
	//Check if message is of appropriate size
	if len(update.Message.Text) > maxMessageLen {
		if err := b.processMessageTooLong(ctx, update); err != nil {
			return err
		}
	}
	user, err := b.store.GetUser(ctx, strconv.Itoa(int(update.Message.Chat.ID)))
	if err != nil {
		return err
	}
	switch user.State {
	case waitingForWord:
		if err := b.processWord(ctx, update); err != nil {
			return err
		}
	default:
		if err := b.processUnknownCommand(ctx, update); err != nil {
			return err
		}
	}
	return nil
}

func (b *Bot) processWord(ctx context.Context, update *models.Update) error {
	user, err := b.store.GetUser(ctx, strconv.Itoa(int(update.Message.Chat.ID)))
	if err != nil {
		return err
	}
	res, err := b.geminiClient.Request(ctx, gemini.FormatRequestString(user.Level, user.SentenceLanguage, update.Message.Text, update.Message.From.LanguageCode))
	if err != nil {
		return err
	}
	//First sentence is in target language second is in user's language
	sentences := strings.Split(res, ";")
	//Check if the sentences were not generated
	//TODO: fix this shit
	if len(sentences) < 2 {
		if _, err := b.b.SendMessage(ctx, &tgbotapi.SendMessageParams{ChatID: update.Message.Chat.ID, Text: b.messages.BadRequest[language(update.Message.From)]}); err != nil {
			return err
		}
	}
	audio, err := b.tts.Generate(ctx, sentences[0], user.SentenceLanguage)
	if err != nil {
		return err
	}
	if _, err := b.b.SendMessage(ctx, &tgbotapi.SendMessageParams{ChatID: update.Message.Chat.ID, Text: fmt.Sprintf(b.messages.ResponseMsg[language(update.Message.From)], sentences[0], sentences[1]), ParseMode: models.ParseModeMarkdown}); err != nil {
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

	return nil
}

func (b *Bot) processMessageTooLong(ctx context.Context, update *models.Update) error {
	if _, err := b.b.SendMessage(ctx, &tgbotapi.SendMessageParams{ChatID: update.Message.Chat.ID, Text: b.messages.TooLong[language(update.Message.From)]}); err != nil {
		return err
	}
	return nil
}
