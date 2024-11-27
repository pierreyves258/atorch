package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pierreyves258/atorch"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Current float64 `json:"current" binding:"required"`
	Cutoff  float64 `json:"cutoff" binding:"required"`
}

func Start(dl24 *atorch.PX100) gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := Config{}

		if err := c.BindJSON(&cfg); err != nil {
			log.Error().Err(err).Msg("Cannot bind json")
			return
		}

		err := dl24.SetData(atorch.Reset, nil, true)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err})
			return
		}

		err = dl24.SetData(atorch.SetCurrent, cfg.Current, false)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err})
			return
		}

		err = dl24.SetData(atorch.SetCutoff, cfg.Cutoff, false)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err})
			return
		}

		err = dl24.SetData(atorch.SetOutput, true, false)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err})
			return
		}

		c.JSON(http.StatusOK, gin.H{"result": "running"})
	}
}

func Reset(dl24 *atorch.PX100) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := dl24.SetData(atorch.Reset, nil, true)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err})
			return
		}

		c.JSON(http.StatusOK, gin.H{"result": "ok"})
	}
}

func SetCurrent(dl24 *atorch.PX100) gin.HandlerFunc {
	return func(c *gin.Context) {
		current := c.Param("value")

		fCurrent, err := strconv.ParseFloat(current, 64)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err})
			return
		}

		log.Debug().Msgf("SetCurrent %+v", current)
		err = dl24.SetData(atorch.SetCurrent, fCurrent, false)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err})
			return
		}

		curr, err := dl24.GetData(atorch.GetCurrentLimit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err})
			return
		}

		c.JSON(http.StatusOK, gin.H{"result": "ok", "current": curr})
	}
}
