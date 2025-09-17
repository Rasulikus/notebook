package app

import (
	"time"

	"github.com/Rasulikus/notebook/internal/api/handler"
	"github.com/Rasulikus/notebook/internal/api/middleware"
	"github.com/Rasulikus/notebook/internal/config"
	"github.com/Rasulikus/notebook/internal/repository"
	noteRepository "github.com/Rasulikus/notebook/internal/repository/note"
	"github.com/Rasulikus/notebook/internal/repository/session"
	"github.com/Rasulikus/notebook/internal/repository/user"
	"github.com/Rasulikus/notebook/internal/service/auth"
	"github.com/Rasulikus/notebook/internal/service/note"

	"github.com/gin-gonic/gin"
)

const (
	accessTTL  = 15 * time.Minute
	refreshTTL = 6 * time.Hour
)

func App() *gin.Engine {
	cfg := config.LoadConfig()

	db, err := repository.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	userRepo := user.NewRepository(db.DB)
	sessinoRepo := session.NewRepository(db.DB)
	jwtService := auth.NewTokenManager([]byte(cfg.Auth.Secret), accessTTL, refreshTTL, sessinoRepo)
	authService := auth.NewService(userRepo, jwtService)
	authHandler := handler.NewAuthHanlder(authService, jwtService, refreshTTL, false)

	noteRepo := noteRepository.NewRepository(db.DB)
	noteService := note.NewService(noteRepo)
	noteHanlder := handler.NewNoteHandler(noteService)

	router := gin.Default()
	authApi := router.Group("/auth")
	{
		authApi.POST("/register", authHandler.Register)
		authApi.POST("/login", authHandler.Login)
		authApi.POST("/refresh", authHandler.Refresh)
		authApi.POST("/logout", authHandler.Logout)
	}

	noteApi := router.Group("/notes", middleware.AuthMiddleware(jwtService))
	{
		noteApi.POST("", noteHanlder.Create)
		noteApi.GET("", noteHanlder.List)
		noteApi.GET("/:id", noteHanlder.GetByID)
		noteApi.DELETE("/:id", noteHanlder.DeleteByID)
	}

	return router
}
