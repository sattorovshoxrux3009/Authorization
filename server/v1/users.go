package v1

import (
	"log"
	"net/http"
	"time"

	"GitHub.com/sattorovshohruh3009/Authorization/server/models"
	"GitHub.com/sattorovshohruh3009/Authorization/storage/repo"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func (h *handlerV1) CreateUser(ctx *gin.Context) {
	var req models.CreateUser
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	existingUser, err := h.strg.User().Get(ctx, req.Username)
	if err == nil && existingUser != nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"error": "This username already exists",
		})
		return
	}

	user, err := h.strg.User().Create(ctx, &repo.UserCreate{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ichki xato yuz berdi :(",
		})
		return
	}

	ctx.JSON(http.StatusCreated, models.User{
		Id:       user.Id,
		Username: user.Username,
		Password: user.Password,
	})
}

// JWT token yaratish
// JWT token yaratish
func CreateJWTToken(userId int, username string) (string, error) {
	claims := jwt.MapClaims{
		"sub":     username,
		"user_id": userId, // user_id ni qo'shish
		"exp":     time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("my_secret_key"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (h *handlerV1) LoginUser(ctx *gin.Context) {
	var req models.LoginUser
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
		})
		return
	}
	user, err := h.strg.User().Get(ctx, req.Username)
	if err != nil || user == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid username or password",
		})
		return
	}
	if user.Password != req.Password {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid username or password",
		})
		return
	}
	token, err := CreateJWTToken(user.Id, user.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	err = h.strg.User().UpdateToken(ctx, user.Username, token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update token in database",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
