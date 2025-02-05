package server

import (
	v1 "GitHub.com/sattorovshohruh3009/Authorization/server/v1"
	"GitHub.com/sattorovshohruh3009/Authorization/storage"
	"github.com/gin-gonic/gin"
)

type Options struct {
	Strg storage.StorageI
}

func NewServer(opts *Options) *gin.Engine {
	router := gin.New()

	//ROUTER
	handler := v1.New(&v1.HandlerV1{
		Strg: opts.Strg,
	})

	//Users
	router.POST("/v1/users", handler.CreateUser)
	router.POST("/v1/login", handler.LoginUser)
	// Himoyalangan API-lar (Middleware ishlaydi)
	protected := router.Group("/v1")
	protected.Use(handler.AuthMiddleware()) // JWT tokenni tekshirish
	{
		protected.GET("/user-subjects", handler.GetUserSubjects) // Foydalanuvchining fanlari

	}
	return router
}
