package jwt

import (
	"github.com/DmitriySolopenkov/distribCalc.rev2/pkg/config"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type CustomClaims struct {
	ID int
	jwt.StandardClaims
}

type Token struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}

type TokenJwt struct {
}

func New() *TokenJwt {
	return &TokenJwt{}
}

func (j *TokenJwt) CreateUserToken(ID int) (*Token, error) {
	expiresAt := time.Now().Add(time.Hour * 24 * time.Duration(config.Get().JwtExpires)).Unix()
	claims := &CustomClaims{
		ID,
		jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Get().JwtSecret))

	return &Token{Token: tokenString, ExpiresAt: expiresAt}, err
}

func (j *TokenJwt) ParseToken(token string) (*CustomClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Get().JwtSecret), nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*CustomClaims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}

func (j *TokenJwt) JwtClaims(c *gin.Context) (*CustomClaims, error) {
	token := c.GetHeader("Authorization")
	claims, err := j.ParseToken(token)
	return claims, err
}

func (j *TokenJwt) JwtUserId(c *gin.Context) int {
	claims, err := j.JwtClaims(c)
	if err != nil {
		return 0
	}
	return claims.ID
}
