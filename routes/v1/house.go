package v1

import (
	"database/sql"
	"main/common"
	"main/model"
	"time"

	"github.com/gin-gonic/gin"
)

type HouseIn struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}
type HouseDeleteIn struct {
	ID int `json:"id"`
}

func HouseAdd(c *gin.Context) {
	in := HouseIn{}
	if err := c.BindJSON(&in); err != nil {
		common.ErrorResposne(c, err.Error())
		return
	}
	houseModel := model.House{}
	houseModel.Name = in.Name
	houseModel.Price = in.Price
	house := common.DbConnections.Get("house")
	if err := house.Create(&houseModel).Error; err != nil {
		common.ErrorResposne(c, "新增失败")
		return
	}
	common.SuccessResposne(c, gin.H{"message": "新增成功"})
}

type HouseEditIn struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func HouseEdit(c *gin.Context) {
	edit := HouseEditIn{}
	if err := c.BindJSON(&edit); err != nil {
		common.ErrorResposne(c, err.Error())
		return
	}
	houseModel := model.House{}
	houseModel.Name = edit.Name
	houseModel.Price = edit.Price
	house := common.DbConnections.Get("house")
	if err := house.Model(&houseModel).Where("id=?", edit.ID).Updates(houseModel).Error; err != nil {
		common.ErrorResposne(c, err.Error())
		return
	}
	common.SuccessResposne(c, gin.H{"message": "修改成功"})
}

func HouseDelete(c *gin.Context) {
	hd := HouseDeleteIn{}
	if err := c.BindJSON(&hd); err != nil {
		common.ErrorResposne(c, err.Error())
		return
	}
	house := common.DbConnections.Get("house")
	if err := house.Delete(&model.House{ID: hd.ID}).Error; err != nil {
		common.ErrorResposne(c, err.Error())
		return
	}
	common.SuccessResposne(c, gin.H{"message": "删除成功"})
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
	Id               int    `json:"id"`
	Name             string `json:"name"`
	Price            int    `json:"price"`
	ContractPrice    int    `json:""contract_price`
	CardName         string `json:"card_name"`
	Month            int    `json:"month"`
	LastContractTime string `json:"last_contract_time"`
	Money            int    `json:"money"`
	Pay              bool   `json:"pay"`
}

//对于进入来说,应该显示名称,出租价格,当前住户,交租价格,上次交租日期,租金,是否应该交租
func HousenIndex(c *gin.Context) {
	house := common.DbConnections.Get("house")
	rows := make([]*HouseContract, 0)
	now := time.Now().In(common.GetLoc())
	if err := house.Table("house h").
		Joins("left join contract c on h.id=c.house_id").Where("start_time<=? and end_time>=?", now.Format("2006-01-02"), now.Format("2006-01-02")).
		Select("h.id,h.name,h.price,c.card_name,c.price as contract_price,c.month").Scan(&rows).Error; err != nil {
		common.ErrorResposne(c, err.Error())
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
		common.ErrorResposne(c, err.Error())
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
		Select("money,contract_at,house_id").Scan(&logs).Error; err != nil {
		common.ErrorResposne(c, err.Error())
		return
	}
	logsMap := make(map[int]*model.ContractLog)
	for _, log := range logs {
		logsMap[log.HouseId] = log
	}
	contracts := make([]*model.Contract, 0)

	house.Table("contract").Where("start_time<=? and end_time>=? and house_id in (?)", now.Format("2006-01-02"), now.Format("2006-01-02"), houseId).Find(&contracts)

	for _, item := range items {
		if log, ok := logsMap[item.Id]; ok {
			if common.FindIn(log.HouseId, houseId) {
				item.LastContractTime = log.ContractAt.Format("2006-01-02")
				item.Money = log.Money
				if log.ContractAt.AddDate(0, item.Month, 0).Before(now) {
					item.Pay = true
				}
			}
		}
	}
	common.SuccessResposne(c, gin.H{"data": items})
}

func HouseRow(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		common.ErrorResposne(c, "id不能为空")
		return
	}
	house := model.House{}
	ep := common.DbConnections.Get("ep")
	if err := ep.Where("id=?", id).First(&house).Error; err != nil {
		if err == sql.ErrNoRows {
			common.ErrorResposne(c, "查询失败")
			return
		}
	}
	common.SuccessResposne(c, house)
}
