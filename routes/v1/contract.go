package v1

import (
	"main/common"
	"main/model"
	"time"

	"github.com/gin-gonic/gin"
)

type ContracIn struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Price     int    `json:"price"`
	Month     int    `json:"month"`
	HouseId   int    `json:"house_id,string"`
	Id        int    `json:"id"`
	CardName  string `json:"card_name"`
	CardNum   string `json:"card_num"`
}

func ContractAdd(c *gin.Context) {
	in := ContracIn{}
	if err := c.BindJSON(&in); err != nil {
		common.ErrorResposne(c, err.Error())
		return
	}
	contractModel := model.Contract{
		Price:    in.Price,
		Month:    in.Month,
		HouseId:  in.HouseId,
		CardName: in.CardName,
		CardNum:  in.CardNum,
	}
	if statrTime, err := common.ParseTime("2006-01-02", in.StartTime); err != nil {
		common.ErrorResposne(c, err.Error())
		return
	} else {
		contractModel.StartTime = statrTime
	}
	if EndTime, err := common.ParseTime("2006-01-02", in.EndTime); err != nil {
		common.ErrorResposne(c, err.Error())
		return
	} else {
		contractModel.EndTime = EndTime
	}
	house := common.DbConnections.Get("house")
	if err := house.Create(&contractModel).Error; err != nil {
		common.ErrorResposne(c, "新增失败")
		return
	}
	common.SuccessResposne(c, gin.H{"message": "新增成功"})
	return

}

func ContractEdit(c *gin.Context) {
	in := ContracIn{}
	if err := c.BindJSON(&in); err != nil {
		common.ErrorResposne(c, err.Error())
		return
	}
	contractModel := model.Contract{
		Price:    in.Price,
		Month:    in.Month,
		HouseId:  in.HouseId,
		ID:       in.Id,
		CardName: in.CardName,
		CardNum:  in.CardNum,
	}
	if in.StartTime != "" {
		if statrTime, err := common.ParseTime("2006-01-02", in.StartTime); err != nil {
			common.ErrorResposne(c, err.Error())
			return
		} else {
			contractModel.StartTime = statrTime
		}
	}
	if in.EndTime != "" {
		if EndTime, err := common.ParseTime("2006-01-02", in.EndTime); err != nil {
			common.ErrorResposne(c, err.Error())
			return
		} else {
			contractModel.EndTime = EndTime
		}
	}
	house := common.DbConnections.Get("house")
	if house.Model(&contractModel).Update(contractModel).Error != nil {
		common.ErrorResposne(c, "修改成功")
		return
	}
	common.SuccessResposne(c, gin.H{"message": "修改成功"})
	return
}

type ContractDeleteIn struct {
	Id int `json:"id"`
}

func ContractDelete(c *gin.Context) {
	in := ContractDeleteIn{}
	if err := c.BindJSON(&in); err != nil {
		common.ErrorResposne(c, err.Error())
		return
	}
	house := common.DbConnections.Get("house")
	if house.Delete(model.Contract{ID: in.Id}).Error != nil {
		common.ErrorResposne(c, "删除失败")
		return
	}
	common.SuccessResposne(c, gin.H{"message": "删除成功"})
}

type ContractIndexOut struct {
	Id        int    `json:"id"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Price     int    `json:"price"`
	Month     int    `json:"month"`
	HouseName string `json:"house_name"`
	CardName  string `json:"card_name"`
	CardNum   string `json:"card_num"`
}

type ContractIndexIn struct {
	HouseId int `json:"house_id,string"`
}

func ContractIndex(c *gin.Context) {
	in := ContractIndexIn{}
	if err := c.BindJSON(&in); err != nil {
		common.ErrorResposne(c, err.Error())
		return
	}

	house := common.DbConnections.Get("house")
	rows, err := house.Table("contract c").
		Joins("inner join house h on h.id=c.house_id and h.id=?", in.HouseId).
		Select("c.id,c.start_time,c.end_time,c.price,c.month,h.name,c.card_num,c.card_name").Rows()
	if err != nil {
		common.ErrorResposne(c, err.Error())
		return
	}
	var startTime, endTime time.Time
	var items []*ContractIndexOut
	for rows.Next() {
		item := ContractIndexOut{}
		rows.Scan(&item.Id, &startTime, &endTime, &item.Price, &item.Month, &item.HouseName, &item.CardNum, &item.CardName)
		item.StartTime = startTime.Format("2006-01-02")
		item.EndTime = endTime.Format("2006-01-02")
		items = append(items, &item)
	}
	common.SuccessResposne(c, gin.H{"data": items})
}
