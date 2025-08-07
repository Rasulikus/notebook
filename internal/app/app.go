package app

import (
	"github.com/Rasulikus/notebook/internal/config"
	"github.com/gin-gonic/gin"
)

func App() *gin.Engine {
	cfg := config.LoadConfig()

	//db := db.NewDb(cfg)
	router := gin.Default()

	return router
}
