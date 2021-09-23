package v1

import (
	authv1 "github.com/tuanden0/learn_ent/proto/gen/go/v1/auth"
)

func mapLoginResponse() *authv1.LoginResponse {
	return &authv1.LoginResponse{
		AccessToken: "not implement JWT function yet",
	}
}
