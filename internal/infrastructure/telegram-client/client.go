package telegram_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"tlgbs/internal/infrastructure/models"
)

type TgClient struct {
	token      string
	httpClient *http.Client
}

func NewTgClient(token string) *TgClient {
	return &TgClient{
		token:      token,
		httpClient: &http.Client{},
	}
}

func (c *TgClient) SendMessage(chatID int64, text string) error {

	body, err := json.Marshal(models.SendMsgRequest{
		ChatID: chatID,
		Text:   text,
	})
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	urlApiTelegramBot := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", c.token)
	resp, err := c.httpClient.Post(urlApiTelegramBot, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}

	bodyBytes, _ := io.ReadAll(resp.Body)

	switch resp.StatusCode {
	case http.StatusOK:
		// всё хорошо
		return nil

	case http.StatusBadRequest:
		return fmt.Errorf("telegram API: возможно, неверные параметры: %s", string(bodyBytes))

	case http.StatusUnauthorized:
		return fmt.Errorf("telegram API: неверный токен бота")

	case http.StatusForbidden:
		return fmt.Errorf("telegram API: бот не имеет доступа к чату %d", chatID)

	case http.StatusTooManyRequests:
		return fmt.Errorf("telegram API: превышен лимит запросов")

	case http.StatusInternalServerError:
		return fmt.Errorf("telegram API: ошибка на стороне Telegram")

	default:
		return fmt.Errorf("telegram API: unexpected status %d — response: %s", resp.StatusCode, string(bodyBytes))
	}
}
