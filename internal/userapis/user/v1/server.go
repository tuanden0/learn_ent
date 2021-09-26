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
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

func setupGrpcServerOptions() (opts []grpc.ServerOption) {
	return opts
}

func setupServeMuxOptions() (opts []runtime.ServeMuxOption) {
	return opts
}

func setupClientDialOpts() []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithInsecure(),
	}
}

func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
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
	grpcServer := grpc.NewServer(setupGrpcServerOptions()...)

	// Register gRPC service
	userv1.RegisterUserServiceServer(grpcServer, srv)
	userv1.RegisterUserAuthenServiceServer(grpcServer, srv)

	// Create HTTP server same port with gRPC
	mux := runtime.NewServeMux(setupServeMuxOptions()...)

	// Register gRPC Gateway service
	if err := userv1.RegisterUserServiceHandlerFromEndpoint(ctx, mux, addr, setupClientDialOpts()); err != nil {
		return err
	}

	if err := userv1.RegisterUserAuthenServiceHandlerFromEndpoint(ctx, mux, addr, setupClientDialOpts()); err != nil {
		return err
	}

	// Create http server
	gwServer := &http.Server{
		Handler: grpcHandlerFunc(grpcServer, mux),
	}

	// Run http server
	go func() {
		glog.Info("HTTP server is running")
		if err := gwServer.Serve(lis); err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// Handle shutdown signal
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		glog.Infof("got %v signal, graceful shutdown server", <-c)
		cancel()
		gwServer.Shutdown(ctx)
		close(errChan)
	}()

	return <-errChan
}
