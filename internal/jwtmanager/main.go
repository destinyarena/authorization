package jwtmanager

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

type (
	jwtManager struct {
		SigningMethod jwt.SigningMethod
		Secret        string
		Logger        *logrus.Logger
	}

	// JWTManager handles are the exported functions replated to jwtManager
	JWTManager interface {
		Sign(jwt.Claims) (string, error)
		Decrypt(string, jwt.Claims) error
	}
)

// NewDefault generates config and returns a new JWTManager
func NewDefault(logger *logrus.Logger) (JWTManager, error) {
	cfg, err := newEnvConfig()
	if err != nil {
		return nil, err
	}

	jwt, err := New(cfg, logger)
	return jwt, err
}

// New returns a new jwt manager
func New(cfg *Config, logger *logrus.Logger) (JWTManager, error) {
	if len(cfg.Secret) == 0 {
		return nil, errors.New("Invalid secret")
	}

	return &jwtManager{
		SigningMethod: jwt.GetSigningMethod("HS256"),
		Secret:        cfg.Secret,
		Logger:        logger,
	}, nil
}
