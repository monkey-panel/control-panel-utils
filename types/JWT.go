package types

import (
	"time"

	"github.com/monkey-panel/control-panel-utils/utils"

	"github.com/golang-jwt/jwt/v5"
)

type (
	JWT    string
	Claims struct {
		ID ID `json:"id"`
		jwt.RegisteredClaims
	}
	RefreshToken struct {
		Token          JWT  `json:"token"`
		ExpirationTime Time `json:"expiration"`
	}
)

func NewJWT(userID ID) (token *RefreshToken) {
	return CreateJWT(Claims{ID: userID}, time.Duration(int64(time.Hour)*utils.Config().JWTTimeout))
}

func CreateJWT(claims Claims, newTime time.Duration) (token *RefreshToken) {
	expirationTime := time.Now().Add(newTime)

	claims.RegisteredClaims = jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(expirationTime)}

	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims).SignedString([]byte(utils.Config().JWTKey))
	if err != nil {
		return nil
	}

	return &RefreshToken{JWT(tokenString), NewTime(expirationTime)}
}

// JWT to string
func (j JWT) String() string { return string(j) }

func (j JWT) Info() *Claims {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(j.String(), claims, func(token *jwt.Token) (any, error) {
		return []byte(utils.Config().JWTKey), nil
	})

	if token == nil || (token != nil && !token.Valid) || err != nil {
		return nil
	}

	return claims
}

// Refresh Token
func (j *JWT) Refresh(newTime time.Duration) *RefreshToken {
	claims := j.Info()
	if claims == nil {
		return nil
	}

	if token := CreateJWT(*claims, newTime); token != nil {
		*j = JWT(token.Token)
		return token
	}

	return nil
}
