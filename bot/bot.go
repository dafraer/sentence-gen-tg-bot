package bot

import (
	"github.com/dafraer/sentence-gen-tg-bot/db"
	"github.com/dafraer/sentence-gen-tg-bot/gemini"
	"github.com/go-telegram/bot"
)

type Bot struct {
	b            *bot.Bot
	store        *db.Store
	geminiClient *gemini.Client
}

func New(token string, store *db.Store, geminiClient *gemini.Client) (*Bot, error) {
	b, err := bot.New(token)
	if err != nil {
		return nil, err
	}
	return &Bot{b, store, geminiClient}, nil
}
