package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	accessSecret  = []byte("my_chat_access_secret_123")
	refreshSecret = []byte("my_chat_refresh_secret_456")
)

const (
	AccessTokenDuration  = 15 * time.Minute
	RefreshTokenDuration = 7 * 24 * time.Hour
)

type Claims struct {
	UserId    string `json:"user_id"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

func GenerateTokenPair(userId string) (accessToken string, refreshToken string, err error) {
	accessClaims := &Claims{
		UserId:    userId,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenDuration)),
			Issuer:    "my-chat",
		},
	}
	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(accessSecret)
	if err != nil {
		return "", "", err
	}
	refreshClaims := &Claims{
		UserId:    userId,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreshTokenDuration)),
			Issuer:    "my-chat",
		},
	}
	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(refreshSecret)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

// 解析并验证token
func ParseAccessToken(tokenString string) (*Claims, error) {
	return parseToken(tokenString, accessSecret, "access")
}
func ParseRefreshToken(tokenString string) (*Claims, error) {
	return parseToken(tokenString, refreshSecret, "refresh")
}
func parseToken(tokenString string, secret []byte, tokenType string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	//类型断言
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		if claims.TokenType != tokenType {
			return nil, errors.New("invalid token type")
		}
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
