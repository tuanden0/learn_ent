package v1

import (
	"github.com/tuanden0/learn_ent/internal/userapis/user/v1/ent"
	userv1 "github.com/tuanden0/learn_ent/proto/gen/go/v1/user"
)

func mapCreateResponse(u *ent.User) *userv1.CreateResponse {
	return &userv1.CreateResponse{
		Id:       uint64(u.ID),
		Username: u.Username,
		Email:    u.Email,
		Role:     userv1.Role(u.Role),
	}
}

func mapRetrieveResponse(u *ent.User) *userv1.RetrieveResponse {
	return &userv1.RetrieveResponse{
		Id:       uint64(u.ID),
		Username: u.Username,
		Email:    u.Email,
		Role:     userv1.Role(u.Role),
	}
}

func mapUpdateResponse(u *ent.User) *userv1.UpdateResponse {
	return &userv1.UpdateResponse{
		Id:       uint64(u.ID),
		Username: u.Username,
		Email:    u.Email,
		Role:     userv1.Role(u.Role),
	}
}

func mapListResponse(us []*ent.User) *userv1.ListResponse {

	l := make([]*userv1.UserList, 0)
	for _, u := range us {
		l = append(l, &userv1.UserList{
			Id:       uint64(u.ID),
			Username: u.Username,
			Email:    u.Email,
			Role:     userv1.Role(u.Role),
		})
	}

	return &userv1.ListResponse{Users: l}
}
