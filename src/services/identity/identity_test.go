package identity

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsPremium(t *testing.T) {
	testCases := []struct {
		description     string
		input           UserInfo
		expectedOutcome bool
	}{
		{
			description: "Happy path should return true when user has premium role",
			input: UserInfo{
				ID:    1,
				Name:  "Joe",
				Email: "joe@test.com",
				Roles: []string{"premium"},
			},
			expectedOutcome: true,
		},
		{
			description: "Should return false when user has no premium role",
			input: UserInfo{
				ID:    1,
				Name:  "Joe",
				Email: "joe@test.com",
				Roles: []string{"standard"},
			},
			expectedOutcome: false,
		},
		{
			description: "Should return false when user has no roles",
			input: UserInfo{
				ID:    1,
				Name:  "Joe",
				Email: "joe@test.com",
				Roles: []string{},
			},
			expectedOutcome: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			got := tc.input.IsPremium()
			assert.Equal(t, tc.expectedOutcome, got)
		})
	}
}

func TestGetUser(t *testing.T) {

	testCases := []struct {
		description     string
		responseBody    string
		statusCode      int
		expectedOutcome UserInfo
		expectError     bool
	}{
		{
			description: "Happy path should return user when valid response received",
			statusCode:  http.StatusOK,
			responseBody: `{
				"id": 7564,
				"name": "Joe Bloggs",
				"email": "joe.bloggs@noemail.notreal",
				"roles": ["premium"]
			}`,
			expectedOutcome: UserInfo{
				ID:    7564,
				Name:  "Joe Bloggs",
				Email: "joe.bloggs@noemail.notreal",
				Roles: []string{"premium"},
			},
			expectError: false,
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
			description: "Sad path should error when non 200 status code received",
			statusCode:  http.StatusUnauthorized,
			expectError: true,
		},
		{
			description: "Sad path should error when missing required id field",
			statusCode:  http.StatusOK,
			responseBody: `{
				"id": 0,
				"name": "Joe Bloggs",
				"email": "joe.bloggs@noemail.notreal",
				"roles": ["premium"]
			}`,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {

			// Arrange: fake downstream identity service
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.statusCode)
				w.Write([]byte(tc.responseBody))
			}))
			defer server.Close()

			httpClient := server.Client()

			svc := NewService(server.URL, httpClient)

			// Act
			user, err := svc.GetUserInfo(
				context.Background(),
				"fake-token",
			)

			// Assert error cases
			if tc.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedOutcome, user)
		})
	}
}
