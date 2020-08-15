package main

import (
	"main/common"
	v1 "main/routes/v1"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Options() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodOptions {
			common.OptionsResponse(c)
			c.Abort()
			return
		}
		c.Next()
	}
}
func Cose() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}
func main() {
	common.DbConnections = common.NewDbConnection()
	common.GlobalConfig = common.NewConfig()
	common.GlobalConfig.SetDefault("server.address", ":8080")
	if err := common.GlobalConfig.Load("yaml", "./env.yml"); err != nil {
		panic(err.Error())
	}
	r := gin.Default()
	r.Use(Options(), Cose())
	//卧槽,,,必须在分组之前,...为啥....
	version1 := r.Group("v1")
	{
		version1.POST("/contract/add", v1.ContractAdd)
		version1.POST("/contract/edit", v1.ContractEdit)
		version1.GET("/contract/index", v1.ContractIndex)
		version1.POST("/contract/cancel", v1.ContractCancle)
		version1.GET("/contract/contract-valid", v1.ContractValid)

		version1.GET("house/row", v1.HouseRow)
		version1.POST("/house/add", v1.HouseAdd)
		version1.POST("/house/edit", v1.HouseEdit)
		version1.POST("/house/delete", v1.HouseDelete)
		version1.POST("/house/index", v1.HousenIndex)
		version1.GET("/house/index", v1.HousenIndex)
		version1.GET("house/simple", v1.HouseSimple)

		version1.POST("/log/add", v1.ContractLogAdd)
		version1.POST("/log/delete", v1.ContractLogDelete)
		version1.POST("/log/index", v1.ContractLogIndex)
		version1.GET("/log/index", v1.ContractLogIndex)
	}
	r.Run(common.GlobalConfig.GetString("server.address"))
}
