package availability

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Availability struct {
	VideoID int
	From    time.Time
	To      time.Time
}

type AvailabilityService struct {
	client  *http.Client
	baseURL string
}

func NewService(baseURL string, client *http.Client) *AvailabilityService {
	return &AvailabilityService{
		client:  client,
		baseURL: baseURL,
	}
}

func (a Availability) IsValid(now time.Time) bool {
	return ((now.Equal(a.From) || now.After(a.From)) &&
		(now.Equal(a.To) || now.Before(a.To)))
}

func (as *AvailabilityService) GetAvailability(ctx context.Context, videoID, token string) (Availability, error) {

	url := fmt.Sprintf("%s/availabilityinfo/%s", as.baseURL, videoID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Availability{}, errors.New(fmt.Sprintf("failed to build availability http request %v", err))
	}

	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", token))
	resp, err := as.client.Do(req)
	if err != nil {
		return Availability{}, errors.New(fmt.Sprintf("error calling availability service %v", err))
	}

	if resp.StatusCode != http.StatusOK {
		return Availability{}, errors.New(fmt.Sprintf("availability service returned status %v", resp.StatusCode))
	}

	apiResp := AvailabilityResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return Availability{}, errors.New(fmt.Sprint("error decoding availability api response %w", err))
	}

	avail, err := apiResp.MapToDomain()
	if err != nil {
		return Availability{}, fmt.Errorf("error validating / building domain availability %w", err)
	}

	return avail, nil
}
