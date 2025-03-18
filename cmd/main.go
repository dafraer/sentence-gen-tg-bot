package main

import (
	"context"
	"github.com/dafraer/sentence-gen-tg-bot/text"
	"github.com/dafraer/sentence-gen-tg-bot/tts"
	"os"

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
	ctx := context.TODO()

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

	//Create bot
	myBot, err := bot.New(token, store, geminiClient, ttsClient, msgs)
	if err != nil {
		panic(err)
	}

	if err := myBot.Run(ctx); err != nil {
		panic(err)
	}

}
