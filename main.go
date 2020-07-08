package main

import (
	"main/common"
	v1 "main/routes/v1"

	"github.com/gin-gonic/gin"
)

func main() {
	common.DbConnections = common.NewDbConnection()
	common.GlobalConfig = common.NewConfig()
	if err := common.GlobalConfig.Load("yaml", "./env.yml"); err != nil {
		panic("配置读取失败")
	}
	r := gin.Default()
	version1 := r.Group("v1")
	{
		version1.POST("/contract/add", v1.ContractAdd)
		version1.POST("/contract/edit", v1.ContractEdit)
		version1.POST("/contract/delete", v1.ContractDelete)
		version1.POST("/contract/index", v1.ContractIndex)
		version1.GET("/contract/index", v1.ContractIndex)

		version1.POST("/house/add", v1.HouseAdd)
		version1.POST("/house/edit", v1.HouseEdit)
		version1.POST("/house/delete", v1.HouseDelete)
		version1.POST("/house/index", v1.HousenIndex)
		version1.GET("/house/index", v1.HousenIndex)

		version1.POST("/log/add", v1.ContractLogAdd)
		version1.POST("/log/delete", v1.ContractLogDelete)
		version1.POST("/log/index", v1.ContractLogIndex)
		version1.GET("/log/index", v1.ContractLogIndex)
	}
	r.Run(":8080")
}
