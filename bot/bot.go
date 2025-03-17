package bot

import (
	"context"
	"log"
	"strings"

	"github.com/dafraer/sentence-gen-tg-bot/db"
	"github.com/dafraer/sentence-gen-tg-bot/gemini"
	"github.com/go-telegram/bot"
	tgbotapi "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const (
	english            = "en"
	russian            = "ru"
	maxMessageLen      = 100 //bytes
	waitingForLanguage = iota
	waitingForWord
)

type Bot struct {
	b            *tgbotapi.Bot
	store        *db.Store
	geminiClient *gemini.Client
}

func New(token string, store *db.Store, geminiClient *gemini.Client) (*Bot, error) {
	bot := &Bot{store: store, geminiClient: geminiClient}
	b, err := tgbotapi.New(token, tgbotapi.WithDefaultHandler(bot.defaultHandler))
	if err != nil {
		return nil, err
	}
	bot.b = b
	return bot, nil
}

func (b *Bot) Run(ctx context.Context) error {
	b.b.Start(ctx)
	return nil
}

func (b *Bot) defaultHandler(ctx context.Context, _ *bot.Bot, update *models.Update) {
	switch {
	case update.CallbackQuery != nil:
		if err := b.processCallbackQuery(ctx, update); err != nil {
			log.Println(err)
		}
	case update.Message != nil:
		if strings.HasPrefix(update.Message.Text, "/") {
			if err := b.processCommand(ctx, update); err != nil {
				log.Println(err)
			}
		} else {
			if err := b.processMessage(ctx, update); err != nil {
				log.Println(err)
			}
		}
	}
}

// Returns user's language code if its russian or english, else returns english
func language(user *models.User) string {
	switch user.LanguageCode {
	case "ru":
		return russian
	default:
		return english
	}
}

// levelsMarkup returns inline keyboard markup for selecting language level
func levelsMarkup() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "A1", CallbackData: "A1"},
			}, {
				{Text: "A2", CallbackData: "A2"},
			}, {
				{Text: "B1", CallbackData: "B1"},
			}, {
				{Text: "B2", CallbackData: "B2"},
			}, {
				{Text: "C1", CallbackData: "C1"},
			}, {
				{Text: "C2", CallbackData: "C2"},
			},
		},
	}
}
