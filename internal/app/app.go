package app

import (
	"time"

	"github.com/Rasulikus/notebook/internal/api/handler"
	"github.com/Rasulikus/notebook/internal/api/middleware"
	"github.com/Rasulikus/notebook/internal/config"
	"github.com/Rasulikus/notebook/internal/repository"
	noteRepository "github.com/Rasulikus/notebook/internal/repository/note"
	"github.com/Rasulikus/notebook/internal/repository/session"
	tagRepository "github.com/Rasulikus/notebook/internal/repository/tag"
	"github.com/Rasulikus/notebook/internal/repository/user"
	"github.com/Rasulikus/notebook/internal/service/auth"
	"github.com/Rasulikus/notebook/internal/service/note"
	"github.com/Rasulikus/notebook/internal/service/tag"

	"github.com/gin-gonic/gin"
)

const (
	accessTTL  = 15 * time.Minute
	refreshTTL = 24 * time.Hour
)

func App() *gin.Engine {
	cfg := config.LoadConfig()

	db, err := repository.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	userRepo := user.NewRepository(db.DB)
	sessinoRepo := session.NewRepository(db.DB)
	authService := auth.NewService(userRepo, auth.TokenConfig{
		Secret:      []byte(cfg.Auth.Secret),
		AccessTTL:   accessTTL,
		RefreshTTL:  refreshTTL,
		SessionRepo: sessinoRepo,
	})
	authHandler := handler.NewAuthHandler(authService, refreshTTL, false)

	tagRepo := tagRepository.NewRepository(db.DB)
	tagService := tag.NewService(tagRepo)
	tagHandler := handler.NewTagHandler(tagService)

	noteRepo := noteRepository.NewRepository(db.DB)
	noteService := note.NewService(noteRepo, tagRepo)
	noteHandler := handler.NewNoteHandler(noteService)

	router := gin.Default()
	authApi := router.Group("/auth")
	{
		authApi.POST("/register", authHandler.Register)
		authApi.POST("/login", authHandler.Login)
		authApi.POST("/refresh", authHandler.Refresh)
		authApi.POST("/logout", authHandler.Logout)
	}

	noteApi := router.Group("/notes", middleware.AuthMiddleware(authService))
	{
		noteApi.POST("", noteHandler.Create)
		noteApi.GET("", noteHandler.List)
		noteApi.GET("/:id", noteHandler.GetByID)
		noteApi.PATCH("/:id", noteHandler.UpdateByID)
		noteApi.DELETE("/:id", noteHandler.DeleteByID)
	}

	tagApi := router.Group("/tags", middleware.AuthMiddleware(authService))
	{
		tagApi.POST("", tagHandler.Create)
		tagApi.GET("", tagHandler.List)
		tagApi.GET("/:id", tagHandler.GetByID)
		tagApi.PATCH("/:id", tagHandler.UpdateByID)
		tagApi.DELETE("/:id", tagHandler.DeleteByID)
	}

	return router
}
