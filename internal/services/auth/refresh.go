package auth

import (
	"context"
	"fmt"
	"time"
)

func (a *auth) RefreshToken(ctx context.Context, refresh, fingerprint string) (string, string, error) {

	// user_id, role, fingerprint
	racer, err := a.repo.Auth.Fingerprint(ctx, fingerprint, refresh)
	if err != nil {
		return "", "", fmt.Errorf("%w", err)
	}

	// here is the repo layer to remove from session
	err = a.repo.Auth.DeleteSession(ctx, fingerprint, refresh)
	if err != nil {
		return "", "", fmt.Errorf("%w", err)
	}

	// check the token expires or not
	if err = a.parseRefreshToken(refresh); err != nil {
		return "", "", fmt.Errorf("unable to parse refresh token %w", err)
	}

	// check the fingerprint are the same
	if racer.Fingerprint != fingerprint {
		return "", "", fmt.Errorf("fingerprint does not match")
	}

	racer.LastLogin = time.Now()

	racer.RefreshToken, err = a.generateRefreshToken()
	if err != nil {
		return "", "", fmt.Errorf("%w", err)
	}

	// or new refresh and store in the session
	//update racer table refresh and fingerprint
	if err = a.repo.Auth.InsertSession(ctx, racer); err != nil {
		return "", "", fmt.Errorf("%w", err)
	}

	// create endpoint token and create refresh token and send back to the frontend
	access, err := a.generateAccessToken(racer.ID.String(), racer.Role)
	if err != nil {
		return "", "", fmt.Errorf("%w", err)
	}

	return access, racer.RefreshToken, nil
}
