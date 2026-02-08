package identity

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Service struct {
	baseURL string
	client  *http.Client
}

func NewService(baseURL string, client *http.Client) *Service {
	return &Service{
		baseURL: baseURL,
		client:  client,
	}
}

func (s *Service) GetUserInfo(ctx context.Context, token string) (UserInfo, error) {

	url := fmt.Sprintf("%s/identity/userinfo", s.baseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return UserInfo{}, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "bearer "+token)

	resp, err := s.client.Do(req)
	if err != nil {
		return UserInfo{}, fmt.Errorf("call identity service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return UserInfo{}, fmt.Errorf(
			"identity service returned status %d",
			resp.StatusCode,
		)
	}

	var apiResp IdentityResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return UserInfo{}, fmt.Errorf("decode response: %w", err)
	}

	return apiResp.ToDomain()
}
