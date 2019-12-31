package user

import (
	"context"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type (
	repository interface {
		Insert(ctx context.Context, user User) error
		FindUserByEmail(ctx context.Context, email string) (*User, error)
		GetPasswordByEmail(ctx context.Context, email string) (string, error)
	}
	Service struct {
		repo repository
	}
)

func NewService(repo repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Register(ctx context.Context, req RegisterRequest) (string, error) {
	if err := validator.New().Struct(req); err != nil {
		logrus.Errorf("validation failed, error: %v", err)
		return "", err
	}
	dbuser, err := s.repo.FindUserByEmail(ctx, req.Email)
	if err != nil && err != ErrUserNotFound {
		logrus.Errorf("Failed to find user, error: %v", err)
		return "", err
	}
	if dbuser != nil {
		return "", ErrUserAlreadyExist
	}
	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logrus.Errorf("failed to gen password, err: %v", err)
		return "", fmt.Errorf("failed to register")
	}
	user := User{
		ID:        uuid.New().String(),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Gender:    req.Gender,
		Password:  string(password),
		CreatedAt: time.Now(),
	}
	if err := s.repo.Insert(ctx, user); err != nil {
		logrus.Errorf("failed to insert user, err: %v", err)
		return "", fmt.Errorf("failed to register: %w", err)
	}
	return user.ID, nil
}
func (s *Service) Login(ctx context.Context, req LoginRequest) (string, error) {

}
