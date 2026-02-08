package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"urlresolver/src/services/availability"
	availabilitymocks "urlresolver/src/services/availability/mocks"
	"urlresolver/src/services/identity"
	identitymocks "urlresolver/src/services/identity/mocks"

	"github.com/stretchr/testify/assert"
)

func TestGetVideo(t *testing.T) {

	now := time.Now()

	testCases := []struct {
		description string

		videoID    string
		authHeader string

		mockUser         identity.UserInfo
		mockUserErr      error
		mockAvailability availability.Availability
		mockAvailErr     error

		expectedStatus            int
		expectedIdentityInput     string
		expectedAvailabilityToken string
		expectedAvailabilityID    string
	}{
		{
			description: "Happy path premium user",
			videoID:     "001",
			authHeader:  "bearer testtoken",

			mockUser: identity.UserInfo{
				Roles: []string{"premium"},
			},
			mockAvailability: availability.Availability{
				From: now.Add(-time.Hour),
				To:   now.Add(time.Hour),
			},
			expectedIdentityInput:     "testtoken",
			expectedAvailabilityToken: "testtoken",
			expectedAvailabilityID:    "001",
			expectedStatus:            http.StatusOK,
		},
		{
			description:    "Unauthorized when header missing",
			videoID:        "001",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			description:           "Identity service failure",
			videoID:               "001",
			authHeader:            "bearer token",
			expectedIdentityInput: "token",
			mockUserErr:           assert.AnError,
			expectedStatus:        http.StatusBadGateway,
		},
		{
			description:               "Availability service failure",
			videoID:                   "001",
			authHeader:                "bearer token",
			mockUser:                  identity.UserInfo{},
			mockAvailErr:              assert.AnError,
			expectedStatus:            http.StatusBadGateway,
			expectedIdentityInput:     "token",
			expectedAvailabilityToken: "token",
			expectedAvailabilityID:    "001",
		},
		{
			description: "Video not available",
			videoID:     "001",
			authHeader:  "bearer token",
			mockUser:    identity.UserInfo{},
			mockAvailability: availability.Availability{
				From: now.Add(time.Hour),
				To:   now.Add(2 * time.Hour),
			},
			expectedStatus:            http.StatusForbidden,
			expectedIdentityInput:     "token",
			expectedAvailabilityToken: "token",
			expectedAvailabilityID:    "001",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {

			mockIdentity := &identitymocks.MockIdentityService{
				User: tc.mockUser,
				Err:  tc.mockUserErr,
			}

			mockAvailability := &availabilitymocks.MockAvailabilityService{
				Availability: tc.mockAvailability,
				Err:          tc.mockAvailErr,
			}

			handler := &Handler{
				identitySvc:     mockIdentity,
				availabilitySvc: mockAvailability,
				playbackBaseURL: "https://cdn.test",
			}

			req := httptest.NewRequest(
				http.MethodGet,
				"/video/"+tc.videoID,
				nil,
			)

			if tc.authHeader != "" {
				req.Header.Set("Authorization", tc.authHeader)
			}

			rec := httptest.NewRecorder()

			handler.GetVideo(rec, req)
			assert.Equal(t, tc.expectedIdentityInput, mockIdentity.Input)
			assert.Equal(t, tc.expectedAvailabilityToken, mockAvailability.Token)
			assert.Equal(t, tc.expectedAvailabilityID, mockAvailability.VideoID)
			assert.Equal(t, tc.expectedStatus, rec.Code)
		})
	}
}
