package v1

import (
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

type userClaims struct {
	jwt.RegisteredClaims
	UserId   uint64 `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

func newUserClaims(id uint64, username, role string, exp time.Time) *userClaims {
	return &userClaims{
		UserId:   id,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{
				Time: exp,
			},
		},
	}
}

func (u *userClaims) GetUserId() uint64 {
	return u.UserId
}

func (u *userClaims) GetUsername() string {
	if u.Username == "" {
		return "Anonymous"
	}
	return u.Username
}

func (u *userClaims) GetRole() string {
	if u.Role == "" {
		return "ROLE_UNSPECIFIED"
	}
	return u.Role
}

func (u *userClaims) GetExp() *jwt.NumericDate {
	return u.ExpiresAt
}

func (u *userClaims) IsAnonymous() bool {
	return u.GetRole() == "ROLE_UNSPECIFIED" ||
		u.GetUsername() == "Anonymous" ||
		u.GetExp() == nil
}
