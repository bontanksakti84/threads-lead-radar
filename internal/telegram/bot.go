package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Bot struct {
	token  string
	chatID string
	client *http.Client
}

func NewBot(token, chatID string) *Bot {
	return &Bot{
		token:  token,
		chatID: chatID,
		client: &http.Client{},
	}
}

type sendMessageRequest struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode,omitempty"`
}

type telegramResponse struct {
	OK          bool   `json:"ok"`
	Description string `json:"description"`
}

func (b *Bot) SendMessage(ctx context.Context, text string) error {
	if b.token == "" {
		return fmt.Errorf("telegram bot token is empty")
	}

	if b.chatID == "" {
		return fmt.Errorf("telegram chat ID is empty")
	}

	payload := sendMessageRequest{
		ChatID:    b.chatID,
		Text:      text,
		ParseMode: "HTML",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal telegram request: %w", err)
	}

	url := fmt.Sprintf(
		"https://api.telegram.org/bot%s/sendMessage",
		b.token,
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewReader(body),
	)
	if err != nil {
		return fmt.Errorf("create telegram request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := b.client.Do(req)
	if err != nil {
		return fmt.Errorf("send telegram request: %w", err)
	}
	defer resp.Body.Close()

	var result telegramResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("decode telegram response: %w", err)
	}

	if !result.OK {
		return fmt.Errorf(
			"telegram API error: %s",
			result.Description,
		)
	}

	return nil
}
