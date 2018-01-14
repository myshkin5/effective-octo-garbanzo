package identity

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/myshkin5/effective-octo-garbanzo/logs"
)

type Validator struct {
	publicKeys map[string]*rsa.PublicKey
}

type CustomClaims struct {
	Org string `json:"custom:org"`
	jwt.StandardClaims
}

func (c *CustomClaims) Valid() error {
	err := c.StandardClaims.Valid()
	if err != nil {
		return err
	}

	if len(c.Org) == 0 {
		return errors.New("org not specified in claims")
	}

	return nil
}

const bearerPrefix = "bearer "

func NewValidator(publicKeys map[string]*rsa.PublicKey) *Validator {
	return &Validator{
		publicKeys: publicKeys,
	}
}

func (v *Validator) IsValid(authHeader string) (ok bool, org string) {
	if !strings.HasPrefix(strings.ToLower(authHeader), bearerPrefix) {
		logs.Logger.Infof("Authentication header lacks %sprefix", bearerPrefix)
		return false, ""
	}

	tokenString := authHeader[len(bearerPrefix):]
	var claims CustomClaims
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		keyId, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("key id (kid) is not a string in headers: %v", token.Header)
		}

		publicKey, ok := v.publicKeys[keyId]
		if !ok {
			return nil, fmt.Errorf("no public key found for key id (kid) %v", keyId)
		}

		return publicKey, nil
	})
	if err != nil {
		logs.Logger.Infof("Error parsing authentication header, %v", err)
		return false, ""
	}

	return true, claims.Org
}
