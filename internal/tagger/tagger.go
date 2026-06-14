package tagger

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kamuridesu/readeck-manager/internal/config"
	"github.com/kamuridesu/readeck-manager/internal/readeck"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/shared"
)

type Tagger struct {
	Client *openai.Client
	Model  string
}

type AIResult struct {
	Tags []string `json:"tags"`
}

func New(cfg *config.Config) *Tagger {
	opts := []option.RequestOption{
		option.WithAPIKey(cfg.OpenAIKey),
	}

	if cfg.OpenAIUrl != "" {
		opts = append(opts, option.WithBaseURL(cfg.OpenAIUrl))
	}

	client := openai.NewClient(opts...)

	return &Tagger{
		Client: &client,
		Model:  cfg.OpenAIModel,
	}
}

func (o *Tagger) GenerateLabels(ctx context.Context, bm *readeck.Bookmark, labels []string, snippet string) ([]string, error) {
	prompt := fmt.Sprintf(TAGGER_PROMPT, bm.Title, bm.SiteName, bm.Description, snippet, strings.Join(labels, ", "))
	prompt += "\n\nYou MUST respond with a valid JSON object in this exact format:\n{\"tags\": [\"label1\", \"label2\"]}"

	format := shared.NewResponseFormatJSONObjectParam()

	chatCompletion, err := o.Client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{

		Model: shared.ChatModel(o.Model),
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONObject: &format,
		},

		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("You are a precise classification engine. You must output ONLY a JSON object. Do not output markdown, reasoning, or bullet points. Output nothing but JSON."),
			openai.UserMessage(prompt),
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get AI response: %w", err)
	}

	if len(chatCompletion.Choices) == 0 {
		return nil, fmt.Errorf("ai returned no completion choices")
	}

	responseText := chatCompletion.Choices[0].Message.Content

	cleanResponse := strings.TrimSpace(responseText)
	cleanResponse = strings.TrimPrefix(cleanResponse, "```json")
	cleanResponse = strings.TrimPrefix(cleanResponse, "```")
	cleanResponse = strings.TrimSuffix(cleanResponse, "```")
	cleanResponse = strings.TrimSpace(cleanResponse)

	var result AIResult
	if err := json.Unmarshal([]byte(cleanResponse), &result); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w, content is: %s", err, responseText)
	}

	var validLabels []string
	for _, label := range result.Tags {
		for _, available := range labels {
			if strings.EqualFold(strings.TrimSpace(label), available) {
				validLabels = append(validLabels, available)
				break
			}
		}
	}

	return validLabels, nil
}
