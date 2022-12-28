package user

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"

	"github.com/taiwan-voting-guide/backend/pg"
)

var (
	errUserAlreadyExists = errors.New("user already exists")
	errUserNotFound      = errors.New("user not found")
)

type User struct {
	Id        string
	Name      string
	Email     string
	AvatarURL string

	GoogleId string
}

// Get a user by id. Return errUserNotFound if the user does not exist.
func Get(ctx context.Context, id string) (*User, error) {
	conn, err := pg.Connect(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer conn.Close(ctx)

	row := conn.QueryRow(ctx, "SELECT name, email, avatar_url, google_id FROM users WHERE id = $1", id)

	var user User
	if err := row.Scan(id, &user.Name, &user.Email, &user.AvatarURL, &user.GoogleId); err != pgx.ErrNoRows {
		return nil, errUserNotFound
	} else if err != nil {
		log.Println(err)
		return nil, err
	}

	return &user, nil
}

// Create a new user. Return errUserAlreadyExists if the user already exists.
func Create(ctx context.Context, user *User) error {
	conn, err := pg.Connect(ctx)
	if err != nil {
		log.Println(err)
		return err
	}
	defer conn.Close(ctx)

	if tag, err := conn.Exec(ctx, "INSERT INTO users (id, name, email, avatar_url, google_id) VALUES ($1, $2, $3, $4, $5)", user.Id, user.Name, user.Email, user.AvatarURL, user.GoogleId); tag.RowsAffected() == 0 {
		return errUserAlreadyExists
	} else if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// Update an existing user. Return errUserNotFound if the user does not exist.
func Update(ctx context.Context, user *User) error {
	conn, err := pg.Connect(ctx)
	if err != nil {
		log.Println(err)
		return err
	}
	defer conn.Close(ctx)

	if commandTag, err := conn.Exec(ctx, "UPDATE users SET name = $1, email = $2, avatar_url = $3, google_id = $4 WHERE id = $5 AND active = true", user.Name, user.Email, user.AvatarURL, user.GoogleId, user.Id); commandTag.RowsAffected() == 0 {
		return errUserNotFound
	} else if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// Deactivate a user by id. Return errUserNotFound if the user does not exist.
func Deactivate(ctx context.Context, id string) error {
	conn, err := pg.Connect(ctx)
	if err != nil {
		log.Println(err)
		return err
	}
	defer conn.Close(ctx)

	if commandTag, err := conn.Exec(ctx, "UPDATE users SET active = false WHERE id = $1", id); commandTag.RowsAffected() == 0 {
		return errUserNotFound
	} else if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// Reactivate a user by id. Return errUserNotFound if the user does not exist.
func Activate(ctx context.Context, id string) error {
	conn, err := pg.Connect(ctx)
	if err != nil {
		log.Println(err)
		return err
	}
	defer conn.Close(ctx)

	if commandTag, err := conn.Exec(ctx, "UPDATE users SET active = true WHERE id = $1", id); commandTag.RowsAffected() == 0 {
		return errUserNotFound
	} else if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
