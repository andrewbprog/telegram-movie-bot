package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
)

const urlAddRating = "/movies/rating"

type RatingRequest struct {
	UserID  string  `json:"user_id"`
	MovieID int64   `json:"movie_id"`
	Rating  float64 `json:"rating"`
}

func (c *Client) AddRating(ctx context.Context, userID string, movieID int64, rating float64, logger *zap.Logger) error {

	body := RatingRequest{
		UserID:  userID,
		MovieID: movieID,
		Rating:  rating,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal body error: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+urlAddRating, bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}

		switch resp.StatusCode {
		case http.StatusBadRequest:
			return fmt.Errorf("invalid rating data: %s", string(body))

		case http.StatusUnauthorized:
			return errors.New("unauthorized: please check authentication credentials")

		case http.StatusForbidden:
			return errors.New("access denied: user is not allowed to rate this movie")

		case http.StatusNotFound:
			return fmt.Errorf("movie not found: cannot rate non-existing movie (id=%d)", movieID)

		case http.StatusConflict:
			return errors.New("rating already exists for this movie")

		case http.StatusInternalServerError:
			return errors.New("movie service internal error, please try again later")

		case http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
			return errors.New("movie service unavailable: please try again later")

		default:
			return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
		}
	}
	logger.Info("rating successfully sent to gateway",
		zap.String("user_id", userID),
		zap.Int64("movie_id", movieID),
		zap.Float64("rating", rating),
	)
	return nil
}
