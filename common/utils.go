package common

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func ErrorResposne(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, gin.H{"message": msg})
}
func SuccessResposne(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

func GetLoc() *time.Location {
	lo := GlobalConfig.GetString("location")
	if lo == "" {
		lo = "Local"
	}
	location, err := time.LoadLocation(lo)
	if err != nil {
		location, _ := time.LoadLocation("UTC")
		return location
	}
	return location
}

func ParseTime(layout, parseTime string) (time.Time, error) {
	return time.ParseInLocation(layout, parseTime, GetLoc())
}

func FindIn(target int, sl []int) bool {
	for _, n := range sl {
		if n == target {
			return true
		}
	}
	return false
}
