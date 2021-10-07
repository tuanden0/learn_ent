package v1

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type Service interface {
	Ping(c echo.Context) error
	Create(c echo.Context) error
	Retrieve(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
}

type service struct {
	log  echo.Logger
	repo Repository
}

func NewService(log echo.Logger, repo Repository) Service {
	return &service{log: log, repo: repo}
}

func (s *service) Ping(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

func (s *service) Create(c echo.Context) error {

	// Get user
	user := c.Get("user").(*jwt.Token)
	claims, ok := user.Claims.(*userClaims)
	if !ok {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "access denied"})
	}

	// Bind data to CreateItemRequest struct
	in := &CreateItemRequest{}
	if err := c.Bind(in); err != nil {
		s.log.Errorf("failed to bind create item input %v", err)
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
	}
	in.Owner = claims.Username

	// Update data to mongo db
	out, err := s.repo.Create(context.Background(), in)
	if err != nil {
		s.log.Errorf("failed to create item %v", err)
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, out)
}

func (s *service) Retrieve(c echo.Context) error {

	// Get user
	user := c.Get("user").(*jwt.Token)
	_, ok := user.Claims.(*userClaims)
	if !ok {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "access denied"})
	}

	id := c.Param("id")
	out, err := s.repo.Retrieve(context.Background(), id)
	if err != nil {
		s.log.Errorf("failed to retrieve item %v", err)
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, out)
}

func (s *service) Update(c echo.Context) error {

	// Get user
	user := c.Get("user").(*jwt.Token)
	claims, ok := user.Claims.(*userClaims)
	if !ok {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "access denied"})
	}

	id := c.Param("id")
	in := &UpdateItemRequest{}
	if err := c.Bind(in); err != nil {
		s.log.Errorf("failed to bind update item input %v", err)
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
	}
	in.Owner = claims.Username

	out, err := s.repo.Update(context.Background(), id, in)
	if err != nil {
		s.log.Errorf("failed to update item %v", err)
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, out)
}

func (s *service) Delete(c echo.Context) error {

	// Get user
	user := c.Get("user").(*jwt.Token)
	_, ok := user.Claims.(*userClaims)
	if !ok {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "access denied"})
	}

	id := c.Param("id")
	if err := s.repo.Delete(context.Background(), id); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
