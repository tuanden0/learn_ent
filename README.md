# LEARN ENT

## How to run

```bash
# Log to stderr, file and using version 2
go run cmd/userapis/user/v1/main.go -v=2 -alsologtostderr=1 -log_dir=log

# Log only stderr and using version 2
go run cmd/userapis/user/v1/main.go -logtostderr=true -v=2
```

## Handle errors

[use of closed network connection](https://github.com/grpc-ecosystem/grpc-gateway/issues/727)
