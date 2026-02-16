package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"telegram-movie-bot/internal/gateway/models"
)

const urlGetRecommendations = "/movies/recommendations?user_id="

func (c *Client) GetRecommendations(ctx context.Context, userID string) ([]models.Movie, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+urlGetRecommendations+userID, nil)
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
			return nil, errors.New("empty response from recommendation service")
		}

		var movies []models.Movie
		if err := json.Unmarshal(body, &movies); err != nil {
			return nil, fmt.Errorf("decode recommendations: %w", err)
		}

		if len(movies) == 0 {
			return nil, fmt.Errorf("no recommendations available for user %s", userID)
		}

		return movies, nil

	case http.StatusBadRequest:
		return nil, errors.New("invalid request: check user_id or query parameters")

	case http.StatusForbidden:
		return nil, errors.New("access denied: user is not allowed to get recommendations")

	case http.StatusNotFound:
		return nil, fmt.Errorf("recommendations not found: user %s has no data", userID)

	case http.StatusTooManyRequests:
		return nil, errors.New("too many requests: the request limit for the Recommendations API has been exceeded")

	case http.StatusInternalServerError:
		return nil, errors.New("recommendation service internal error, please try again later")

	case http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return nil, errors.New("recommendation service unavailable: please try again later")

	default:
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}
}
