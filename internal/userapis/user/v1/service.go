package v1

import (
	"context"

	userv1 "github.com/tuanden0/learn_ent/proto/gen/go/v1/user"
)

type Service interface {
	userv1.UserServiceServer
}

type service struct {
	userv1.UnimplementedUserServiceServer
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) Ping(ctx context.Context, in *userv1.PingRequest) (*userv1.PingResponse, error) {
	return &userv1.PingResponse{}, nil
}

func (s *service) Create(ctx context.Context, in *userv1.CreateRequest) (*userv1.CreateResponse, error) {

	u, err := s.repo.Create(ctx, in)
	if err != nil {
		return nil, err
	}

	return mapCreateResponse(u), nil
}

func (s *service) Retrieve(ctx context.Context, in *userv1.RetrieveRequest) (*userv1.RetrieveResponse, error) {

	u, err := s.repo.Retrieve(ctx, in.GetId())
	if err != nil {
		return nil, err
	}

	return mapRetrieveResponse(u), nil
}

func (s *service) Update(ctx context.Context, in *userv1.UpdateRequest) (*userv1.UpdateResponse, error) {

	u, err := s.repo.Update(ctx, in)
	if err != nil {
		return nil, err
	}

	return mapUpdateResponse(u), nil
}

func (s *service) Delete(ctx context.Context, in *userv1.DeleteRequest) (*userv1.DeleteResponse, error) {

	err := s.repo.Delete(ctx, in.GetId())
	if err != nil {
		return nil, err
	}

	return &userv1.DeleteResponse{}, nil
}

func (s *service) List(ctx context.Context, in *userv1.ListRequest) (*userv1.ListResponse, error) {

	us, err := s.repo.List(ctx, in)
	if err != nil {
		return nil, err
	}

	return mapListResponse(us), nil
}
