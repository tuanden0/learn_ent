package v1

import (
	"context"
	"os"
	"testing"

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
		Role:     0,
	}

	u, _ := r.Create(context.Background(), in)

	return u
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
