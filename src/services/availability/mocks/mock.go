package availabilitymocks

import (
	"context"
	"urlresolver/src/services/availability"
)

type MockAvailabilityService struct {
	// Configurable outputs
	Availability availability.Availability
	Err          error

	// Captured inputs
	Called  bool
	VideoID string
	Token   string
}

func (m *MockAvailabilityService) GetAvailability(
	ctx context.Context,
	videoID string,
	token string,
) (availability.Availability, error) {

	m.Called = true
	m.VideoID = videoID
	m.Token = token

	return m.Availability, m.Err
}
