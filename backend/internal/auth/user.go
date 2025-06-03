package auth

import (
	"context"
	"errors"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/user"
)

// Devuelve el ID del usuario autenticado desde el contexto de Clerk.
func GetAuthenticatedUserID(ctx context.Context) (string, error) {
	claims, ok := clerk.SessionClaimsFromContext(ctx)
	if !ok {
		return "", errors.New("unauthorized: no session claims")
	}

	user, err := user.Get(ctx, claims.Subject)
	if err != nil {
		return "", err
	}

	return user.ID, nil
}
