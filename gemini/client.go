package gemini

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

const (
	geminiFlash   = "gemini-1.5-flash"
	requestString = "Give me a simple %s sentence in %s with the word %s. Don't explain me anything just give me the sentence and translation to %s seperated by ; symbol"
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

// Request sends a text-only request to the gemini-1.5-flash model
func (c *Client) Request(ctx context.Context, request string) (string, error) {
	//Specify model
	model := c.client.GenerativeModel(geminiFlash)

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

func FormatRequestString(level, sentenceLanguage, word, language string) string {
	return fmt.Sprintf(requestString, level, sentenceLanguage, word, language)
}
