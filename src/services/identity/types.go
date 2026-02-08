package identity

import (
	"fmt"
	"slices"
)

type IdentityResponse struct {
	ID    int      `json:"id"`
	Name  string   `json:"name"`
	Email string   `json:"email"`
	Roles []string `json:"roles"`
}

type UserInfo struct {
	ID    int
	Name  string
	Email string
	Roles []string
}

func (u UserInfo) IsPremium() bool {
	return slices.Contains(u.Roles, "premium")
}

func (r IdentityResponse) ToDomain() (UserInfo, error) {
	if r.ID == 0 {
		return UserInfo{}, fmt.Errorf("invalid identity response: missing id")
	}

	return UserInfo{
		ID:    r.ID,
		Name:  r.Name,
		Email: r.Email,
		Roles: r.Roles,
	}, nil
}
