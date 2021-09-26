package v1

import (
	authv1 "github.com/tuanden0/learn_ent/proto/gen/go/v1/auth"
)

func mapLoginResponse(token string) *authv1.LoginResponse {
	return &authv1.LoginResponse{
		AccessToken: token,
	}
}

func mapVerifyResponse(u *userClaims) *authv1.VerifyResponse {
	return &authv1.VerifyResponse{
		Id:       u.GetUserId(),
		Username: u.GetUsername(),
		Role:     u.GetRole(),
		Exp:      u.GetExp().Unix(),
	}
}
