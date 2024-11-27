package router

import (
	"github.com/gin-gonic/gin"
	"github.com/pierreyves258/atorch"
	"github.com/pierreyves258/atorch/cmd/dl24_api/handlers"
)

func Init(dl24 *atorch.PX100) *gin.Engine {

	router := gin.New()

	router.GET("/config", handlers.GetConfig(dl24))
	router.POST("/current/:value", handlers.SetCurrent(dl24))
	router.POST("/reset", handlers.Reset(dl24))
	router.POST("/start", handlers.Start(dl24))

	return router
}
