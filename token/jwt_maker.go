package token

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

// JwtMaker is a JSON Web Token maker.
type JwtMaker struct {
	secretKey string
}

func NewJwtMaker(secretKey string) (Maker, error) {
	if len(secretKey) < 32 {
		return nil, fmt.Errorf("Invalid key size: must be at least 32 characters")
	}
	return &JwtMaker{secretKey}, nil
}
func (maker *JwtMaker) CreateToken(username string, duration int64) (string, *Payload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", payload, err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(maker.secretKey))
	return token, payload, err
}

// VerifyToken checks if the token is valid or not.
func (maker *JwtMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("invalid token")
		}
		return []byte(maker.secretKey), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && verr.Errors == jwt.ValidationErrorExpired {
			return nil, fmt.Errorf("token has expired: %w", err)
		}
		return nil, fmt.Errorf("invalid token")
	}
	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, fmt.Errorf("invalid token")
	}
	return payload, nil
}
