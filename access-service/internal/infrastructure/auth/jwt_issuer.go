// Package auth implements the access.TokenIssuer driven port with JWT
// (library-docs/07-api/authentication.md).
package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTIssuer struct {
	secret []byte
	expiry time.Duration
}

func NewJWTIssuer(secret string, expiry time.Duration) *JWTIssuer {
	return &JWTIssuer{secret: []byte(secret), expiry: expiry}
}

// Issue signs a JWT with the Administrator's ID as `sub` — there is no `role`
// claim, since v1 has a single role and no RBAC (02-domain/domain-map.md).
func (j *JWTIssuer) Issue(administratorID string) (string, int, error) {
	now := time.Now().UTC()
	claims := jwt.MapClaims{
		"sub": administratorID,
		"iat": now.Unix(),
		"exp": now.Add(j.expiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(j.secret)
	if err != nil {
		return "", 0, err
	}

	return signed, int(j.expiry.Seconds()), nil
}
