package bot

import (
	"context"
	"errors"
	"fmt"
	"github.com/dafraer/sentence-gen-tg-bot/text"
	"github.com/dafraer/sentence-gen-tg-bot/tts"
	"go.uber.org/zap"
	"strings"
	"time"

	"github.com/dafraer/sentence-gen-tg-bot/db"
	"github.com/dafraer/sentence-gen-tg-bot/gemini"
	"github.com/go-telegram/bot"
	tgbotapi "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const (
	premiumCallback     = "premium"
	geminiFlash         = "gemini-1.5-flash"
	geminiPro           = "gemini-2.0-pro-exp-02-05"
	freeSentencesAmount = 50
	premiumPrice        = 1 //Premium subscription price in Telegram Stars
	english             = "en"
	russian             = "ru"
	maxMessageLen       = 100 //bytes
)

type Bot struct {
	b            *tgbotapi.Bot
	store        *db.Store
	geminiClient *gemini.Client
	tts          *tts.Client
	messages     *text.Messages
	logger       *zap.SugaredLogger
}

// New creates a new bot
func New(token string, store *db.Store, geminiClient *gemini.Client, ttsClient *tts.Client, messages *text.Messages, logger *zap.SugaredLogger) (*Bot, error) {
	bot := &Bot{store: store, geminiClient: geminiClient, tts: ttsClient, messages: messages, logger: logger}
	b, err := tgbotapi.New(token, tgbotapi.WithDefaultHandler(bot.defaultHandler))
	if err != nil {
		return nil, err
	}
	bot.b = b
	return bot, nil
}

func (b *Bot) Run(ctx context.Context) {
	b.b.Start(ctx)
}

// defaultHandler routes request to the bot
func (b *Bot) defaultHandler(ctx context.Context, _ *bot.Bot, update *models.Update) {
	//Check if the update is a preCheckoutQuery, callbackQuery or message
	switch {
	case update.PreCheckoutQuery != nil:
		if err := b.processPreCheckoutQuery(ctx, update); err != nil {
			b.logger.Errorw("error processing pre checkout query", "error", err)
		}
	case update.CallbackQuery != nil:
		if err := b.processCallbackQuery(ctx, update); err != nil {
			b.logger.Errorw("error processing callback query", "error", err)
		}
	case update.Message != nil:
		//Check if the message is successful payment, command or just a message
		switch {
		case update.Message.SuccessfulPayment != nil:
			if err := b.processSuccessfulPayment(ctx, update); err != nil {
				b.logger.Errorw("error processing successful payment", "error", err)
				if _, err := b.b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: update.Message.Chat.ID,
					Text:   b.messages.FailedPayment[language(update.Message.From)],
				}); err != nil {
					b.logger.Errorw("error sending message", "error", err)
				}
			}
		case strings.HasPrefix(update.Message.Text, "/"):
			if err := b.processCommand(ctx, update); err != nil {
				b.logger.Errorw("error processing command", "error", err)
			}
		default:
			if err := b.processMessage(ctx, update); err != nil {
				b.logger.Errorw("error processing message", "error", err)
			}
		}
	}
}

func (b *Bot) sendInvoice(ctx context.Context, user *models.User, title, desc string) error {
	_, err := b.b.SendInvoice(ctx, &tgbotapi.SendInvoiceParams{
		ChatID:      user.ID,
		Title:       title,
		Description: desc,
		Currency:    "XTR",
		Payload:     "premium",
		Prices: []models.LabeledPrice{
			{
				Label:  b.messages.Premium[language(user)],
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

// premium returns true if user is still premium
func premium(user *db.User) bool {
	return !time.Unix(user.PremiumUntil, 0).Before(time.Now())
}

// canGenerate returns true if non-premium user can generate new sentences
func canGenerate(user *db.User) bool {
	return user.FreeSentences > 0 || time.Unix(user.LastUsed, 0).Before(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local))
}

// parseSentences parses gemini response into 2 sentences, returns error if fails
func parseSentences(resp string) (string, string, error) {
	//First sentence is in target language second is in user's language
	sentences := strings.Split(resp, ";")
	//Check if the sentences were not generated
	if len(sentences) < 2 {
		return "", "", errors.New(fmt.Sprintf("error parsing gemini response into sentences. Gemini repsonse: %s", resp))
	}
	return sentences[0], sentences[1], nil
}
