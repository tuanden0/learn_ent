package v1

type UserInfo struct {
	Id       uint64 `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Role     string `json:"role,omitempty"`
}
