package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// JWT ni yaratishda ishlatilgan secret key
var SecretKey = []byte("my_secret_key")

// Tokenni tekshiruvchi middleware
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Headerdan tokenni olish
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token kerak"})
			ctx.Abort()
			return
		}

		// "Bearer " qismi bor yoki yo‘qligini tekshiramiz
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader { // Agar token "Bearer " bilan boshlanmasa
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token noto‘g‘ri formatda"})
			ctx.Abort()
			return
		}

		// Tokenni tekshirish
		claims := &jwt.StandardClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return SecretKey, nil
		})

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Tokenni parse qilishda xatolik: " + err.Error()})
			ctx.Abort()
			return
		}

		// Token yaroqsizligi tekshirildi
		if !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token yaroqsiz"})
			ctx.Abort()
			return
		}

		// Token muddati tugagan bo‘lsa
		if claims.ExpiresAt < time.Now().Unix() {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token eskirgan"})
			ctx.Abort()
			return
		}

		// Token ichidagi `username` ni saqlaymiz
		ctx.Set("username", claims.Subject)

		// So‘rovni davom ettiramiz
		ctx.Next()
	}
}
