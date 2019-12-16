package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var Singleton *tk

type Claims struct {
	PlayLoad string `json:"playLoad"`
	jwt.StandardClaims
}

type tk struct {
	secret    []byte
	expiresAt time.Duration
}

func (this tk) TokenCreate(playLoad string) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		playLoad,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(this.expiresAt).Unix(),
		},
	}).SignedString(this.secret)
}

func (this tk) TokenParse(tokenString string) (*Claims, error) {
	if tokenString == "" {
		return nil, errors.New("tokenString 无效")
	}
	token, err := jwt.ParseWithClaims(tokenString, new(Claims), func(tokenString *jwt.Token) (interface{}, error) {
		return this.secret, nil
	})
	if token != nil {
		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			return claims, nil
		}
	}
	return nil, err
}

func (this *tk) TokenRefresh(tokenString string) (string, error) {
	t, err := this.TokenParse(tokenString)
	if err != nil {
		return "", err
	}
	return this.TokenCreate(t.PlayLoad)
}

func New(secret string, expiresAt time.Duration) *tk {
	if Singleton != nil {
		return Singleton
	}
	Singleton = &tk{[]byte(secret), expiresAt}
	return Singleton
}
