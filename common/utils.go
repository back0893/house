package common

import (
	"time"

	"github.com/gin-gonic/gin"
)

func ErrorResposne(msg string)gin.H{
	return gin.H{"message":msg}
}
func SuccessResposne(msg string,data interface{})gin.H{
	return gin.H{"message":msg,"data":data}
}

func GetLoc()*time.Location{
	lo:=GlobalConfig.GetString("location")
	if lo==""{
		lo="Local"
	}
	location,err:=time.LoadLocation(lo)
	if err!=nil{
		location,_:=time.LoadLocation("UTC")
		return location
	}
	return location
}

func ParseTime(layout,parseTime string)(time.Time,error){
	return time.ParseInLocation(layout,parseTime,GetLoc())
}