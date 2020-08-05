package common

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func ErrorResposne(c *gin.Context, msg string) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusBadRequest, gin.H{"message": msg})
}
func SuccessResposne(c *gin.Context, data interface{}) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, data)
}

func OptionsResponse(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Content-Type")
	c.Header("Access-Control-Allow-Methods", "POST,GET,OPTIONS,DELETE")
	c.JSON(http.StatusOK, nil)
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
