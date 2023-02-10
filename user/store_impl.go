package user

import (
	"context"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/taiwan-voting-guide/backend/model"
	"github.com/taiwan-voting-guide/backend/pg"
)

func New() Store {
	return &impl{}
}

type impl struct{}

func (im *impl) CreateByAuthResult(ctx context.Context, result *model.AuthResult) (*model.User, error) {
	conn, err := pg.Connect(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer conn.Close(ctx)

	u, err := result.ToUser()
	if err != nil {
		return nil, err
	}

	id, err := createNewId()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	u.Id = id

	if tag, err := conn.Exec(ctx, "INSERT INTO users (id, name, email, avatar_url, google_id) VALUES ($1, $2, $3, $4, $5)", u.Id, u.Name, u.Email, u.AvatarURL, u.GoogleId); err != nil {
		return nil, err
	} else if tag.RowsAffected() == 0 {
		return nil, ErrUserExist
	}

	return u, nil
}

func createNewId() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	return strings.ReplaceAll(id.String(), "-", ""), nil
}

func authResultGoogleToUser(ctx context.Context, result *model.AuthResultGoogle) *model.User {
	return &model.User{
		Email:     result.Payload.Claims["email"].(string),
		Name:      result.Payload.Claims["given_name"].(string),
		AvatarURL: result.Payload.Claims["picture"].(string),
		GoogleId:  result.Payload.Subject,
	}
}

func (im *impl) Get(ctx context.Context, id string) (*model.User, error) {
	conn, err := pg.Connect(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer conn.Close(ctx)

	row := conn.QueryRow(ctx, "SELECT name, email, avatar_url, google_id FROM users WHERE id = $1", id)

	var user model.User
	user.Id = id
	if err := row.Scan(&user.Name, &user.Email, &user.AvatarURL, &user.GoogleId); err == pgx.ErrNoRows {
		return nil, ErrUserNotFound
	} else if err != nil {
		log.Println(err)
		return nil, err
	}

	return &user, nil
}

func (im *impl) GetByAuthResult(ctx context.Context, result *model.AuthResult) (*model.User, error) {
	switch result.Type {
	case model.AuthTypeGoogle:
		return getByAuthResultGoogle(ctx, result.Google)
	default:
		return nil, ErrAuthTypeInvalid
	}
}

func getByAuthResultGoogle(ctx context.Context, result *model.AuthResultGoogle) (*model.User, error) {
	conn, err := pg.Connect(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer conn.Close(ctx)

	row := conn.QueryRow(ctx, "SELECT id, name, email, avatar_url FROM users WHERE google_id = $1", result.Payload.Subject)

	var user model.User
	if err := row.Scan(&user.Id, &user.Name, &user.Email, &user.AvatarURL); err == pgx.ErrNoRows {
		return nil, ErrUserNotFound
	} else if err != nil {
		log.Println(err)
		return nil, err
	}

	return &user, nil
}

func (im *impl) Update(ctx context.Context, user *model.User) error {
	conn, err := pg.Connect(ctx)
	if err != nil {
		log.Println(err)
		return err
	}
	defer conn.Close(ctx)

	if commandTag, err := conn.Exec(ctx, "UPDATE users SET name = $1, email = $2, avatar_url = $3, google_id = $4 WHERE id = $5 AND active = true", user.Name, user.Email, user.AvatarURL, user.GoogleId, user.Id); commandTag.RowsAffected() == 0 {
		return ErrUserNotFound
	} else if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (im *impl) Deactivate(ctx context.Context, id string) error {
	conn, err := pg.Connect(ctx)
	if err != nil {
		log.Println(err)
		return err
	}
	defer conn.Close(ctx)

	if commandTag, err := conn.Exec(ctx, "UPDATE users SET active = false WHERE id = $1", id); commandTag.RowsAffected() == 0 {
		return ErrUserNotFound
	} else if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (im *impl) Activate(ctx context.Context, id string) error {
	conn, err := pg.Connect(ctx)
	if err != nil {
		log.Println(err)
		return err
	}
	defer conn.Close(ctx)

	if commandTag, err := conn.Exec(ctx, "UPDATE users SET active = true WHERE id = $1", id); commandTag.RowsAffected() == 0 {
		return ErrUserNotFound
	} else if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
