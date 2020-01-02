package user

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type (
	repository interface {
		Insert(ctx context.Context, user User) error
		FindUserByEmail(ctx context.Context, email string) (*User, error)
	}
	Service struct {
		repo repository
	}
)

var jwtKey = []byte("my_secret_key")

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
	if err := validator.New().Struct(req); err != nil {
		logrus.Errorf("validation failed, error: %v", err)
		return "", err
	}
	// Create the JWT key used to create the signature
	dbuser, err := s.repo.FindUserByEmail(ctx, req.Email)
	if err != nil && err != ErrUserNotFound {
		logrus.Errorf("Failed to find user, error: %v", err)
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(dbuser.Password), []byte(req.Password)); err != nil {
		logrus.Errorf("Password wrong")
		return "", nil
	}

	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &Claims{
		Email: req.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		logrus.Errorf("Signing string fail")
		return "", err
	}
	return tokenString, nil
}

func (s *Service) ShowInfo(ctx context.Context, tokenString string) (*User, error) {
	// Parse the JWT string and store the result in `claims`.
	claims := &Claims{}
	user := &User{}
	tokenValidType := strings.Replace(tokenString, "Bearer ", "", 7)
	token, err := jwt.ParseWithClaims(tokenValidType, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		logrus.Errorf("Can not compare token, error: %v", err)
		return user, jwt.ErrSignatureInvalid
	}
	if !token.Valid {
		logrus.Errorf("the token is invalid, error: %v", err)

		return user, jwt.ErrInvalidKey
	}
	// Create the JWT key used to create the signature
	dbuser, err := s.repo.FindUserByEmail(ctx, claims.Email)
	if err != nil && err != ErrUserNotFound {
		logrus.Errorf("Failed to find user, error: %v", err)
		return user, err
	}
	return dbuser, nil
}
