package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *handlerV1) GetUserSubjects(ctx *gin.Context) {
	// Middleware orqali saqlangan username ni olish
	username, exists := ctx.Get("username")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Username topilmadi"})
		return
	}

	// Username bo‘yicha foydalanuvchining ID sini olish
	user, err := h.strg.User().Get(ctx, username.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Foydalanuvchi topilmadi"})
		return
	}

	// User ID bo‘yicha fanlarni olish
	subjects, err := h.strg.Users_Subjects().GetByUserID(ctx, user.Id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Fanlar topilmadi"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"subjects": subjects})
}
