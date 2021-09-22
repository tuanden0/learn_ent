package main

import (
	"flag"

	"github.com/golang/glog"
	v1 "github.com/tuanden0/learn_ent/internal/userapis/user/v1"
)

var (
	addr = flag.String("ip", "0.0.0.0:8000", "server address:port")
)

func init() {
	flag.Parse()
}

func main() {

	// Make sure to write log to file
	defer glog.Flush()

	// Create user db
	userDB, err := v1.GetDatabase()
	if err != nil {
		glog.Fatal(err)
	}
	defer userDB.Close()

	// Create user repo
	userRepo := v1.NewRepoManager(userDB)

	// Create new service
	srv := v1.NewService(userRepo)

	if err := v1.RunServer(srv, *addr); err != nil {
		glog.Fatal(err)
	}
}
