package auth

import (
	"context"
	"fmt"
)

func (a *auth) AdminRefresh(ctx context.Context, refresh string) (string, string, error) {
	id, err := a.repo.Auth.AdminByRefresh(ctx, refresh)
	if err != nil {
		return "", "", fmt.Errorf("%w", err)
	}

	err = parseRefreshToken(refresh)
	if err != nil {
		return "", "", fmt.Errorf("unable to parse refresh token %w", err)
	}

	refreshNew, err := a.generateRefreshToken()
	if err != nil {
		return "", "", fmt.Errorf("%w", err)
	}

	access, err := a.generateAccessToken(id, "admin")
	if err != nil {
		return "", "", fmt.Errorf("%w", err)
	}

	err = a.repo.Auth.UpdateRefresh(ctx, id, refreshNew)
	if err != nil {
		return "", "", fmt.Errorf("%w", err)
	}

	return access, refreshNew, nil
}
