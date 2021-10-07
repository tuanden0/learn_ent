package main

import (
	"context"
	"flag"

	"github.com/labstack/echo/v4"
	v1 "github.com/tuanden0/learn_ent/internal/itemapis/item/v1"
)

var (
	addr = flag.String("ip", "0.0.0.0:8002", "server address:port")
)

func init() {
	flag.Parse()
}

func main() {

	// Init echo server
	e := echo.New()

	// Init database
	client, err := v1.GetClientConnection()
	if err != nil {
		e.Logger.Fatalf("failed to connect mongodb: %v", err)
	}

	// Get collection
	coll := client.Database("item").Collection("itemData")

	// Handle close db
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			e.Logger.Fatalf("failed to disconnect to mongodb %v", err)
		}
	}()

	// Create item repository
	repoManager := v1.NewRepoManager(coll)

	// Create item service
	itemService := v1.NewService(e.Logger, repoManager)

	// Run server
	if err := v1.RunServer(e, itemService, *addr); err != nil {
		e.Logger.Error(err)
	}
}
