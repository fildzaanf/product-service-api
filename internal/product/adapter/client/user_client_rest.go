package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"product-service-api/internal/product/application/port"
	"product-service-api/pkg/middleware"
)

type userRESTClient struct {
	baseURL string
	client  *http.Client
}

func NewUserRESTClient(host, port string) port.UserQueryClientInterface {
	return &userRESTClient{
		baseURL: fmt.Sprintf("http://%s:%s", host, port), 
		client: &http.Client{
			Timeout: 3 * time.Second, 
		},
	}
}
func (c *userRESTClient) GetUserByID(ctx context.Context, userID string) error {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/users/%s", c.baseURL, userID),
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed create request: %w", err)
	}

	token := ""
	if t, ok := ctx.Value(middleware.ClaimTokenJWT).(string); ok && t != "" {
		token = t
	}


	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token)) 
	}

	req.Header.Set("Content-Type", "application/json")


	response, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer response.Body.Close()


	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed read response body: %w", err)
	}

	var userResponse struct {
		Meta struct {
			Success bool   `json:"success"`
			Message string `json:"message"`
		} `json:"meta"`
		Results interface{} `json:"results"`
	}
	if err := json.Unmarshal(body, &userResponse); err != nil {
		return fmt.Errorf("failed unmarshal response: %w", err)
	}

	if !userResponse.Meta.Success {
		return fmt.Errorf("user not valid: %s", userResponse.Meta.Message)
	}

	return nil
}
