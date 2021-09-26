package v1

import (
	"context"
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
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

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock(),
	}
	cc, err := grpc.DialContext(ctx, userServiceUrl, opts...)
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

func generateJWTToken(ctx context.Context, in *userv1.UserLoginResponse) (string, error) {

	claims := newUserClaims(
		in.GetId(),
		in.GetUsername(),
		in.GetRole(),
		time.Now().Add(tokenDuration),
	)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(secret))

	return ss, err
}

func verifyJWTToken(ctx context.Context, token string) (*userClaims, error) {

	tk, err := jwt.ParseWithClaims(
		token,
		&userClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}
			return []byte(secret), nil
		},
	)

	if err != nil {
		return nil, err
	}

	err = tk.Claims.Valid()
	if err != nil {
		return nil, err
	}

	claims, ok := tk.Claims.(*userClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

func parseUserOrAnonymousFromCTX(ctx context.Context) *userClaims {

	userClaim := ctx.Value(JWTKey)
	if userClaim != nil {
		if userClaim, ok := userClaim.(userClaims); ok {
			return &userClaim
		}
	}

	return &userClaims{
		Username: "Anonymous",
		Role:     "ROLE_UNSPECIFIED",
	}
}

func parseUsersOrNilFromCTX(ctx context.Context) *userClaims {

	userClaim := ctx.Value(JWTKey)
	if userClaim != nil {
		if userClaim, ok := userClaim.(userClaims); ok {
			return &userClaim
		}
	}

	return nil
}
