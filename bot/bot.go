package bot

import (
	"context"
	"github.com/dafraer/sentence-gen-tg-bot/text"
	"github.com/dafraer/sentence-gen-tg-bot/tts"
	"log"
	"strings"
	"time"

	"github.com/dafraer/sentence-gen-tg-bot/db"
	"github.com/dafraer/sentence-gen-tg-bot/gemini"
	"github.com/go-telegram/bot"
	tgbotapi "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const (
	freeSentencesAmount   = 50
	premiumPrice          = 1 //Premium subscription price in Telegram Stars
	english               = "en"
	russian               = "ru"
	maxMessageLen         = 100 //bytes
	waitingForPreferences = iota
	waitingForWord
)

type Bot struct {
	b            *tgbotapi.Bot
	store        *db.Store
	geminiClient *gemini.Client
	tts          *tts.Client
	messages     *text.Messages
}

func New(token string, store *db.Store, geminiClient *gemini.Client, ttsClient *tts.Client, messages *text.Messages) (*Bot, error) {
	bot := &Bot{store: store, geminiClient: geminiClient, tts: ttsClient, messages: messages}
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
	case update.PreCheckoutQuery != nil:
		if err := b.processPreCheckoutQuery(ctx, update); err != nil {
			log.Println(err)
		}
	case update.CallbackQuery != nil:
		if err := b.processCallbackQuery(ctx, update); err != nil {
			log.Println(err)
		}
	case update.Message != nil:
		switch {
		case update.Message.SuccessfulPayment != nil:
			if err := b.processSuccessfulPayment(ctx, update); err != nil {
				log.Println(err)
				if _, err := b.b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: update.Message.Chat.ID,
					Text:   b.messages.FailedPayment[language(update.Message.From)],
				}); err != nil {
					log.Println(err)
				}
			}
		case strings.HasPrefix(update.Message.Text, "/"):
			if err := b.processCommand(ctx, update); err != nil {
				log.Println(err)
			}
		default:
			if err := b.processMessage(ctx, update); err != nil {
				log.Println(err)
			}
		}
	}
}

func (b *Bot) sendInvoice(ctx context.Context, update *models.Update, title, desc string) error {
	_, err := b.b.SendInvoice(ctx, &tgbotapi.SendInvoiceParams{
		ChatID:      update.Message.Chat.ID,
		Title:       title,
		Description: desc,
		Currency:    "XTR",
		Payload:     "premium",
		Prices: []models.LabeledPrice{
			{
				Label:  b.messages.Premium[language(update.Message.From)],
				Amount: premiumPrice,
			},
		},
	})
	return err
}
func (b *Bot) processPreCheckoutQuery(ctx context.Context, update *models.Update) error {
	_, err := b.b.AnswerPreCheckoutQuery(ctx, &bot.AnswerPreCheckoutQueryParams{
		PreCheckoutQueryID: update.PreCheckoutQuery.ID,
		OK:                 true,
		ErrorMessage:       "",
	})
	return err
}

func (b *Bot) processSuccessfulPayment(ctx context.Context, update *models.Update) error {
	user, err := b.store.GetUser(ctx, update.Message.Chat.ID)
	if err != nil {
		return err
	}
	from := max(user.PremiumUntil, time.Now().Unix())
	if err := b.store.UpdateUserPremium(ctx, update.Message.Chat.ID, time.Unix(from, 0).Add(time.Hour*24*30).Unix()); err != nil {
		return err
	}
	_, err = b.b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   b.messages.SuccessfulPayment[language(update.Message.From)],
	})
	return err
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

// languageMarkupEn returns inline keyboard markup for selecting language in english
func languageMarkupEn() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Spanish", CallbackData: "es-ES"},
			}, {
				{Text: "French", CallbackData: "fr-FR"},
			}, {
				{Text: "German", CallbackData: "de-DE"},
			}, {
				{Text: "Turkish", CallbackData: "tr-TR"},
			}, {
				{Text: "Greek", CallbackData: "el-GR"},
			}, {
				{Text: "Russian", CallbackData: "ru-RU"},
			}, {
				{Text: "Japanese", CallbackData: "ja-JP"},
			}, {
				{Text: "Korean", CallbackData: "ko-KR"},
			}, {
				{Text: "Arabic", CallbackData: "ar-XA"},
			}, {
				{Text: "Italian", CallbackData: "it-IT"},
			},
		},
	}
}

func languageMarkupRu() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Английский", CallbackData: "en-US"},
			}, {
				{Text: "Испанский", CallbackData: "es-ES"},
			}, {
				{Text: "Французский", CallbackData: "fr-FR"},
			}, {
				{Text: "Немецкий", CallbackData: "de-DE"},
			}, {
				{Text: "Турецкий", CallbackData: "tr-TR"},
			}, {
				{Text: "Греческий", CallbackData: "el-GR"},
			}, {
				{Text: "Японский", CallbackData: "ja-JP"},
			}, {
				{Text: "Корейский", CallbackData: "ko-KR"},
			}, {
				{Text: "Арабский", CallbackData: "ar-XA"},
			}, {
				{Text: "Итальянский", CallbackData: "it-IT"},
			},
		},
	}
}
