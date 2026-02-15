package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"tlgbs/internal/gateway/models"
)

const urlGetRandomMovie = "/movies/random?user_id="

func (c *Client) GetRandomMovie(ctx context.Context, userID string) (*models.Movie, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+urlGetRandomMovie+userID, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		if len(body) == 0 {
			return nil, errors.New("empty response from movie service")
		}

		var movie models.Movie
		if err := json.Unmarshal(body, &movie); err != nil {
			return nil, fmt.Errorf("failed to decode movie response: %w", err)
		}
		return &movie, nil

	case http.StatusBadRequest:
		return nil, errors.New("invalid request: check user_id or request parameters")

	case http.StatusForbidden:
		return nil, errors.New("access denied: user is not allowed to get random movies")

	case http.StatusNotFound:
		return nil, fmt.Errorf("movie not found: no random movie available for user %s", userID)

	case http.StatusInternalServerError:
		return nil, errors.New("movie service internal error, please try again later")

	case http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return nil, errors.New("movie service unavailable: please try again later")

	default:
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}
}
