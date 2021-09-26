package v1

import (
	"context"

	"github.com/golang/glog"
	authv1 "github.com/tuanden0/learn_ent/proto/gen/go/v1/auth"
)

type Service interface {
	authv1.AuthenServiceServer
}

type service struct {
	repo      Repository
	validator validate
}

func NewService(repo Repository, validator validate) Service {
	return &service{
		repo:      repo,
		validator: validator,
	}
}

func (s *service) Ping(ctx context.Context, in *authv1.PingRequest) (*authv1.PingResponse, error) {
	return &authv1.PingResponse{}, nil
}

func (s *service) Login(ctx context.Context, in *authv1.LoginRequest) (*authv1.LoginResponse, error) {

	if err := s.validator.Login(ctx, in); err != nil {
		glog.Errorf("user login error %v", err)
		return nil, err
	}

	u, err := s.repo.Login(ctx, in)
	if err != nil {
		return nil, err
	}

	token, err := generateJWTToken(ctx, u)
	if err != nil {
		glog.Error(err)
		return nil, errGenerateToken
	}

	return mapLoginResponse(token), nil
}

func (s *service) Verify(ctx context.Context, in *authv1.VerifyRequest) (*authv1.VerifyResponse, error) {

	if err := s.validator.Verify(ctx, in); err != nil {
		glog.Errorf("token verify error %v", err)
		return nil, err
	}

	u, err := verifyJWTToken(ctx, in.GetToken())
	if err != nil {
		glog.Error(err)
		return nil, errVerifyToken
	}

	return mapVerifyResponse(u), nil
}
