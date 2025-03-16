package bot

import (
	"context"

	"github.com/dafraer/sentence-gen-tg-bot/db"
	"github.com/dafraer/sentence-gen-tg-bot/gemini"
	"github.com/go-telegram/bot"
	tgbotapi "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Bot struct {
	b            *tgbotapi.Bot
	store        *db.Store
	geminiClient *gemini.Client
}

func New(token string, store *db.Store, geminiClient *gemini.Client) (*Bot, error) {
	b, err := tgbotapi.New(token)
	if err != nil {
		return nil, err
	}
	return &Bot{b, store, geminiClient}, nil
}

func (b *Bot) Run(ctx context.Context) error {
	b.b.RegisterHandler(bot.HandlerTypeMessageText, "/start", tgbotapi.MatchTypeExact, b.startHandler)
	b.b.Start(ctx)
	return nil
}

func (b *Bot) startHandler(ctx context.Context, _ *bot.Bot, update *models.Update) {
	b.store.SaveUser(ctx, &db.User{ChatId: update.Message.Chat.ID, UserName: update.Message.From.Username, Lang: update.Message.From.LanguageCode})
	b.b.SendMessage(ctx, &tgbotapi.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "Hello Friend"})
}
