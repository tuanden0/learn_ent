package v1

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	errGenerateToken = status.Error(codes.Internal, "unable to generate token")
	errVerifyToken   = status.Error(codes.InvalidArgument, "token invalid")
	errUserVerify    = status.Error(codes.Unauthenticated, "unauthenticated")
)
