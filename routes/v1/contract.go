package v1

import (
	"main/common"
	"main/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ContracIn struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Price     int    `json:"price"`
	Month     int    `json:"month"`
	HouseId   int    `json:"house_id"`
	Id        int    `json:"id"`
}

func ContractAdd(c *gin.Context) {
	in := ContracIn{}
	if err := c.BindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResposne(err.Error()))
		return
	}
	contractModel := model.Contract{
		Price:   in.Price,
		Month:   in.Month,
		HouseId: in.HouseId,
	}
	if statrTime, err := common.ParseTime("2006-01-02 15:04:05", in.StartTime); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResposne(err.Error()))
		return
	} else {
		contractModel.StartTime = statrTime
	}
	if EndTime, err := common.ParseTime("2006-01-02 15:04:05", in.EndTime); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResposne(err.Error()))
		return
	} else {
		contractModel.EndTime = EndTime
	}
	house := common.DbConnections.Get("house")
	if !house.NewRecord(contractModel) {
		c.JSON(http.StatusBadRequest, common.ErrorResposne("新增失败"))
		return
	}
	c.JSON(http.StatusBadRequest, common.SuccessResposne("新增成功", nil))
	return

}

func ContractEdit(c *gin.Context) {
	in := ContracIn{}
	if err := c.BindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResposne(err.Error()))
		return
	}
	contractModel := model.Contract{
		Price:   in.Price,
		Month:   in.Month,
		HouseId: in.HouseId,
		Id:      in.Id,
	}
	if statrTime, err := common.ParseTime("2006-01-02 15:04:05", in.StartTime); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResposne(err.Error()))
		return
	} else {
		contractModel.StartTime = statrTime
	}
	if EndTime, err := common.ParseTime("2006-01-02 15:04:05", in.EndTime); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResposne(err.Error()))
		return
	} else {
		contractModel.EndTime = EndTime
	}
	house := common.DbConnections.Get("house")
	if house.Update(contractModel).Error != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResposne("更新失败"))
		return
	}
	c.JSON(http.StatusBadRequest, common.SuccessResposne("更新成功", nil))
	return
}

type ContractDeleteIn struct {
	Id int `json:"id"`
}

func ContractDelete(c *gin.Context) {
	in := ContractDeleteIn{}
	if err := c.BindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResposne(err.Error()))
	}
	house := common.DbConnections.Get("house")
	if house.Delete(model.Contract{Id: in.Id}).Error != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResposne("删除失败"))
		return
	}
	c.JSON(http.StatusBadRequest, common.ErrorResposne("删除失败"))
}

type ContractIndexOut struct {
	Id        int
	StartTime string
	EndTime   string
	Price     int
	Month     int
	HouseName string
}

type ContractIndexIn struct {
	HouseId int `json:"house_id"`
}

func ContractIndex(c *gin.Context) {
	in := ContractIndexIn{}
	if err := c.BindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResposne(err.Error()))
		return
	}

	house := common.DbConnections.Get("house")
	rows, err := house.Table("contract c").Joins("innher join house h on h.id=c.house_id and h.id=?", in.HouseId).Select("c.id,c.start_time,c.end_time,c.price,c.month,h.name").Rows()
	if err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResposne(err.Error()))
		return
	}
	var startTime, endTime time.Time
	var items []*ContractIndexOut
	for rows.Next() {
		item := ContractIndexOut{}
		rows.Scan(&item.Id, &startTime, &endTime, &item.Price, &item.Month, &item.HouseName)
		item.StartTime = startTime.Format("2006-01-02")
		item.EndTime = endTime.Format("2006-01-02")
		items = append(items, &item)
	}
	c.JSON(http.StatusOK, common.SuccessResposne("", items))
}
