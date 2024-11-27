package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pierreyves258/atorch"
)

func GetConfig(dl24 *atorch.PX100) gin.HandlerFunc {
	return func(c *gin.Context) {
		voltage, err := dl24.GetData(atorch.GetVoltage)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err})
			return
		}
		current, err := dl24.GetData(atorch.GetCurrent)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err})
			return
		}
		ison, err := dl24.GetData(atorch.GetIsOn)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err})
			return
		}

		c.JSON(http.StatusOK, gin.H{"voltage": voltage, "current": current, "ison": ison})
	}
}
