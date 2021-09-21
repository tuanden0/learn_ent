package v1

import (
	"context"
	"fmt"
	"strconv"

	"github.com/tuanden0/learn_ent/internal/userapis/user/v1/ent"
	"github.com/tuanden0/learn_ent/internal/userapis/user/v1/ent/user"
	userv1 "github.com/tuanden0/learn_ent/proto/gen/go/v1/user"
)

type Repository interface {
	Create(ctx context.Context, u *userv1.CreateRequest) (*ent.User, error)
	Retrieve(ctx context.Context, id uint64) (*ent.User, error)
	Update(ctx context.Context, u *userv1.UpdateRequest) (*ent.User, error)
	Delete(ctx context.Context, id uint64) error
	List(ctx context.Context, in *userv1.ListRequest) ([]*ent.User, error)
}

type repoManager struct {
	client *ent.Client
}

func NewRepoManager(client *ent.Client) Repository {
	return &repoManager{
		client: client,
	}
}

func (r *repoManager) Create(ctx context.Context, u *userv1.CreateRequest) (*ent.User, error) {

	user, err := r.client.User.
		Create().
		SetUsername(u.GetUsername()).
		SetEmail(u.GetEmail()).
		SetPassword(u.GetPassword()).
		SetRole(uint32(u.GetRole())).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *repoManager) Retrieve(ctx context.Context, id uint64) (*ent.User, error) {

	u, err := r.client.User.Get(ctx, int(id))
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *repoManager) Update(ctx context.Context, u *userv1.UpdateRequest) (*ent.User, error) {

	user := r.client.User.UpdateOneID(int(u.GetId()))

	if u.GetEmail() != nil {
		user = user.SetEmail(u.GetEmail().GetValue())
	}

	if u.GetPassword() != nil {
		user = user.SetPassword(u.GetPassword().GetValue())
	}

	ux, err := user.Save(ctx)
	if err != nil {
		return nil, err
	}

	return ux, nil
}

func (r *repoManager) Delete(ctx context.Context, id uint64) error {
	return r.client.User.DeleteOneID(int(id)).Exec(ctx)
}

func (r *repoManager) List(ctx context.Context, in *userv1.ListRequest) ([]*ent.User, error) {

	users := r.client.User.Query()

	// Only support = filter
	fs := in.GetFilters()
	for _, f := range fs {
		if f.GetMethod() == "=" {
			if f.GetKey() == "username" {
				users = users.Where(user.UsernameEQ(f.GetValue()))
			} else if f.GetKey() == "email" {
				users = users.Where(user.EmailContains(f.GetKey()))
			} else if f.GetKey() == "id" {
				id, err := strconv.ParseUint(f.GetValue(), 10, 64)
				if err != nil {
					return nil, fmt.Errorf("invalid id")
				}
				users = users.Where(user.IDEQ(int(id)))
			}
		}
	}

	pg := in.GetPagination()
	limit := int(pg.GetLimit())
	if limit > 100 {
		limit = 100
	} else if limit <= 0 {
		limit = 5
	}

	page := int(pg.GetPage())
	if page <= 0 {
		page = 1
	}

	users = users.Limit(limit).Offset(limit * (page - 1))

	sort := in.GetSort()
	isAsc := sort.GetIsAsc()
	key := sort.GetKey()
	if key == "" {
		key = "id"
	}
	if isAsc {
		users = users.Order(ent.Asc(key))
	} else {
		users = users.Order(ent.Desc(key))
	}

	us, err := users.All(ctx)
	if err != nil {
		return nil, err
	}

	return us, nil
}
