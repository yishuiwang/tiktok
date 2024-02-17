package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"tiktok/models"
	"time"
)

type Claims struct {
	UserId int64
	jwt.StandardClaims
}

// ReleaseToken 生成token
func ReleaseToken(user models.UserLogin) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		UserId: user.Id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "tiktok",
			Subject:   "user token",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("tiktok"))
}

// ParseToken 解析token
func ParseToken(tokenString string) (*Claims, bool) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("tiktok"), nil
	})
	if err != nil {
		return nil, false
	}
	claims, ok := token.Claims.(*Claims)
	return claims, ok
}

// JWTMiddleWare 鉴权中间件，鉴权并设置user_id
func JWTMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.Query("token")
		if tokenStr == "" {
			tokenStr = c.PostForm("token")
		}

		// 从header中获取token
		if tokenStr == "" {
			c.JSON(200, models.CommonResponse{
				StatusCode: 401,
				StatusMsg:  "token is empty",
			})
			c.Abort()
			return
		}

		// 解析token
		claims, ok := ParseToken(tokenStr)
		if !ok {
			c.JSON(200, models.CommonResponse{
				StatusCode: 401,
				StatusMsg:  "token is invalid",
			})
			c.Abort()
			return
		}

		// 验证token是否过期
		if time.Now().Unix() > claims.ExpiresAt {
			c.JSON(200, models.CommonResponse{
				StatusCode: 401,
				StatusMsg:  "token is expired",
			})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserId)
		c.Next()
	}
}
