package models

import "github.com/dgrijalva/jwt-go"

// CustomClaims 加密解密都需要这个结构.
type CustomClaims struct {
	// 额外的payload
	ID         uint
	NickName   string
	AutorityId uint // role

	jwt.StandardClaims // 继承基本结构
}
