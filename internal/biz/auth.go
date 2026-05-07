package biz

import (
	"context"
	"errors"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type contextKey string

const (
	UserKey contextKey = "user"
)

type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TenantID string `json:"tenant_id"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrUserNotFound = errors.New("user not found")
)

type AuthRepo interface {
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	CreateUser(ctx context.Context, user *User) (*User, error)
}

type AuthUseCase struct {
	repo AuthRepo
	jwtKey []byte
	log    *log.Helper
}

func NewAuthUseCase(repo AuthRepo, logger log.Logger) *AuthUseCase {
	return &AuthUseCase{
		repo: repo,
		jwtKey: []byte("grc-secret-key"), // Should be from config
		log:    log.NewHelper(logger),
	}
}

func (uc *AuthUseCase) Login(ctx context.Context, username, password string) (*User, string, error) {
	uc.log.Infof("Attempting login for: %s", username)
	user, err := uc.repo.GetUserByUsername(ctx, username)
	if err != nil {
		uc.log.Errorf("User not found: %s", username)
		return nil, "", ErrUnauthorized
	}

	uc.log.Infof("Verifying password for: %s", username)
	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		uc.log.Errorf("Password mismatch for: %s", username)
		return nil, "", ErrUnauthorized
	}

	uc.log.Infof("Generating token for: %s", username)
	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   user.ID,
		"username":  user.Username,
		"tenant_id": user.TenantID,
		"role":      user.Role,
		"exp":       time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(uc.jwtKey)
	if err != nil {
		return nil, "", err
	}

	return user, tokenString, nil
}

func (uc *AuthUseCase) Register(ctx context.Context, user *User) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)
	user.ID = uuid.New()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	return uc.repo.CreateUser(ctx, user)
}
