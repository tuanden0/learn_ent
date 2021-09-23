package v1

import (
	"context"

	authv1 "github.com/tuanden0/learn_ent/proto/gen/go/v1/auth"
)

type Service interface {
	authv1.AuthenServiceServer
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Ping(ctx context.Context, in *authv1.PingRequest) (*authv1.PingResponse, error) {
	return &authv1.PingResponse{}, nil
}

func (s *service) Login(ctx context.Context, in *authv1.LoginRequest) (*authv1.LoginResponse, error) {

	u, err := s.repo.Login(ctx, in)
	if err != nil {
		return nil, err
	}

	_ = u

	return mapLoginResponse(), nil
}
