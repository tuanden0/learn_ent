package main

import (
	"context"
	"flag"
	"time"

	"github.com/golang/glog"
	userv1 "github.com/tuanden0/learn_ent/proto/gen/go/v1/user"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

var (
	addr = flag.String("ip", "0.0.0.0:8000", "server address:port")
)

func init() {
	flag.Lookup("v").Value.Set("2")
	flag.Lookup("logtostderr").Value.Set("true")
	flag.Parse()
}

func main() {

	// Make sure to write log to file
	defer glog.Flush()

	// Create client connection
	ctx, cancer := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancer()
	cc, err := grpc.DialContext(ctx, *addr, setupClientDialOpts()...)
	if err != nil {
		glog.Fatal(err)
	}
	defer cc.Close()
	c := userv1.NewUserServiceClient(cc)

	createUser(c)
	retrieveUser(c)
}

func retrieveUser(c userv1.UserServiceClient) {

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	// Retrieve user
	u1, err := c.Retrieve(ctx, &userv1.RetrieveRequest{
		Id: 1,
	})

	errDetails := handleServerError(err)
	if len(errDetails) != 0 {
		glog.Error(errDetails)
		return
	}

	glog.Infof("Retrieve user id 1: %v", u1)
}

func createUser(c userv1.UserServiceClient) {

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Create user
	u, err := c.Create(ctx, &userv1.CreateRequest{
		Username: "client",
		Password: "123123",
		Email:    "client.com",
		Role:     1,
	})

	errDetails := handleServerError(err)
	if len(errDetails) != 0 {
		glog.Error(errDetails)
		return
	}

	glog.Infof("Create user success: %v", u)
}

func setupClientDialOpts() []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithInsecure(),
	}
}

func handleServerError(err error) map[string]string {

	if err != nil {
		res := make(map[string]string)
		st, ok := status.FromError(err)
		if ok {
			if len(st.Details()) != 0 {
				for _, detail := range st.Details() {
					switch t := detail.(type) {
					case *errdetails.BadRequest:
						for _, violation := range t.GetFieldViolations() {
							res[violation.GetField()] = violation.GetDescription()
						}
					}
				}
			} else {
				res[st.Code().String()] = st.Message()
			}
		}

		return res
	}

	return nil
}
