package v1

import (
	"context"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// JWT ni yaratishda ishlatilgan secret key
var SecretKey = []byte("my_secret_key")

func (h *handlerV1) AuthMiddleware() gin.HandlerFunc {
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
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			return SecretKey, nil
		})

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token error " + err.Error()})
			ctx.Abort()
			return
		}

		// Token yaroqsizligi tekshirildi
		if !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token yaroqsiz"})
			ctx.Abort()
			return
		}

		userIDFloat, ok := claims["user_id"].(float64) // JWT float64 qaytarishi mumkin
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token tarkibida user_id yo‘q yoki noto‘g‘ri formatda"})
			ctx.Abort()
			return
		}

		userID := int(userIDFloat)
		storedToken, err := h.strg.User().GetToken(context.TODO(), userID)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Foydalanuvchi topilmadi yoki token mavjud emas"})
			ctx.Abort()
			return
		}
		// Token database dagi token bilan mos kelishini tekshirish
		if storedToken != tokenString {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Noto‘g‘ri token"})
			ctx.Abort()
			return
		}

		// Token ichidagi `username`, `user_id` va `role` ni saqlaymiz
		ctx.Set("username", claims["sub"])
		ctx.Set("user_id", claims["user_id"]) // user_id
		ctx.Set("role", claims["role"])       // role

		// So‘rovni davom ettiramiz
		ctx.Next()
	}
}
