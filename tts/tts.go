package tts

import (
	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
	"context"
	"log"
)

type Client struct {
	tts *texttospeech.Client
}

func New(ctx context.Context) (*Client, error) {
	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return &Client{client}, nil
}

func (c *Client) Close() error {
	if err := c.tts.Close(); err != nil {
		return err
	}
	return nil
}

func (c *Client) Generate(ctx context.Context, text string, languageCode string) (*texttospeechpb.SynthesizeSpeechResponse, error) {
	// Perform the text-to-speech request on the text input with the selected
	// voice parameters and audio file type.
	req := texttospeechpb.SynthesizeSpeechRequest{
		// Set the text input to be synthesized.
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{Text: text},
		},
		// Build the voice request, select the language code ("en-US") and the SSML
		// voice gender ("neutral").
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: languageCode,
			SsmlGender:   texttospeechpb.SsmlVoiceGender_NEUTRAL,
		},
		// Select the type of audio file you want returned.
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_MP3,
		},
	}

	resp, err := c.tts.SynthesizeSpeech(ctx, &req)
	if err != nil {
		log.Fatal(err)
	}
	return resp, nil
}
