package identititymocks

import (
	"context"
	"urlresolver/src/services/identity"
)

type MockIdentityService struct {
	User identity.UserInfo
	Err  error

	Input string
}

func (m *MockIdentityService) GetUserInfo(ctx context.Context, token string) (identity.UserInfo, error) {
	m.Input = token
	return m.User, m.Err
}
