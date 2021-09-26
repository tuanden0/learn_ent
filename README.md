# LEARN ENT

This is repository is not only learning about Ent ORM but also gRPC, gRPC gateway and more things. It will be implement by owner idea to know more a bout go world.

## How to run

```bash
# Default was auto set -logtostderr=true -v=2
go run cmd/userapis/user/v1/server/server.go
go run cmd/authapis/auth/v1/server/server.go

# Log to stderr, file and using version 2
go run cmd/userapis/user/v1/server/server.go -v=2 -alsologtostderr=1 -log_dir=log

# Log only stderr and using version 2
go run cmd/userapis/user/v1/server/server.go -logtostderr=true -v=2
```

## Todos

Implement validate user credential on User service

Re-Implement validate conditions on Auth service to use token to call verify method

## Bugs
P1: Implement authentication between service User and Auth but it take to much time (1 - 4 second) to verify each other.

## Resolve
P1: Using http handler (grpcHandlerFunc) to handle connection instead of running gRPC server.
- We make everything more easier if we use TLS instead of WithInsecure.
- If use TLS, we don't need to implement HTTP handler (grpcHandlerFunc) to check if request is HTTP or gRPC and the speed is same as use grpcHandlerFunc.

## Handle errors

[use of closed network connection](https://github.com/grpc-ecosystem/grpc-gateway/issues/727)

[Bug P1](https://github.com/philips/grpc-gateway-example/issues/22#issuecomment-490733965)