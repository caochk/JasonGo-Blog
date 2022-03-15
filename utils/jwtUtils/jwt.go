package jwtUtils

import "github.com/golang-jwt/jwt"

const (
	//HS256 signed key
	SIGNED_KEY = "jsongoforpassionate"
)

// BuildToken 创建带签名令牌
func BuildToken(claims jwt.Claims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)                  // 带着自定义的声明创建了一个token
	if signedToken, err := token.SignedString([]byte(SIGNED_KEY)); err != nil { // 用自定义密钥加密
		return ""
	} else {
		return signedToken
	}
}

// ParseToken 解析令牌
func ParseToken(signedToken string) (*jwt.Claims, error) {
	token, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(SIGNED_KEY), nil
	})
	// 未完
}

//func BuildClaims(...interface{}) jwt.Claims {
//
//}
