package main

import (
	"context"
	"fmt"
	"os"

	"github.com/dafraer/sentence-gen-tg-bot/gemini"
)

func main() {
	if len(os.Args) < 2 {
		panic("API key must be passed aas an argument")
	}
	ctx := context.TODO()
	client, err := gemini.New(ctx, os.Args[1])
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			panic(err)
		}
	}()

	resp, err := client.Request(ctx, "Give me a simple A1 sentence in russian with the word сосать. dont explain me anything just give me the sentence and trsnalation seperated by ; symbol")
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
