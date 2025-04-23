package tts

import (
	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const georgianAPIEndPoint = "https://api.narakeet.com/text-to-speech/mp3"

// !!! Unofficial API - might break
const tatarAPIEndPoint = "https://issai.nu.edu.kz/tatartts/?speaker=female&text="

type Client struct {
	tts    *texttospeech.Client
	apiKey string
}

// New creates new tts client
func New(ctx context.Context, apiKey string) (*Client, error) {
	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return &Client{client, apiKey}, nil
}

// Close closes tts client
func (c *Client) Close() error {
	if err := c.tts.Close(); err != nil {
		return err
	}
	return nil
}

// Generate generates mp3 audio based on the text and language provided
func (c *Client) Generate(ctx context.Context, text string, languageCode string) ([]byte, error) {
	if languageCode == "ka-GE" {
		return c.generateGeorgian(ctx, text)
	}

	if languageCode == "tatar" {
		return c.generateTatar(ctx, text)
	}
	// Perform the text-to-speech request on the text input with the selected voice parameters and audio file type.
	req := texttospeechpb.SynthesizeSpeechRequest{
		// Set the text input to be synthesized.
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{Text: text},
		},
		// Build the voice request, select the language code (e.g. "en-US") and the SSML voice gender ("neutral").
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: languageCode,
			SsmlGender:   texttospeechpb.SsmlVoiceGender_NEUTRAL,
		},
		// Select the type of audio file you want returned.
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_MP3,
		},
	}

	//Generate speech
	resp, err := c.tts.SynthesizeSpeech(ctx, &req)
	if err != nil {
		log.Fatal(err)
	}
	return resp.AudioContent, nil
}

// Separate function for georgian because Google doesn't have georgian tts
func (c *Client) generateGeorgian(ctx context.Context, text string) ([]byte, error) {
	//Create new request
	req, err := http.NewRequestWithContext(ctx, "POST", georgianAPIEndPoint, strings.NewReader(text))
	if err != nil {
		return nil, err
	}

	//Set headers
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("accept", "application/octet-stream")

	//Make a request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	//Get mp3 data from the request
	audioContent, err := io.ReadAll(resp.Body)
	return audioContent, err
}

// !!! Unofficial API - might break
func (c *Client) generateTatar(ctx context.Context, text string) ([]byte, error) {
	//Create new request
	req, err := http.NewRequestWithContext(ctx, "GET", tatarAPIEndPoint+url.QueryEscape(text), http.NoBody)
	if err != nil {
		return nil, err
	}

	//Make a request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	//Get b64 from the response
	var b64 string
	if err := json.NewDecoder(resp.Body).Decode(&b64); err != nil {
		return nil, err
	}

	//Return decoded mp3 and error
	return base64.StdEncoding.DecodeString(b64)
}
