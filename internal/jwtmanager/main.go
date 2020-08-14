package jwtmanager

import (
	"crypto/rsa"
	"io/ioutil"

	"github.com/dgrijalva/jwt-go"
)

type (
	jwtManager struct {
		SigningMethod jwt.SigningMethod
		VerifyKey     *rsa.PublicKey
		SignKey       *rsa.PrivateKey
	}

	// JWTManager handles are the exported functions replated to jwtManager
	JWTManager interface {
		Sign(interface{}) (string, error)
		Decrypt(string, interface{}) error
	}
)

// NewDefault generates config and returns a new JWTManager
func NewDefault() (JWTManager, error) {
	cfg, err := newEnvConfig()
	if err != nil {
		return nil, err
	}

	jwt, err := New(cfg)
	return jwt, err
}

// New returns a new jwt manager
func New(cfg *Config) (JWTManager, error) {
	signBytes, err := ioutil.ReadFile(cfg.PrivKeyPath)
	if err != nil {
		return nil, err
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		return nil, err
	}

	verifyBytes, err := ioutil.ReadFile(cfg.PubKeyPath)
	if err != nil {
		return nil, err
	}

	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		return nil, err
	}

	return &jwtManager{
		SigningMethod: jwt.GetSigningMethod("RS256"),
		VerifyKey:     verifyKey,
		SignKey:       signKey,
	}, nil
}
