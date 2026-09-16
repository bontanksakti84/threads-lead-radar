package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const defaultGroqModel = "openai/gpt-oss-20b"

type GroqClassifier struct {
	apiKey     string
	model      string
	baseURL    string
	httpClient *http.Client
}

func NewGroqClassifier() (*GroqClassifier, error) {
	apiKey := strings.TrimSpace(os.Getenv("GROQ_API_KEY"))

	if apiKey == "" {
		return nil, fmt.Errorf("GROQ_API_KEY is not set")
	}

	model := strings.TrimSpace(os.Getenv("GROQ_MODEL"))
	if model == "" {
		model = defaultGroqModel
	}

	return &GroqClassifier{
		apiKey:  apiKey,
		model:   model,
		baseURL: "https://api.groq.com/openai/v1",
		httpClient: &http.Client{
			Timeout: 60 * 1000 * 1000 * 1000,
		},
	}, nil
}

type groqChatRequest struct {
	Model          string          `json:"model"`
	Messages       []groqMessage   `json:"messages"`
	Temperature    float64         `json:"temperature"`
	ResponseFormat groqResponseFmt `json:"response_format"`
}

type groqMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type groqResponseFmt struct {
	Type       string         `json:"type"`
	JSONSchema groqJSONSchema `json:"json_schema"`
}

type groqJSONSchema struct {
	Name   string         `json:"name"`
	Strict bool           `json:"strict"`
	Schema map[string]any `json:"schema"`
}

type groqChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func (g *GroqClassifier) Classify(
	ctx context.Context,
	content string,
) (Classification, error) {

	content = strings.TrimSpace(content)

	if content == "" {
		return Classification{}, fmt.Errorf("content is empty")
	}

	systemPrompt := `
You are a lead qualification classifier for a software development agency in Indonesia.

Analyze the Threads post and classify ONLY what is supported by the text.

Do not invent budget, urgency, business details, or hiring intent.

Categories:
- website
- mobile_app
- custom_software
- erp
- hrm
- pos
- automation
- ai
- ecommerce
- maintenance
- other

Intents:
- hire_developer
- looking_for_jasa
- project_inquiry
- recommendation
- product_research
- general_question
- other

Lead scoring:
0-29   = unlikely to become a service lead
30-49  = weak potential
50-69  = possible lead
70-84  = strong lead
85-100 = very strong lead

A person asking for recommendations is NOT automatically hiring.

A person asking whether anyone has experience with a product is normally product_research, not a service lead.

Set needs_developer=true only when the post indicates a need for software development, programming, a developer, programmer, or implementation work.

Set has_budget=true only when the post explicitly mentions a budget, price, rate, cost, or spending amount.

Set has_urgency=true only when urgency is explicitly stated.

The summary must be concise and factual.
`

	userPrompt := fmt.Sprintf(
		"Classify this Threads post:\n\n%s",
		content,
	)

	request := groqChatRequest{
		Model: g.model,
		Messages: []groqMessage{
			{
				Role:    "system",
				Content: systemPrompt,
			},
			{
				Role:    "user",
				Content: userPrompt,
			},
		},
		Temperature: 0,
		ResponseFormat: groqResponseFmt{
			Type: "json_schema",
			JSONSchema: groqJSONSchema{
				Name:   "lead_classification",
				Strict: true,
				Schema: classificationSchema(),
			},
		},
	}

	body, err := json.Marshal(request)
	if err != nil {
		return Classification{}, fmt.Errorf(
			"marshal Groq request: %w",
			err,
		)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		g.baseURL+"/chat/completions",
		bytes.NewReader(body),
	)
	if err != nil {
		return Classification{}, fmt.Errorf(
			"create Groq request: %w",
			err,
		)
	}

	req.Header.Set(
		"Authorization",
		"Bearer "+g.apiKey,
	)
	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return Classification{}, fmt.Errorf(
			"call Groq API: %w",
			err,
		)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return Classification{}, fmt.Errorf(
			"read Groq response: %w",
			err,
		)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return Classification{}, fmt.Errorf(
			"Groq API returned %s: %s",
			resp.Status,
			string(responseBody),
		)
	}

	var response groqChatResponse

	if err := json.Unmarshal(
		responseBody,
		&response,
	); err != nil {
		return Classification{}, fmt.Errorf(
			"decode Groq response: %w",
			err,
		)
	}

	if len(response.Choices) == 0 {
		return Classification{}, fmt.Errorf(
			"Groq response contains no choices",
		)
	}

	var result Classification

	if err := json.Unmarshal(
		[]byte(response.Choices[0].Message.Content),
		&result,
	); err != nil {
		return Classification{}, fmt.Errorf(
			"decode classification JSON: %w",
			err,
		)
	}

	if err := ValidateClassification(result); err != nil {
		return Classification{}, fmt.Errorf(
			"invalid classification: %w",
			err,
		)
	}

	return result, nil
}

func classificationSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"is_lead": map[string]any{
				"type": "boolean",
			},
			"category": map[string]any{
				"type": "string",
				"enum": []string{
					"website",
					"mobile_app",
					"custom_software",
					"erp",
					"hrm",
					"pos",
					"automation",
					"ai",
					"ecommerce",
					"maintenance",
					"other",
				},
			},
			"intent": map[string]any{
				"type": "string",
				"enum": []string{
					"hire_developer",
					"looking_for_jasa",
					"project_inquiry",
					"recommendation",
					"product_research",
					"general_question",
					"other",
				},
			},
			"lead_score": map[string]any{
				"type":    "integer",
				"minimum": 0,
				"maximum": 100,
			},
			"has_budget": map[string]any{
				"type": "boolean",
			},
			"has_urgency": map[string]any{
				"type": "boolean",
			},
			"needs_developer": map[string]any{
				"type": "boolean",
			},
			"summary": map[string]any{
				"type": "string",
			},
		},
		"required": []string{
			"is_lead",
			"category",
			"intent",
			"lead_score",
			"has_budget",
			"has_urgency",
			"needs_developer",
			"summary",
		},
		"additionalProperties": false,
	}
}
