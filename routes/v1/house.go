package v1

import (
	"database/sql"
	"main/common"
	"main/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HouseIn struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	ID    int     `json:"id"`
}
type HouseDeleteIn struct {
	ID int `json:"id"`
}

func HouseAdd(c *gin.Context) {
	ha := HouseIn{}
	if err := c.BindJSON(&ha); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	houseModel := model.House{}
	houseModel.Name = ha.Name
	houseModel.Price = ha.Price
	house := common.DbConnections.Get("house")
	if ok := house.NewRecord(houseModel); !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "新增失败"})
		return
	}
	c.JSON(http.StatusOK, common.SuccessResposne("新增成功", nil))
}

func HouseEdit(c *gin.Context) {
	edit := HouseIn{}
	if err := c.BindJSON(&edit); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResposne(err.Error()))
		return
	}
	houseModel := model.House{}
	houseModel.Name = edit.Name
	houseModel.Price = edit.Price
	house := common.DbConnections.Get("house")
	if err := house.Model(&houseModel).Where("id=?", edit.ID).Updates(houseModel).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, common.SuccessResposne("修改成功", nil))
}

func HouseDelete(c *gin.Context) {
	hd := HouseDeleteIn{}
	if err := c.BindJSON(&hd); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResposne(err.Error()))
		return
	}
	house := common.DbConnections.Get("house")
	if err := house.Delete(&model.House{Id: hd.ID}).Error; err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResposne(err.Error()))
		return
	}
	c.JSON(http.StatusOK, common.SuccessResposne("删除成功", nil))
}

type HouseContract struct {
	Id            int
	Name          string
	Price         int
	ContractPrice int
	CardName      string
	Month         int
}

type HouseListOut struct {
	Id               int
	Name             string
	Price            int
	ContractPrice    int
	CardName         string
	Month            int
	LastContractTime string
	Money            int
	Pay              bool
}

//对于进入来说,应该显示名称,出租价格,当前住户,交租价格,上次交租日期,租金,是否应该交租
func HousenIndex(c *gin.Context) {
	house := common.DbConnections.Get("house")
	rows := make([]*HouseContract, 0)
	if err := house.Table("house h").Joins("left join contract c on h.id=c.house_id").Select("h.id,h.name,h.price,c.card_name,c.price as contract_price,c.month").Scan(&rows).Error; err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResposne(err.Error()))
		return
	}

	var items []*HouseListOut

	var houseId []int
	for _, row := range rows {
		houseId = append(houseId, row.Id)
		items = append(items, &HouseListOut{
			Id:            row.Id,
			Name:          row.Name,
			Price:         row.Price,
			ContractPrice: row.ContractPrice,
			CardName:      row.CardName,
			Month:         row.Month,
		})
	}
	//获得最新交租日志
	logRowsId, err := house.Table("contract_log").Where("house_id in (?)", houseId).
		Select("max(id)").Group("id").Rows()
	if err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResposne(err.Error()))
		return
	}
	var ids []int32
	var id sql.NullInt32
	for logRowsId.Next() {
		logRowsId.Scan(&id)
		if id.Valid {
			ids = append(ids, id.Int32)
		}
	}
	logs := make([]*model.ContractLog, 0)
	if err := house.Table("contract_log").Where("id in (?)", ids).
		Select("money,contract_at,house_id").Scan(logs).Error; err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResposne(err.Error()))
		return
	}
	logsMap := make(map[int]*model.ContractLog)
	for _, log := range logs {
		logsMap[log.HouseId] = log
	}
	now := time.Now()
	for _, item := range items {
		if log, ok := logsMap[item.Id]; ok {
			item.LastContractTime = log.ContractAt.Format("2006-01-02")
			item.Money = log.Monoey
			if log.ContractAt.AddDate(0, item.Month, 0).Before(now) {
				item.Pay = true
			}
		}
	}
	c.JSON(http.StatusBadRequest, common.SuccessResposne("", items))
}
