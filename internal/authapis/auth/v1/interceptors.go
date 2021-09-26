package v1

import (
	"context"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func logUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	userClaim := parseUserOrAnonymousFromCTX(ctx)
	glog.Infof("%v has call %v", userClaim.GetUsername(), info.FullMethod)
	return handler(ctx, req)
}

func authInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {

	// Bypass login
	loginPath := "/auth.v1.AuthenService/Login"
	if info.FullMethod == loginPath {
		return handler(ctx, req)
	}

	// Get metadata from in comming context
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	// Parse authorization
	values, ok := md["authorization"]
	if !ok || len(values) == 0 {
		return ctx, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	accessToken := values[0]
	claims, err := verifyJWTToken(ctx, accessToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}

	// Assign metadata to new context
	newCTX := context.WithValue(ctx, JWTKey, userClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: claims.GetExp(),
		},
		UserId:   claims.GetUserId(),
		Username: claims.GetUsername(),
		Role:     claims.GetRole(),
	})

	return handler(newCTX, req)
}
