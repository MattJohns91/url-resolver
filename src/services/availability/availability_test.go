package availability

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsValid(t *testing.T) {
	today := time.Date(
		2026, // year
		time.February,
		8,            // day
		14, 30, 0, 0, // hour, min, sec, nanosec
		time.UTC,
	)
	testCases := []struct {
		description     string
		input           Availability
		time            time.Time
		expectedOutcome bool
	}{
		{
			description: "Happy path should return true when valid",
			input: Availability{
				VideoID: 1,
				From:    today.AddDate(0, 0, -3),
				To:      today.AddDate(0, 0, 3),
			},
			time:            today,
			expectedOutcome: true,
		},
		{
			description: "Happy path should return false when invalid",
			input: Availability{
				VideoID: 1,
				From:    today.AddDate(0, 0, 7),
				To:      today.AddDate(0, 0, 3),
			},
			time:            today,
			expectedOutcome: false,
		},
	}

	for _, tc := range testCases {

		got := tc.input.IsValid(tc.time)
		assert.Equal(t, tc.expectedOutcome, got)

	}

}

func TestGetAvailability(t *testing.T) {

	testCases := []struct {
		description     string
		responseBody    string
		statusCode      int
		expectedOutcome Availability
		expectError     bool
	}{
		{
			description: "Happy path should return availability when valid response received",
			statusCode:  http.StatusOK,
			responseBody: `{
				"video_id": 1,
				"availability_window": {
					"from": "2006-01-02T15:04:05.000",
					"to": "2006-06-02T15:04:05.000"
				}
			}`,
			expectedOutcome: Availability{
				VideoID: 1,
				From:    time.Date(2006, time.January, 2, 15, 4, 5, 0, time.UTC),
				To:      time.Date(2006, time.June, 2, 15, 4, 5, 0, time.UTC),
			},
			expectError: false,
		},
		{
			description: "Sad path should error when valid received but invalid date format",
			statusCode:  http.StatusOK,
			responseBody: `{
				"video_id": 1,
				"availability_window": {
					"from": "2006-01-02",
					"to": "2006-06-02"
				}
			}`,
			expectError: true,
		},
		{
			description: "Sad path should error when invalid JSON received",
			statusCode:  http.StatusOK,
			responseBody: `{
				9
			}`,
			expectError: true,
		},
		{
			description: "Sad path should error when  non 200 status code received",
			statusCode:  http.StatusUnauthorized,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {

			// Arrange: fake downstream service
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.statusCode)
				w.Write([]byte(tc.responseBody))
			}))
			defer server.Close()

			httpClient := server.Client()

			svc := NewService(server.URL, httpClient)

			// Act
			availability, err := svc.GetAvailability(
				context.Background(),
				"1",
				"fake-token",
			)

			// Assert error cases
			if tc.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			assert.Equal(t, tc.expectedOutcome, availability)
		})
	}
}
