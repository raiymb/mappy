package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/raiymb/mappy/config"
)

// CustomClaims gets embedded in every token we sign.
type CustomClaims struct {
	UID  string `json:"uid"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

// Pair contains the freshly‑signed access/refresh tokens.
type Pair struct {
	Access  string
	Refresh string
}

// NewPair issues a 15‑min access & 30‑day refresh token.
func NewPair(uid, role string, cfg config.JWT) (Pair, error) {
	now := time.Now().UTC()

	accessClaims := CustomClaims{
		UID:  uid,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(cfg.AccessTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        util.UUID(), // JTI
		},
	}
	refreshClaims := CustomClaims{
		UID:  uid,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(cfg.RefreshTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        util.UUID(),
		},
	}

	access, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(cfg.Secret))
	if err != nil {
		return Pair{}, err
	}
	refresh, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(cfg.Secret))
	if err != nil {
		return Pair{}, err
	}
	return Pair{Access: access, Refresh: refresh}, nil
}

// Parse verifies signature + expiration and returns custom claims.
func Parse(tokenStr string, secret string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(t *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}
	return claims, nil
}
