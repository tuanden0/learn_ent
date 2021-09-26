package main

import (
	"flag"

	"github.com/golang/glog"
	v1 "github.com/tuanden0/learn_ent/internal/authapis/auth/v1"
)

var (
	addr = flag.String("ip", "0.0.0.0:8001", "server address:port")
)

func init() {
	flag.Lookup("v").Value.Set("2")
	flag.Lookup("logtostderr").Value.Set("true")
	flag.Parse()
}

func main() {

	defer glog.Flush()

	// Create auth repo
	authRepo := v1.NewRepoManager()

	// Create auth validator
	authValidator := v1.NewValidate()
	if err := authValidator.Init(); err != nil {
		glog.Fatal(err)
	}

	// Create auth service
	srv := v1.NewService(authRepo, authValidator)

	if err := v1.RunServer(srv, *addr); err != nil {
		glog.Fatal(err)
	}
}
