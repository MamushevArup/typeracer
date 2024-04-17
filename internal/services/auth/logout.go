package auth

import (
	"context"
	"fmt"
)

func (a *auth) Logout(ctx context.Context, refresh string) error {

	err := parseRefreshToken(refresh)
	if err != nil {
		return fmt.Errorf("unable to parse refresh token %w", err)
	}

	err = a.repo.Auth.DeleteRefreshSession(ctx, refresh)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
