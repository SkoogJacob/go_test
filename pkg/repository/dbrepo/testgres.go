package dbrepo

import (
	"database/sql"
	"errors"
	"web_test/pkg/data"
)

type TestRepo struct {
	Users []*data.User
}

func (r *TestRepo) Connection() *sql.DB {
	panic("not applicable for this stub db") // TODO: Implement
}

func (r *TestRepo) AllUsers() ([]*data.User, error) {
	return r.Users, nil
}

func (r *TestRepo) GetUser(id int) (*data.User, error) {
	for _, u := range r.Users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (r *TestRepo) GetUserByEmail(email string) (*data.User, error) {
	for _, u := range r.Users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (r *TestRepo) UpdateUser(user data.User) error {
	for _, u := range r.Users {
		if u.ID == user.ID {
			*u = user
			return nil
		}
	}
	return errors.New("user not found")
}

func (r *TestRepo) DeleteUser(id int) error {
	panic("not implemented") // TODO: Implement
}

func (r *TestRepo) InsertUser(u data.User) (int, error) {
	panic("not implemented") // TODO: Implement
}

func (r *TestRepo) ResetPassword(id int, password string) error {
	panic("not implemented") // TODO: Implement
}

func (r *TestRepo) InsertUserImage(i data.UserImage) (int, error) {
	panic("not implemented") // TODO: Implement
}
