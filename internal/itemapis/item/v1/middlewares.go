package v1

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4/middleware"
)

type userClaims struct {
	jwt.RegisteredClaims
	UserId   uint64 `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

var configJWT = middleware.JWTConfig{
	Claims:     &userClaims{},
	SigningKey: secretAuthenToken,
}
