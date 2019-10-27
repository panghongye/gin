package jwt

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// var Jwt *tk

type tk struct {
	secret    []byte
	expiresAt int64
}

func (this *tk) TokenCreate(playLoad interface{}) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{ // jwt.StandardClaims
		"exp":      this.expiresAt + time.Now().UnixNano(),
		"playLoad": playLoad,
	})
	tokenString, err := token.SignedString(this.secret)
	return tokenString, err
}

func (this *tk) TokenParse(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return this.secret, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

func (this *tk) TokenRefresh(tokenString string) (string, error) {
	t, err := this.TokenParse(tokenString)
	if err != nil {
		return "", err
	}
	return this.TokenCreate(t["playLoad"])
}

// expiresAt/UnixNano
func New(secret string, expiresAt int64) *tk {
	Jwt := &tk{[]byte(secret), expiresAt}
	return Jwt
}
