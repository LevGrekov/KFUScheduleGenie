package kfuapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const baseURL = "https://auth.kpfu.tyuop.ru/api/v1"

type Client struct {
	HTTPClient *http.Client
}

func NewClient() *Client {
	return &Client{
		HTTPClient: &http.Client{},
	}
}

func (c *Client) doRequest(req *http.Request, v any) error {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("неожиданный статус: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("ошибка чтения ответа: %w", err)
	}

	if err := json.Unmarshal(body, v); err != nil {
		return fmt.Errorf("ошибка парсинга JSON: %w", err)
	}

	return nil
}
