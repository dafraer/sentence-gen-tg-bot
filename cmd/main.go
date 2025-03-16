package main

import (
	"context"
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

	store, err := db.New(ctx)
	if err != nil {
		panic(err)
	}
	geminiClient, err := gemini.New(ctx, geminiAPIKey)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := geminiClient.Close(); err != nil {
			panic(err)
		}
	}()

	myBot, err := bot.New(token, store, geminiClient)
	if err != nil {
		panic(err)
	}

	if err := myBot.Run(ctx); err != nil {
		panic(err)
	}

}
