package tokenizer

import (
	"ara-server/internal/infrastructure/configuration"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type Tokenizer struct {
	config *configuration.AppConfig
}

func NewTokenizer(appConfig *configuration.AppConfig) *Tokenizer {
	return &Tokenizer{
		config: appConfig,
	}
}

func (t *Tokenizer) GenerateToken(payload map[string]interface{}) (string, error) {
	claim := jwt.MapClaims{}
	for k, v := range payload {
		claim[k] = v
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString([]byte(t.config.Auth.JWTSecret))
}

func (t *Tokenizer) VerifyAndParseJWTToken(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(t.config.Auth.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		result := map[string]interface{}{}
		for k, v := range claims {
			result[k] = v
		}
		return result, nil
	}

	return nil, err
}
