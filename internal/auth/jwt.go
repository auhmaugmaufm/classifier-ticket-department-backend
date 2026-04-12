package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTManager struct {
	SecretKey string
	TTL       time.Duration
}

type Claims struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func NewJWTManager(secretKey string, ttlHours int) *JWTManager {
	return &JWTManager{SecretKey: secretKey, TTL: time.Duration(ttlHours) * time.Hour}
}

func (m *JWTManager) GenerateToken(ID string, email string) (string, error) {
	claims := &Claims{
		ID:    ID,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.TTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.SecretKey))
}

func (m *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(m.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}
