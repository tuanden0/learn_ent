package v1

import (
	"context"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tuanden0/learn_ent/internal/userapis/user/v1/ent"
)

func GetDatabase() (*ent.Client, error) {

	client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		return nil, fmt.Errorf("failed opening connection to sqlite: %v", err)
	}

	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		return nil, fmt.Errorf("failed creating schema resources: %v", err)
	}

	return client, nil
}
