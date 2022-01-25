package service

import (
	"fmt"
	"os"
	"time"

	"github.com/anna-osipova/go-wordle/models"
	"github.com/golang-jwt/jwt"
)

//jwt service
type JWTService interface {
	GenerateToken(session *models.Session) string
	ValidateToken(token string) (*jwt.Token, error)
	GetSecretKey() string
}
type CustomClaims struct {
	SessionId string `json:"session_id"`
	jwt.StandardClaims
}

type jwtServices struct {
	secretKey string
	issure    string
}

func JWTAuthService() JWTService {
	return &jwtServices{
		secretKey: getSecretKey(),
		issure:    "Anna",
	}
}

func getSecretKey() string {
	secret := os.Getenv("SECRET")
	if secret == "" {
		secret = "secret"
	}
	return secret
}

func (service *jwtServices) GetSecretKey() string {
	return service.secretKey
}

func (service *jwtServices) GenerateToken(session *models.Session) string {
	claims := &CustomClaims{
		session.ID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
			Issuer:    service.issure,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//encoded string
	t, err := token.SignedString([]byte(service.secretKey))
	if err != nil {
		panic(err)
	}
	return t
}

func (service *jwtServices) ValidateToken(encodedToken string) (*jwt.Token, error) {
	return jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, isvalid := token.Method.(*jwt.SigningMethodHMAC); !isvalid {
			return nil, fmt.Errorf("Invalid token", token.Header["alg"])
		}
		return []byte(service.secretKey), nil
	})
}
