package jwt

import (
	"errors"
	"github.com/chuji555/homework-system/pkg/errcode"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

// Token载荷（存储用户核心信息）
type Claims struct {
	UserID     int64  `json:"user_id"`
	Username   string `json:"username"`
	Role       string `json:"role"`
	Department string `json:"department"`
	jwt.RegisteredClaims
}

// 生成双Token
func GenerateTokens(userID int64, username, role, department string) (accessToken, refreshToken string, err error) {
	// 生成AccessToken
	accessClaims := Claims{
		UserID:     userID,
		Username:   username,
		Role:       role,
		Department: department,
		RegisteredClaims: jwt.RegisteredClaims{
			// 过期时间：当前时间+配置的有效期
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(viper.GetInt("jwt.access_expire")))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(viper.GetString("jwt.secret")))
	if err != nil {
		return "", "", err
	}
	// 生成RefreshToken
	refreshClaims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(viper.GetInt("jwt.refresh_expire")))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(viper.GetString("jwt.secret")))
	if err != nil {
		return "", "", err
	}
	return
}

// 解析AccessToken
func ParseAccessToken(tokenString string) (*Claims, errcode.ErrCode) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("jwt.secret")), nil
	})
	if err != nil {
		// Token过期
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errcode.TokenExpired
		}
		return nil, errcode.AuthError
	}
	// 验证Token有效性
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, errcode.Success
	}
	return nil, errcode.AuthError
}

// 解析RefreshToken
func ParseRefreshToken(tokenString string) (int64, errcode.ErrCode) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("jwt.secret")), nil
	})
	if err != nil || !token.Valid {
		return 0, errcode.AuthError
	}
	claims, _ := token.Claims.(*Claims)
	return claims.UserID, errcode.Success
}
