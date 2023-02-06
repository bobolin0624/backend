package auth

import (
	"context"
	"log"
	"net/http"
	"os"

	"google.golang.org/api/idtoken"
)

func New() Store {
	return &impl{
		httpClient: &http.Client{},
	}
}

type impl struct {
	httpClient *http.Client
}

func (im *impl) Auth(ctx context.Context, info *Info) (*Result, error) {
	switch info.Type {
	case TypeGoogle:
		return im.authWithGoogle(ctx, info.Google)
	default:
		log.Printf("unknown auth type: %d", info.Type)
		return nil, ErrTypeInvalid
	}
}

func (im *impl) authWithGoogle(ctx context.Context, info *InfoGoogle) (*Result, error) {
	payload, err := idtoken.Validate(ctx, info.IdToken, os.Getenv("GOOGLE_CLIENT_ID"))
	if err != nil {
		return nil, ErrTokenAudienceInvalid
	}

	return &Result{
		Type: TypeGoogle,
		Google: &ResultGoogle{
			Payload: payload,
		},
	}, nil
}
