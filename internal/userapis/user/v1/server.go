package v1

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	userv1 "github.com/tuanden0/learn_ent/proto/gen/go/v1/user"
	"google.golang.org/grpc"
)

func setupGrpcServerOptions() (opts []grpc.ServerOption) {
	return opts
}

func setupServeMuxOptions() (opts []runtime.ServeMuxOption) {
	return opts
}

func RunServer(srv Service, addr string) error {

	// Create new context
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Create server listener
	lis, lisErr := net.Listen("tcp", addr)
	if lisErr != nil {
		return lisErr
	}
	defer lis.Close()

	// Create error channel
	errChan := make(chan error)

	// Create grpcServer with options
	grpcServerOpts := setupGrpcServerOptions()
	grpcServer := grpc.NewServer(grpcServerOpts...)

	// Register gRPC service
	userv1.RegisterUserServiceServer(grpcServer, srv)
	userv1.RegisterUserAuthenServiceServer(grpcServer, srv)

	// Run gRPC server
	go func() {
		glog.Info("gRPC server is running")
		if err := grpcServer.Serve(lis); err != nil {
			errChan <- err
		}
	}()

	// Create HTTP server same port with gRPC
	gwServerMuxOpts := setupServeMuxOptions()
	gwmux := runtime.NewServeMux(gwServerMuxOpts...)

	// Register gRPC Gateway service
	if err := userv1.RegisterUserServiceHandlerServer(ctx, gwmux, srv); err != nil {
		return err
	}

	if err := userv1.RegisterUserAuthenServiceHandlerServer(ctx, gwmux, srv); err != nil {
		return err
	}

	gwServer := &http.Server{
		Handler: gwmux,
	}

	go func() {
		glog.Info("HTTP server is running")
		if err := gwServer.Serve(lis); err != nil {
			if !strings.Contains(err.Error(), "use of closed network connection") {
				errChan <- err
			}
		}
	}()

	// Handle shutdown signal
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		glog.Infof("got %v signal, graceful shutdown server", <-c)
		cancel()
		grpcServer.GracefulStop()
		gwServer.Shutdown(ctx)
		close(errChan)
	}()

	return <-errChan
}
