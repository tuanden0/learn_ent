package v1

import (
	"context"
	"time"

	authv1 "github.com/tuanden0/learn_ent/proto/gen/go/v1/auth"
	userv1 "github.com/tuanden0/learn_ent/proto/gen/go/v1/user"
	"google.golang.org/grpc"
)

type Repository interface {
	Login(ctx context.Context, in *authv1.LoginRequest) (*userv1.UserLoginResponse, error)
}

type repoManager struct{}

func NewRepoManager() Repository {
	return &repoManager{}
}

func (r *repoManager) Login(ctx context.Context, in *authv1.LoginRequest) (*userv1.UserLoginResponse, error) {

	// Create connection
	ctx, cancer := context.WithTimeout(ctx, 5*time.Second)
	defer cancer()

	cc, err := grpc.DialContext(ctx, userServiceUrl, setupClientDialOpts()...)
	if err != nil {
		return nil, err
	}
	defer cc.Close()

	// Create client connection
	c := userv1.NewUserAuthenServiceClient(cc)

	// Call client
	res, err := c.Login(ctx, &userv1.UserLoginRequest{
		Username: in.GetUsername(),
		Password: in.GetPassword(),
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}
