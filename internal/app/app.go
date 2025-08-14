package app

import (
	"github.com/Rasulikus/notebook/internal/api/handler"
	"github.com/Rasulikus/notebook/internal/config"
	"github.com/Rasulikus/notebook/internal/repository"
	noteRepository "github.com/Rasulikus/notebook/internal/repository/note"
	"github.com/Rasulikus/notebook/internal/service/note"

	"github.com/gin-gonic/gin"
)

func App() *gin.Engine {
	cfg := config.LoadConfig()

	db, err := repository.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	noteRepo := noteRepository.NewRepository(db.DB)
	noteService := note.NewService(noteRepo)
	noteHanlder := handler.NewNoteHanlder(noteService)

	router := gin.Default()
	noteApi := router.Group("/notes")
	noteHanlder.RegisterNotes(noteApi)

	return router
}
