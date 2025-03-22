package main

import (
	"context"
	"github.com/dafraer/sentence-gen-tg-bot/text"
	"github.com/dafraer/sentence-gen-tg-bot/tts"
	"go.uber.org/zap"
	"os"
	"os/signal"

	"github.com/dafraer/sentence-gen-tg-bot/bot"
	"github.com/dafraer/sentence-gen-tg-bot/db"
	"github.com/dafraer/sentence-gen-tg-bot/gemini"
)

func main() {
	if len(os.Args) < 3 {
		panic("telegram bot token and gemini API key must be passed as arguments")
	}
	token := os.Args[1]
	geminiAPIKey := os.Args[2]

	//Declare context that is marked Done when os.Interrupt is called
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	//Create firestore client
	store, err := db.New(ctx)
	if err != nil {
		panic(err)
	}

	//Create gemini client
	geminiClient, err := gemini.New(ctx, geminiAPIKey)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := geminiClient.Close(); err != nil {
			panic(err)
		}
	}()

	//Create tts client
	ttsClient, err := tts.New(ctx)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := ttsClient.Close(); err != nil {
			panic(err)
		}
	}()

	//Load messages
	msgs := text.Load()

	//Create logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	sugar := logger.Sugar()

	//Create bot
	myBot, err := bot.New(token, store, geminiClient, ttsClient, msgs, sugar)
	if err != nil {
		panic(err)
	}

	//If webhook flag specified run bot using webhook
	webhook := len(os.Args) == 4 && os.Args[3] == "-w"
	if webhook {
		if err := myBot.RunWebhook(ctx, ":8080"); err != nil {
			panic(err)
		}
		return
	}
	myBot.Run(ctx)
}
