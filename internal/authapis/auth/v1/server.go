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
	authv1 "github.com/tuanden0/learn_ent/proto/gen/go/v1/auth"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

func setupGrpcServerOptions() (opts []grpc.ServerOption) {
	opts = append(opts, grpc.ChainUnaryInterceptor(
		authInterceptor,
		logUnaryInterceptor,
	))
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	lis, lisErr := net.Listen("tcp", addr)
	if lisErr != nil {
		return lisErr
	}
	defer lis.Close()

	errChan := make(chan error)

	grpcServer := grpc.NewServer(setupGrpcServerOptions()...)
	authv1.RegisterAuthenServiceServer(grpcServer, srv)

	mux := runtime.NewServeMux(setupServeMuxOptions()...)
	if err := authv1.RegisterAuthenServiceHandlerFromEndpoint(ctx, mux, addr, setupClientDialOpts()); err != nil {
		return err
	}

	gwServer := &http.Server{
		Handler: grpcHandlerFunc(grpcServer, mux),
	}

	go func() {
		glog.Info("HTTP server is running")
		if err := gwServer.Serve(lis); err != nil {
			if !strings.Contains(err.Error(), "use of closed network connection") {
				errChan <- err
			}
		}
	}()

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
