package v1

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	errLogin = status.Errorf(codes.Unauthenticated, "incorrect username/password")
)
