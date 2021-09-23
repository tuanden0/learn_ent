package v1

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/gogo/protobuf/types"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/tuanden0/learn_ent/internal/userapis/user/v1/ent"
	"github.com/tuanden0/learn_ent/internal/userapis/user/v1/ent/enttest"
	"github.com/tuanden0/learn_ent/internal/userapis/user/v1/ent/migrate"
	userv1 "github.com/tuanden0/learn_ent/proto/gen/go/v1/user"
)

var r Repository

func genUser() *ent.User {

	in := &userv1.CreateRequest{
		Username: "test_create",
		Password: "test123456",
		Email:    "test6478@localhost.com",
		Role:     1,
	}

	u, _ := r.Create(context.Background(), in)

	return u
}

func genUsers() []*ent.User {

	us := []*ent.User{}

	for i := 0; i < 10; i++ {
		in := &userv1.CreateRequest{
			Username: fmt.Sprintf("test_create_%v", i),
			Password: "test123456",
			Email:    fmt.Sprintf("test_create_%v@localhost.com", i),
			Role:     1,
		}
		u, _ := r.Create(context.Background(), in)
		us = append(us, u)
	}

	return us
}

func TestMain(m *testing.M) {

	// Mock database
	var t enttest.TestingT
	opts := []enttest.Option{
		enttest.WithMigrateOptions(migrate.WithGlobalUniqueID(true)),
	}
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1", opts...)
	defer client.Close()

	r = NewRepoManager(client)

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestCreateOK(t *testing.T) {

	// Initialize
	in := &userv1.CreateRequest{
		Username: "test",
		Password: "test123",
		Email:    "test@localhost.com",
		Role:     1,
	}

	// Execute
	u, err := r.Create(context.Background(), in)

	// Assertion
	assert.Nil(t, err)
	assert.NotNil(t, u)
	assert.EqualValues(t, in.GetUsername(), u.Username)
	assert.EqualValues(t, in.GetEmail(), u.Email)
	assert.EqualValues(t, in.GetRole(), u.Role)
	assert.NotEqualValues(t, in.GetPassword(), u.Password)
	assert.NotEqualValues(t, 0, u.ID)
}

func TestRetrieveOK(t *testing.T) {

	// Init
	dumpy := genUser()

	// Execute
	u, err := r.Retrieve(context.Background(), uint64(dumpy.ID))

	// Assertion
	assert.Nil(t, err)
	assert.NotNil(t, u)
	assert.EqualValues(t, dumpy.Username, u.Username)
	assert.EqualValues(t, dumpy.Email, u.Email)
	assert.EqualValues(t, dumpy.Role, u.Role)
	assert.EqualValues(t, dumpy.Password, u.Password)
	assert.EqualValues(t, dumpy.ID, u.ID)
}

func TestUpdateOK(t *testing.T) {

	// Init
	dumpy := genUser()
	in := &userv1.UpdateRequest{
		Id: uint64(dumpy.ID),
		Password: &types.StringValue{
			Value: "123456",
		},
		Email: &types.StringValue{
			Value: "updated@local.com",
		},
	}

	// Execute
	u, err := r.Update(context.Background(), in)

	// Assertion
	assert.Nil(t, err)
	assert.NotNil(t, u)
	assert.EqualValues(t, dumpy.ID, u.ID)
	assert.EqualValues(t, dumpy.Username, u.Username)
	assert.NotEqualValues(t, dumpy.Password, u.Password)
	assert.NotEqualValues(t, dumpy.Email, u.Email)
	assert.EqualValues(t, dumpy.Role, u.Role)
}

func TestDeleteOK(t *testing.T) {

	// Init
	dumpy := genUser()

	// Execute
	err := r.Delete(context.Background(), uint64(dumpy.ID))

	// Assertion
	assert.Nil(t, err)
}

func TestListWithEmptyBaseOK(t *testing.T) {

	// Init
	_ = genUsers()
	in := &userv1.ListRequest{}

	// Execute
	us, err := r.List(context.Background(), in)

	// Assertion
	assert.Nil(t, err)
	assert.NotNil(t, us)
	assert.EqualValues(t, 5, len(us))
}

func TestListWithBaseOK(t *testing.T) {

	// Init
	genUs := genUsers()
	in := &userv1.ListRequest{
		Pagination: &userv1.Pagination{
			Limit: 5,
			Page:  1,
		},
		Filters: []*userv1.Filter{
			{
				Key:    "id",
				Value:  "2",
				Method: "=",
			},
		},
		Sort: &userv1.Sort{
			Key:   "id",
			IsAsc: true,
		},
	}

	// Execute
	us, err := r.List(context.Background(), in)

	// Assertion
	assert.Nil(t, err)
	assert.NotNil(t, us)
	assert.EqualValues(t, 1, len(us))

	for _, u := range genUs {
		if u.ID == us[0].ID {
			assert.EqualValues(t, u.Email, us[0].Email)
			assert.EqualValues(t, u.Username, us[0].Username)
			assert.EqualValues(t, u.Password, us[0].Password)
		}
	}
}
