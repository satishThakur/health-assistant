package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// GoogleClaims holds claims returned from the Google tokeninfo endpoint.
type GoogleClaims struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	Audience      string `json:"aud"`
	EmailVerified string `json:"email_verified"`
}

// GoogleVerifier verifies Google ID tokens via the tokeninfo endpoint.
type GoogleVerifier struct {
	clientID   string
	httpClient *http.Client
}

// NewGoogleVerifier creates a GoogleVerifier.
// If clientID is empty, audience validation is skipped (dev mode).
func NewGoogleVerifier(clientID string) *GoogleVerifier {
	return &GoogleVerifier{
		clientID:   clientID,
		httpClient: &http.Client{},
	}
}

// VerifyIDToken calls the Google tokeninfo endpoint and validates the token.
func (v *GoogleVerifier) VerifyIDToken(ctx context.Context, idToken string) (*GoogleClaims, error) {
	tokenInfoURL := "https://oauth2.googleapis.com/tokeninfo?id_token=" + url.QueryEscape(idToken)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, tokenInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating tokeninfo request: %w", err)
	}

	resp, err := v.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("calling tokeninfo endpoint: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tokeninfo returned status %d", resp.StatusCode)
	}

	var claims GoogleClaims
	if err := json.NewDecoder(resp.Body).Decode(&claims); err != nil {
		return nil, fmt.Errorf("decoding tokeninfo response: %w", err)
	}

	// Validate audience against our client ID (skip in dev mode when clientID is empty)
	if v.clientID != "" && claims.Audience != v.clientID {
		return nil, fmt.Errorf("token audience %q does not match expected client ID", claims.Audience)
	}

	if claims.EmailVerified != "true" {
		return nil, fmt.Errorf("google account email is not verified")
	}

	return &claims, nil
}
