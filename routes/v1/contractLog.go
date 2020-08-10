package v1

import (
	"main/common"
	"main/model"

	"github.com/gin-gonic/gin"
)

type ContractLogIn struct {
	ContractId int    `json:"contract_id"`
	ContractAt string `json:"contract_at"`
	Money      int    `json:"money"`
	Remark     string `json:"remark"`
}

func ContractLogAdd(c *gin.Context) {
	in := ContractLogIn{}
	if err := c.BindJSON(&in); err != nil {
		common.ErrorResposne(c, err.Error())
		return
	}
	house := common.DbConnections.Get("house")

	contract := model.Contract{
		ID: in.ContractId,
	}
	if err := house.First(&contract).Error; err != nil {
		common.ErrorResposne(c, err.Error())
		return
	}
	log := model.ContractLog{
		Money:      in.Money,
		Remark:     in.Remark,
		ContractId: contract.ID,
	}
	log.ContractAt, _ = common.ParseTime("2006-01-02", in.ContractAt)
	if err := house.Create(&log).Error; err != nil {
		common.ErrorResposne(c, "记录失败")
		return
	}
	common.SuccessResposne(c, gin.H{"message": "记录成功"})
}

type ContractLogDeleteIn struct {
	Id int `json:"id"`
}

func ContractLogDelete(c *gin.Context) {
	in := ContractDeleteIn{}
	if err := c.BindJSON(&in); err != nil {
		common.ErrorResposne(c, err.Error())
		return
	}
	house := common.DbConnections.Get("house")
	if house.Delete(model.ContractLog{ID: in.Id}).Error != nil {
		common.ErrorResposne(c, "删除失败")
		return
	}
	common.SuccessResposne(c, gin.H{"message": "删除成功"})

}

type ContractLogIndexIn struct {
	ContractId int `json:"contract_id,string"`
}

type ContractLogIndexOut struct {
	ID         int    `json:"id"`
	Money      int    `json:"money"`
	ContractAt string `json:"created_at"`
	Remark     string `json:"remark"`
}

func ContractLogIndex(c *gin.Context) {
	in := ContractLogIndexIn{}
	if err := c.BindJSON(&in); err != nil {
		common.ErrorResposne(c, err.Error())
		return
	}
	house := common.DbConnections.Get("house")
	logs := make([]*model.ContractLog, 0)
	if err := house.Where("contract_id=?", in.ContractId).Find(&logs).Error; err != nil {
		common.ErrorResposne(c, err.Error())
		return
	}
	var items []*ContractLogIndexOut
	for _, log := range logs {
		item := &ContractLogIndexOut{
			ID:         log.ID,
			Money:      log.Money,
			ContractAt: log.ContractAt.Format("2006-01-02"),
			Remark:     log.Remark,
		}
		items = append(items, item)
	}
	common.SuccessResposne(c, gin.H{"data": items})
}
