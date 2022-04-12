package repository

import (
	"api/src/model"
	"database/sql"
	"fmt"
)

// Users represents a connection to the users table.
type Users struct {
	db *sql.DB
}

// NewUsersRepository generates a connection to the users table.
// returns a connection to the users table.
func NewUsersRepository(db *sql.DB) *Users {
	return &Users{db}
}

// CreateUser create an user.
// returns the ID of the created user or an error.
func (repository Users) CreateUser(user model.User) (uint64, error) {

	statement, error := repository.db.Prepare(
		"insert into users(name, nick, email, password) values(?, ?, ?, ?)",
	)
	if error != nil {
		return 0, error
	}

	defer statement.Close()

	result, error := statement.Exec(user.Name, user.Nick, user.Email, user.Password)
	if error != nil {
		return 0, error
	}

	lastInsertId, error := result.LastInsertId()
	if error != nil {
		return 0, error
	}

	return uint64(lastInsertId), nil
}

// ListUsersByFilter list users by a filter.
// returns the list of users found or an error.
func (repository Users) ListUsersByFilter(nameOrNick string) ([]model.User, error) {
	nameOrNick = fmt.Sprintf("%%%s%%", nameOrNick)

	results, error := repository.db.Query(
		"select id, name, nick, email, createdAt from users where name like ? or nick like ?",
		nameOrNick,
		nameOrNick,
	)
	if error != nil {
		return nil, error
	}

	defer results.Close()

	var users []model.User

	for results.Next() {
		var user model.User

		if error = results.Scan(
			&user.ID,
			&user.Name,
			&user.Nick,
			&user.Email,
			&user.CreatedAt,
		); error != nil {
			return nil, error
		}

		users = append(users, user)
	}

	return users, nil
}

// FindUserById search for a user by the given ID.
// return a user or an error.
func (repository Users) FindUserById(userId uint64) (model.User, error) {

	results, error := repository.db.Query(
		"select id, name, nick, email, createdAt from users where id = ?",
		fmt.Sprintf("%d", userId),
	)
	if error != nil {
		return model.User{}, error
	}

	defer results.Close()

	var user model.User

	if results.Next() {
		if error = results.Scan(
			&user.ID,
			&user.Name,
			&user.Nick,
			&user.Email,
			&user.CreatedAt,
		); error != nil {
			return model.User{}, error
		}
	}

	return user, nil
}

// UpdateUserById updates a user by the given ID.
// returns an error if unable to update the user.
func (repository Users) UpdateUserById(userId uint64, user model.User) error {

	statement, error := repository.db.Prepare(
		"update users set name = ?, nick = ?, email = ? where id = ?",
	)
	if error != nil {
		return error
	}

	defer statement.Close()

	if _, error := statement.Exec(user.Name, user.Nick, user.Email, userId); error != nil {
		return error
	}

	return nil
}

// DeleteUserById delete a user by the given ID.
// returns an error if unable to delete the user.
func (repository Users) DeleteUserById(userId uint64) error {

	statement, error := repository.db.Prepare("delete from users where id = ?")
	if error != nil {
		return error
	}

	defer statement.Close()

	if _, error := statement.Exec(userId); error != nil {
		return error
	}

	return nil
}

// FindUserByEmail search for a user by the given email.
// returns a user or an error.
func (repository Users) FindUserByEmail(email string) (model.User, error) {

	results, error := repository.db.Query(
		"select id, password from users where email = ?",
		email,
	)
	if error != nil {
		return model.User{}, error
	}

	defer results.Close()

	var user model.User

	if results.Next() {
		if error = results.Scan(
			&user.ID,
			&user.Password,
		); error != nil {
			return model.User{}, error
		}
	}

	return user, nil
}

// FollowUserById follows a user by the given ID.
// returns an error if unable to follow the user.
func (repository Users) FollowUserById(userId, followerId uint64) error {

	statement, error := repository.db.Prepare(
		"insert ignore into followers(user_id, follower_id) values(?, ?)",
	)
	if error != nil {
		return error
	}

	defer statement.Close()

	if _, error := statement.Exec(userId, followerId); error != nil {
		return error
	}

	return nil
}

// UnfollowUserById unfollow a user by the given ID.
// returns an error if unable to unfollow the user.
func (repository Users) UnfollowUserById(userId, followerId uint64) error {

	statement, error := repository.db.Prepare("delete from followers where user_id = ? and follower_id = ?")
	if error != nil {
		return error
	}

	defer statement.Close()

	if _, error := statement.Exec(userId, followerId); error != nil {
		return error
	}

	return nil
}

// ListFollowersByFollowedUserId lists a user's followers.
// returns a list of followers or an error.
func (repository Users) ListFollowersByFollowedUserId(followedUserId uint64) ([]model.User, error) {

	results, error := repository.db.Query(
		`
			select 
				u.id, u.name, u.nick, u.email, u.createdAt 
			from 
				users u 
				join followers f on f.follower_id = u.id 
			where f.user_id = ?	
		`,
		followedUserId,
	)
	if error != nil {
		return nil, error
	}

	defer results.Close()

	var followers []model.User

	for results.Next() {
		var user model.User

		if error = results.Scan(
			&user.ID,
			&user.Name,
			&user.Nick,
			&user.Email,
			&user.CreatedAt,
		); error != nil {
			return nil, error
		}

		followers = append(followers, user)
	}

	return followers, nil
}

// ListFollowedByFollowerId lists people the user follows.
// returns a list of people the user follows or an error.
func (repository Users) ListFollowedByFollowerId(followedUserId uint64) ([]model.User, error) {

	results, error := repository.db.Query(
		`
			select 
				u.id, u.name, u.nick, u.email, u.createdAt 
			from 
				users u 
				join followers f on f.user_id = u.id 
			where f.follower_id = ?	
		`,
		followedUserId,
	)
	if error != nil {
		return nil, error
	}

	defer results.Close()

	var followers []model.User

	for results.Next() {
		var user model.User

		if error = results.Scan(
			&user.ID,
			&user.Name,
			&user.Nick,
			&user.Email,
			&user.CreatedAt,
		); error != nil {
			return nil, error
		}

		followers = append(followers, user)
	}

	return followers, nil
}

// FindPasswordByUserId search for the user's password by the given ID.
// returns the password or an error.
func (repository Users) FindPasswordByUserId(userId uint64) (string, error) {

	results, error := repository.db.Query(
		"select password from users where id = ?",
		fmt.Sprintf("%d", userId),
	)
	if error != nil {
		return "", error
	}

	defer results.Close()

	var user model.User

	if results.Next() {
		if error = results.Scan(&user.Password); error != nil {
			return "", error
		}
	}

	return user.Password, nil
}

// UpdatePasswordByUserId updates the user's password by the given ID.
// returns an error if unable to update the password.
func (repository Users) UpdatePasswordByUserId(userId uint64, newPasswordHash string) error {

	statement, error := repository.db.Prepare("update users set password = ? where id = ?")
	if error != nil {
		return error
	}

	defer statement.Close()

	if _, error := statement.Exec(newPasswordHash, userId); error != nil {
		return error
	}

	return nil
}
