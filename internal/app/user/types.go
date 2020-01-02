package user

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	MALE Gender = iota
	FEMALE
	OTHER
)

type (
	Gender          int
	RegisterRequest struct {
		FirstName string `validate:"required" json:"first_name"`
		LastName  string `validate:"required" json:"last_name"`
		Gender    Gender `validate:"gte=0,lte=2" json:"gender"`
		Email     string `validate:"required,email" json:"email"`
		Password  string `validate:"required" json:"password"`
	}
	LoginRequest struct {
		Email    string `validate:"required,email" json:"email"`
		Password string `validate:"required" json:"password"`
	}
	Claims struct {
		Email string `validate:"required,email" json:"email"`
		jwt.StandardClaims
	}
	User struct {
		ID        string    `bson:"_id"`
		FirstName string    `bson:"first_name"`
		LastName  string    `bson:"last_name"`
		Gender    Gender    `bson:"gender"`
		Email     string    `bson:"email"`
		Password  string    `bson:"password"`
		Locked    bool      `bson:"locked"`
		CreatedAt time.Time `bson:"created_at"`
		UpdatedAt time.Time `bson:"updated_at"`
	}
)
