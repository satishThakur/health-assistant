package auth

import (
	"encoding/json"
	"log"
	"net/http"
)

// Handler handles authentication endpoints.
type Handler struct {
	googleVerifier *GoogleVerifier
	userRepo       *UserRepository
	tokenService   *TokenService
}

// NewHandler creates a new auth Handler.
func NewHandler(
	googleVerifier *GoogleVerifier,
	userRepo *UserRepository,
	tokenService *TokenService,
) *Handler {
	return &Handler{
		googleVerifier: googleVerifier,
		userRepo:       userRepo,
		tokenService:   tokenService,
	}
}

type googleAuthRequest struct {
	IDToken string `json:"id_token"`
}

type googleAuthResponse struct {
	Token       string `json:"token"`
	UserID      string `json:"user_id"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
}

// HandleGoogleAuth handles POST /api/v1/auth/google
func (h *Handler) HandleGoogleAuth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req googleAuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	if req.IDToken == "" {
		http.Error(w, `{"error":"id_token is required"}`, http.StatusBadRequest)
		return
	}

	claims, err := h.googleVerifier.VerifyIDToken(r.Context(), req.IDToken)
	if err != nil {
		log.Printf("Google token verification failed: %v", err)
		http.Error(w, `{"error":"invalid Google token"}`, http.StatusUnauthorized)
		return
	}

	user, err := h.userRepo.FindOrCreateUserByGoogleID(r.Context(), claims.Sub, claims.Email, claims.Name)
	if err != nil {
		log.Printf("Failed to find or create user: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	token, err := h.tokenService.GenerateToken(user.ID)
	if err != nil {
		log.Printf("Failed to generate token: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(googleAuthResponse{
		Token:       token,
		UserID:      user.ID,
		Email:       user.Email,
		DisplayName: user.DisplayName,
	})

	log.Printf("Google auth successful for user %s (%s)", user.ID, user.Email)
}
