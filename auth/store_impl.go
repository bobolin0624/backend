package auth

import (
	"context"
	"log"
	"net/http"
	"os"

	"google.golang.org/api/idtoken"

	"github.com/taiwan-voting-guide/backend/model"
)

func New() Store {
	return &impl{
		httpClient: &http.Client{},
	}
}

type impl struct {
	httpClient *http.Client
}

func (im *impl) Auth(ctx context.Context, info *model.AuthInfo) (*model.AuthResult, error) {
	switch info.Type {
	case model.AuthTypeGoogle:
		return im.authWithGoogle(ctx, info.Google)
	default:
		log.Printf("unknown auth type: %d", info.Type)
		return nil, ErrTypeInvalid
	}
}

func (im *impl) authWithGoogle(ctx context.Context, info *model.AuthInfoGoogle) (*model.AuthResult, error) {
	payload, err := idtoken.Validate(ctx, info.IdToken, os.Getenv("GOOGLE_CLIENT_ID"))
	if err != nil {
		return nil, ErrTokenAudienceInvalid
	}

	return &model.AuthResult{
		Type: model.AuthTypeGoogle,
		Google: &model.AuthResultGoogle{
			Payload: payload,
		},
	}, nil
}
