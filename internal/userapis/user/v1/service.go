package v1

import (
	"context"

	"github.com/golang/glog"
	userv1 "github.com/tuanden0/learn_ent/proto/gen/go/v1/user"
)

type Service interface {
	userv1.UserServiceServer
	userv1.UserAuthenServiceServer
}

type service struct {
	userv1.UnimplementedUserServiceServer
	repo      Repository
	validator validate
}

func NewService(repo Repository, validator validate) Service {
	return &service{
		repo:      repo,
		validator: validator,
	}
}

func (s *service) Ping(ctx context.Context, in *userv1.PingRequest) (*userv1.PingResponse, error) {
	return &userv1.PingResponse{}, nil
}

func (s *service) Create(ctx context.Context, in *userv1.CreateRequest) (*userv1.CreateResponse, error) {

	if err := s.validator.Create(ctx, in); err != nil {
		glog.Errorf("user create input error %v", err)
		return nil, err
	}

	u, err := s.repo.Create(ctx, in)
	if err != nil {
		return nil, err
	}

	return mapCreateResponse(u), nil
}

func (s *service) Retrieve(ctx context.Context, in *userv1.RetrieveRequest) (*userv1.RetrieveResponse, error) {

	if err := s.validator.Retrieve(ctx, in); err != nil {
		glog.Errorf("user retrieve input error %v", err)
		return nil, err
	}

	u, err := s.repo.Retrieve(ctx, in.GetId())
	if err != nil {
		return nil, err
	}

	return mapRetrieveResponse(u), nil
}

func (s *service) Update(ctx context.Context, in *userv1.UpdateRequest) (*userv1.UpdateResponse, error) {

	if err := s.validator.Update(ctx, in); err != nil {
		glog.Errorf("user update input error %v", err)
		return nil, err
	}

	u, err := s.repo.Update(ctx, in)
	if err != nil {
		return nil, err
	}

	return mapUpdateResponse(u), nil
}

func (s *service) Delete(ctx context.Context, in *userv1.DeleteRequest) (*userv1.DeleteResponse, error) {

	if err := s.validator.Delete(ctx, in); err != nil {
		glog.Errorf("user delete input error %v", err)
		return nil, err
	}

	err := s.repo.Delete(ctx, in.GetId())
	if err != nil {
		return nil, err
	}

	return &userv1.DeleteResponse{}, nil
}

func (s *service) List(ctx context.Context, in *userv1.ListRequest) (*userv1.ListResponse, error) {

	if err := s.validator.List(ctx, in); err != nil {
		glog.Errorf("user list input error %v", err)
		return nil, err
	}

	us, err := s.repo.List(ctx, in)
	if err != nil {
		return nil, err
	}

	return mapListResponse(us), nil
}

func (s *service) Login(ctx context.Context, in *userv1.UserLoginRequest) (*userv1.UserLoginResponse, error) {

	if err := s.validator.Login(ctx, in); err != nil {
		glog.Errorf("user login input error %v", err)
		return nil, err
	}

	u, err := s.repo.Login(ctx, in)
	if err != nil {
		glog.Error(err)
		return nil, errLogin
	}

	return mapUserLoginResponse(u), nil
}
