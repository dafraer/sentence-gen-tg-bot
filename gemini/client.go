package gemini

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

const (
	requestStringEn = `
Generate a simple %s-level sentence in %s using the word %s.  
- The sentence should make it easy to understand the word from context.  
- If the word doesn't exist or if it is from another language, return only "Error"
- Otherwise, return the sentence and its English translation, separated by ";".  
- Do not include any explanations or extra text."`
	requestStringRu = `
Generate a simple %s-level sentence in %s using the word %s.  
- The sentence should make it easy to understand the word from context.  
- If the word doesn't exist or if it is from another language, return only "Error"
- Otherwise, return the sentence and its Russian translation, separated by ";".  
- Do not include any explanations or extra text."`
)

type Client struct {
	client *genai.Client
}

// New creates new gemini client
func New(ctx context.Context, token string) (*Client, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(token))
	if err != nil {
		return nil, err
	}
	return &Client{client: client}, nil
}

// Close closes the client
func (c *Client) Close() error {
	if err := c.client.Close(); err != nil {
		return err
	}
	return nil
}

// Request sends a text-only request to the gemini model specified in geminiVersion parameter
func (c *Client) Request(ctx context.Context, request string, geminiVersion string) (string, error) {
	//Specify model
	model := c.client.GenerativeModel(geminiVersion)

	//Generate content
	resp, err := model.GenerateContent(ctx, genai.Text(request))
	if err != nil {
		return "", err
	}

	//Extract response
	var response strings.Builder
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				response.WriteString(string(part.(genai.Text)))
			}
		}
	}

	return response.String(), nil
}

// FormatRequestString formats request string based on the language
func FormatRequestString(level, sentenceLanguage, word, language string) string {
	if language == "ru" {
		return fmt.Sprintf(requestStringRu, level, sentenceLanguage, word)
	}
	return fmt.Sprintf(requestStringEn, level, sentenceLanguage, word)
}
