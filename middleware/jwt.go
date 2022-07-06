package middleware

import (
	"ctp/databases"
	"ctp/model"
	"ctp/response"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

var jwtKey = []byte("a_secret_key_yyc_test")

type Claims struct {
	UserId uint
	jwt.StandardClaims
}

func ReleaseToken(user model.RegisterParam) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		UserId: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "yyc",
			Subject:   "user token",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	return token, claims, err
}

func AuthMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取 authorization header
		tokenString := ctx.GetHeader("Authorization")
		//validate token formate
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			response.Response(ctx,http.StatusOK, 401,nil, "token err")
			ctx.Abort()
			return
		}
		tokenString = tokenString[7:] //截取字符
		token, claims, err := ParseToken(tokenString)
		if err != nil || !token.Valid {
			response.Response(ctx,http.StatusOK, 401,nil, "token invalid")
			ctx.Abort()
			return
		}
		//token通过验证, 获取claims中的UserID
		userId := claims.UserId
		DB := databases.GetMysqlDB()
		var user model.RegisterParam
		DB.First(&user, userId)
		// 验证用户是否存在
		if user.ID == 0 {
			response.Response(ctx,http.StatusOK, 401,nil, "token-user is not exist")
			ctx.Abort()
			return
		}
		//用户存在 将user信息写入上下文
		ctx.Set("user", user)
		ctx.Next()
	}
}