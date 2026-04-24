package service

import (
	"encoding/json"
	"net/http"
	"grc_be/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type AuthService struct {
	auc *biz.AuthUseCase
	log *log.Helper
}

func NewAuthService(auc *biz.AuthUseCase, logger log.Logger) *AuthService {
	return &AuthService{
		auc: auc,
		log: log.NewHelper(logger),
	}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string          `json:"token"`
	User  *UserResponseDTO `json:"user"`
}

type UserResponseDTO struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
	TenantID string `json:"tenant_id"`
	Role     string `json:"role"`
}

func (s *AuthService) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	s.log.Infof("Login request for user: %s", req.Username)
	user, token, err := s.auc.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		s.log.Errorf("Login failed for user %s: %v", req.Username, err)
		respondError(w, http.StatusUnauthorized, "invalid username or password")
		return
	}

	s.log.Infof("Login success for user: %s", req.Username)
	respondJSON(w, http.StatusOK, &LoginResponse{
		Token: token,
		User: &UserResponseDTO{
			ID:       user.ID.String(),
			Username: user.Username,
			FullName: user.FullName,
			TenantID: user.TenantID.String(),
			Role:     user.Role,
		},
	})
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	TenantID string `json:"tenant_id"`
}

func (s *AuthService) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	tenantID, err := parseUUID(req.TenantID)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid tenant id")
		return
	}

	user := &biz.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		FullName: req.FullName,
		TenantID: tenantID,
	}

	result, err := s.auc.Register(r.Context(), user)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, &UserResponseDTO{
		ID:       result.ID.String(),
		Username: result.Username,
		FullName: result.FullName,
		TenantID: result.TenantID.String(),
		Role:     result.Role,
	})
}
