package v1

import "time"

type Key string

const (
	userServiceUrl     = "0.0.0.0:8000"
	secret             = "vJ1ZnpDoSYMSzGQboSswwuCnnv4T36jl"
	tokenDuration      = 15 * time.Minute
	JWTKey         Key = "user"
)
