# LEARN ENT

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

## Bugs

## Handle errors

[use of closed network connection](https://github.com/grpc-ecosystem/grpc-gateway/issues/727)
