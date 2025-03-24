package bot

import (
	"context"
	"errors"
	"fmt"
	"github.com/dafraer/sentence-gen-tg-bot/text"
	"github.com/dafraer/sentence-gen-tg-bot/tts"
	"go.uber.org/zap"
	"net/http"
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
	geminiProModel      = "gemini-2.0-pro-exp-02-05"
	freeSentencesAmount = 50
	premiumPrice        = 100 //Premium subscription price in Telegram Stars
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
	//Create bot using provided dependencies
	bot := &Bot{store: store, geminiClient: geminiClient, tts: ttsClient, messages: messages, logger: logger}

	//Create telegram bot with a default handler
	b, err := tgbotapi.New(token, tgbotapi.WithDefaultHandler(bot.defaultHandler))
	if err != nil {
		return nil, err
	}
	bot.b = b
	return bot, nil
}

// Run runs the bot using long polling
func (b *Bot) Run(ctx context.Context) {
	b.b.Start(ctx)
}

// RunWebhook runs bot using webhook
func (b *Bot) RunWebhook(ctx context.Context, address string) error {
	//delete webhook before shutdown
	defer func() {
		if _, err := b.b.DeleteWebhook(context.Background(), &tgbotapi.DeleteWebhookParams{DropPendingUpdates: true}); err != nil {
			panic(err)
		}
	}()
	go b.b.StartWebhook(ctx)

	//Set tup server for the webhook
	//Create http server for the webhook
	srv := &http.Server{
		Addr:    address,
		Handler: b.b.WebhookHandler(),
	}

	//Create channel for errors
	ch := make(chan error)

	//Run server in a goroutine
	go func() {
		defer close(ch)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			ch <- err
			return
		}
		ch <- nil
	}()

	//Wait for either shutdown or an error
	select {
	case <-ctx.Done():
		if err := srv.Shutdown(context.Background()); err != nil {
			return err
		}
		err := <-ch
		if err != nil {
			return err
		}
	case err := <-ch:
		return err
	}
	return nil
}

// defaultHandler routes request to the bot
func (b *Bot) defaultHandler(ctx context.Context, _ *bot.Bot, update *models.Update) {
	//Check if the update is a preCheckoutQuery, callbackQuery or message
	switch {
	case update.PreCheckoutQuery != nil:
		b.processPreCheckoutQuery(ctx, update)
	case update.CallbackQuery != nil:
		b.processCallbackQuery(ctx, update)
	case update.Message != nil:
		//Check if the message is successful payment, command or just a message
		switch {
		case update.Message.SuccessfulPayment != nil:
			if err := b.processSuccessfulPayment(ctx, update); err != nil {
				//If we can't process successful payment send user message about it
				b.logger.Errorw("error processing successful payment", "error", err)
				if _, err := b.b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: update.Message.Chat.ID,
					Text:   b.messages.FailedPayment[language(update.Message.From)],
				}); err != nil {
					b.logger.Errorw("error sending message", "error", err)
				}
			}
		case strings.HasPrefix(update.Message.Text, "/"):
			b.processCommand(ctx, update)
		default:
			b.processMessage(ctx, update)
		}
	}
}

// sendInvoice sends invoice for premium subscription to the user
func (b *Bot) sendInvoice(ctx context.Context, user *models.User, title, desc string) error {
	_, err := b.b.SendInvoice(ctx, &tgbotapi.SendInvoiceParams{
		ChatID:      user.ID,
		Title:       title,
		Description: desc,
		Currency:    "XTR", // Telegram stars
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

// processPreCheckoutQuery process pre checkout query that is sent to us to confirm payment
func (b *Bot) processPreCheckoutQuery(ctx context.Context, update *models.Update) {
	_, err := b.b.AnswerPreCheckoutQuery(ctx, &bot.AnswerPreCheckoutQueryParams{
		PreCheckoutQueryID: update.PreCheckoutQuery.ID,
		OK:                 true,
		ErrorMessage:       "",
	})
	if err != nil {
		b.logger.Errorw("error processing pre checkout query", "error", err)
	}
}

// processSuccessfulPayment gives user premium and sends them a message saying that payment has been successful
func (b *Bot) processSuccessfulPayment(ctx context.Context, update *models.Update) error {
	//Get user from the database
	user, err := b.store.GetUser(ctx, update.Message.Chat.ID)
	if err != nil {
		b.logger.Errorw("error getting user from the database", "error", err)
		return err
	}

	//from is a time from which we start the 30 days premium period. It is set to the last day of user's premium if they still have it
	from := max(user.PremiumUntil, time.Now().Unix())

	//Add 30 days to user's premium
	if err := b.store.UpdateUserPremium(ctx, update.Message.Chat.ID, time.Unix(from, 0).Add(time.Hour*24*30).Unix()); err != nil {
		b.logger.Errorw("error updating user's premium", "error", err)
		return err
	}

	//Send message to the user saying that payment has been successful
	_, err = b.b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   b.messages.SuccessfulPayment[language(update.Message.From)],
	})
	if err != nil {
		b.logger.Errorw("error sending message", "error", err)
	}
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

// premium returns true if user is still premium
func premium(user *db.User) bool {
	return !time.Unix(user.PremiumUntil, 0).Before(time.Now())
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
