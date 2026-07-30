package token

import (
	"fmt"

	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

type PasetoMaker struct {
	symmetricKey []byte
	paseto       *paseto.V2
}

func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d bytes", chacha20poly1305.KeySize)
	}

	maker := &PasetoMaker{
		symmetricKey: []byte(symmetricKey),
		paseto:       paseto.NewV2(),
	}
	return maker, nil
}

func (maker *PasetoMaker) CreateToken(username string, duration int64) (string, *Payload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", payload, err
	}
	token, err := maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
	return token, payload, err
}

func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}
	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	err = payload.Valid()
	if err != nil {
		return nil, fmt.Errorf("token has expired: %w", err)
	}
	return payload, nil
}
