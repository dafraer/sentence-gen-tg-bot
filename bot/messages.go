package bot

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/dafraer/sentence-gen-tg-bot/gemini"

	"github.com/dafraer/sentence-gen-tg-bot/db"
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
	case waitingForLanguage:
		if err := b.processLanguage(ctx, update, user); err != nil {
			return err
		}
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

func (b *Bot) processLanguage(ctx context.Context, update *models.Update, user *db.User) error {
	user.SentenceLanguage = update.Message.Text
	user.State = waitingForWord
	if err := b.store.SaveUser(ctx, user); err != nil {
		return err
	}
	if language(update.Message.From) == russian {
		if _, err := b.b.SendMessage(ctx, &tgbotapi.SendMessageParams{ChatID: update.Message.Chat.ID, Text: chooseLevelRu, ReplyMarkup: levelsMarkup()}); err != nil {
			return err
		}
	} else {
		if _, err := b.b.SendMessage(ctx, &tgbotapi.SendMessageParams{ChatID: update.Message.Chat.ID, Text: chooseLevelEn, ReplyMarkup: levelsMarkup()}); err != nil {
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
	if language(update.Message.From) == russian {
		if _, err := b.b.SendMessage(ctx, &tgbotapi.SendMessageParams{ChatID: update.Message.Chat.ID, Text: fmt.Sprintf(responseMsgRu, sentences[0], sentences[1]), ParseMode: models.ParseModeMarkdown}); err != nil {
			return err
		}
	} else {
		if _, err := b.b.SendMessage(ctx, &tgbotapi.SendMessageParams{ChatID: update.Message.Chat.ID, Text: fmt.Sprintf(responseMsgEn, sentences[0], sentences[1]), ParseMode: models.ParseModeMarkdown}); err != nil {
			return err
		}
	}

	return nil
}

func (b *Bot) processMessageTooLong(ctx context.Context, update *models.Update) error {
	if language(update.Message.From) == russian {
		if _, err := b.b.SendMessage(ctx, &tgbotapi.SendMessageParams{ChatID: update.Message.Chat.ID, Text: tooLongMsgRu}); err != nil {
			return err
		}
	} else {
		if _, err := b.b.SendMessage(ctx, &tgbotapi.SendMessageParams{ChatID: update.Message.Chat.ID, Text: tooLongMsgEn}); err != nil {
			return err
		}
	}
	return nil
}
