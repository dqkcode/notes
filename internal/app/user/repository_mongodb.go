package user

import (
	"context"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type (
	MongoDBRepository struct {
		session *mgo.Session
	}
)

func NewMongoDBRepository(session *mgo.Session) *MongoDBRepository {
	return &MongoDBRepository{
		session: session,
	}
}
func (r *MongoDBRepository) Insert(ctx context.Context, user User) error {
	s := r.session.Clone()
	defer s.Close()
	if err := s.DB("").C("users").Insert(user); err != nil {
		return err
	}
	return nil
}
func (r *MongoDBRepository) FindUserByEmail(ctx context.Context, email string) (*User, error) {
	s := r.session.Clone()
	defer s.Close()
	var user User
	if err := s.DB("").C("users").Find(bson.M{
		"email": email,
	}).One(&user); err != nil {
		if err == mgo.ErrNotFound {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}
