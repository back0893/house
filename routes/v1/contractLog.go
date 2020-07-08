package v1

import (
	"main/common"
	"main/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ContractLogIn struct {
	HouseId    int     `json:"house_id"`
	HouseName  string  `json:"house_name"`
	ContractAt string  `json:"contract_at"`
	Money      float64 `json:"money"`
	Remark     string  `json:"remark"`
}

func ContractLogAdd(c *gin.Context) {
	in := ContractLogIn{}
	if err := c.BindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResposne(err.Error()))
		return
	}
	log := model.ContractLog{
		HouseId:   in.HouseId,
		HouseName: in.HouseName,
		Monoey:    int(in.Money * 100),
		Remark:    in.Remark,
	}
	log.ContractAt, _ = common.ParseTime("2006-01-02", in.ContractAt)
	house := common.DbConnections.Get("house")
	if !house.NewRecord(log) {
		c.JSON(http.StatusBadRequest, common.ErrorResposne("记录失败"))
		return
	}
	c.JSON(http.StatusBadRequest, common.SuccessResposne("记录成功", nil))
}

type ContractLogDeleteIn struct {
	Id int `json:"id"`
}

func ContractLogDelete(c *gin.Context) {
	in := ContractDeleteIn{}
	if err := c.BindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResposne(err.Error()))
		return
	}
	house := common.DbConnections.Get("house")
	if house.Delete(model.ContractLog{Id: in.Id}).Error != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResposne("删除失败"))
		return
	}
	c.JSON(http.StatusBadRequest, common.ErrorResposne("删除成功"))

}

type ContractLogIndexIn struct {
	ContractId int `json:"contract_id"`
}

type ContractLogIndexOut struct {
	Id         int
	HouseName  string
	Money      int
	ContractAt string
	remark     string
}

func ContractLogIndex(c *gin.Context) {
	in := ContractLogIndexIn{}
	if err := c.BindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResposne(err.Error()))
		return
	}
	house := common.DbConnections.Get("house")
	logs := make([]*model.ContractLog, 0)
	if err := house.Where("contract_id=?", in.ContractId).Find(&logs).Error; err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResposne(err.Error()))
		return
	}
	var items []*ContractLogIndexOut
	for _, log := range logs {
		item := &ContractLogIndexOut{
			Id:         log.Id,
			HouseName:  log.HouseName,
			Money:      log.Monoey,
			ContractAt: log.ContractAt.Format("2006-01-02"),
			remark:     log.Remark,
		}
		items = append(items, item)
	}
	c.JSON(http.StatusBadRequest, common.SuccessResposne("", items))
}
