package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/weibh/openClusterClaw/config"
	"github.com/weibh/openClusterClaw/internal/model"
)

var (
	// ErrInvalidToken is returned when the token is invalid
	ErrInvalidToken = errors.New("invalid token")
	// ErrExpiredToken is returned when the token has expired
	ErrExpiredToken = errors.New("token has expired")
)

// Claims represents JWT claims
type Claims struct {
	UserID   string      `json:"user_id"`
	Username string      `json:"username"`
	TenantID string      `json:"tenant_id"`
	Role     model.UserRole `json:"role"`
	jwt.RegisteredClaims
}

// JWTService handles JWT token generation and validation
type JWTService struct {
	secret      []byte
	accessExp   time.Duration
	refreshExp  time.Duration
}

// NewJWTService creates a new JWT service
func NewJWTService(cfg *config.Config) *JWTService {
	accessExp := time.Duration(cfg.JWT.ExpireTime) * time.Second
	refreshExp := 7 * 24 * time.Hour // 7 days for refresh token

	return &JWTService{
		secret:      []byte(cfg.JWT.Secret),
		accessExp:   accessExp,
		refreshExp:  refreshExp,
	}
}

// GenerateAccessToken generates a new access token
func (s *JWTService) GenerateAccessToken(user *model.User) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:   user.ID,
		Username: user.Username,
		TenantID: user.TenantID,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessExp)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

// GenerateRefreshToken generates a new refresh token
func (s *JWTService) GenerateRefreshToken(user *model.User) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:   user.ID,
		Username: user.Username,
		TenantID: user.TenantID,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshExp)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

// ValidateToken validates a JWT token and returns claims
func (s *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// GetAccessExpiration returns the access token expiration duration
func (s *JWTService) GetAccessExpiration() int64 {
	return int64(s.accessExp.Seconds())
}

// GetRefreshExpiration returns the refresh token expiration duration
func (s *JWTService) GetRefreshExpiration() int64 {
	return int64(s.refreshExp.Seconds())
}
