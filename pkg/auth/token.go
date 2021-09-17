package auth

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"math/rand"
	"time"
)

// TokenManager provides token issuing and decoding
type TokenManager interface {
	Issue(subject string) (string, error)
	Decode(token string) (string, error)
	Random() (string, error)
}

type JWTManager struct {
	signingKey        string
	accessTokenTTL    time.Duration
	randomTokenLength int
}

func NewJWTManager(signingKey string, accessTokenTTL time.Duration, randomTokenLength int) (*JWTManager, error) {
	if signingKey == "" {
		return nil, errors.New("empty signing key")
	}

	return &JWTManager{
		signingKey:        signingKey,
		accessTokenTTL:    accessTokenTTL,
		randomTokenLength: randomTokenLength,
	}, nil
}

func (m *JWTManager) Issue(subject string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(m.accessTokenTTL).Unix(),
		Subject:   subject,
	})

	return token.SignedString([]byte(m.signingKey))
}

func (m *JWTManager) Decode(token string) (string, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(m.signingKey), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("error get user claims from t")
	}

	return claims["sub"].(string), nil
}

func (m *JWTManager) Random() (string, error) {
	b := make([]byte, m.randomTokenLength)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
