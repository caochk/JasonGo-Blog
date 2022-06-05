package jwtUtils

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"my_blog/utils/respUtils"
)

const (
	//HS256 signed key
	SIGNED_KEY = "jsongoforpassionate" // 签名用密钥（对header+claims的组合字符串进行加密后时用）
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
func ParseToken(signedToken string) (*jwt.Token, error) { // jwt.Parse函数:解析，验证、验证签名并返回分析的令牌
	token, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) { // Parse函数应该只是验证一个token（需要借助签名用密钥），并没有进行解密？同时将前端传过来的字符串类型的令牌转换为了token类型
		return []byte(SIGNED_KEY), nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); !ok {
			return nil, errors.New("无法处理此令牌")
		} else {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, err
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				return nil, errors.New("令牌已过期或尚未激活")
			} else {
				return nil, errors.New("无法处理此令牌")
			}
		}
	}
	if !token.Valid {
		return nil, errors.New("令牌无效")
	}
	return token, nil
}

// Token2Claims 将token中的claims提取出来
func Token2Claims(token jwt.Token) (jwt.MapClaims, *respUtils.Resp) {
	var resp respUtils.Resp
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		respond := resp.NewResp(respUtils.TOKEN_ERR_CODE, "token类型转换错误")
		return nil, respond
	}
	return claims, nil
}

//func BuildClaims(...interface{}) jwt.Claims {
//
//}
