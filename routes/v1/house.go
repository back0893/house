package v1

import (
	"database/sql"
	"main/common"
	"main/model"
	"time"

	"github.com/gin-gonic/gin"
)

type HouseIn struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
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

func HouseEdit(c *gin.Context) {
	edit := HouseIn{}
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
	ContractPrice sql.NullInt32
	Month         sql.NullInt32
	ContractId    sql.NullInt32
	EndTime       time.Time
}

type HouseListOut struct {
	Id               int    `json:"id"`
	Name             string `json:"name"`
	Price            int    `json:"price"`
	HasContract      bool   `json:"has_contract"`
	ContractId       int32  `json:"contract_id"`
	LeaveContractDay int32  `json:"leave_contract_day"`
	PayDay           int    `json:"pay_day"`
	Month            int32  `json:"-"`
	LastContractTime string `json:"last_contract_time"`
	Money            int32  `json:"money"`
}

//对于进入来说,应该显示名称,出租价格,是否出租,剩余租期,是否应该交租,每隔几月交租,上次交租时间,租金(交)
func HousenIndex(c *gin.Context) {
	house := common.DbConnections.Get("house")
	rows := make([]*HouseContract, 0)
	now := time.Now().In(common.GetLoc())
	if err := house.Table("house h").
		Joins("left join contract c on h.id=c.house_id").Where("(c.start_time<=? and c.end_time>=?) or (c.id is null)", now.Format("2006-01-02"), now.Format("2006-01-02")).
		Select("h.id,h.name,h.price,c.price as contract_price,c.month,c.id as `contract_id`,c.end_time").Scan(&rows).Error; err != nil {
		common.ErrorResposne(c, err.Error())
		return
	}

	var items []*HouseListOut

	var contractId []int

	for _, row := range rows {
		item := &HouseListOut{
			Id:    row.Id,
			Name:  row.Name,
			Price: row.Price,
		}
		if row.ContractId.Valid {
			contractId = append(contractId, int(row.ContractId.Int32))
			item.HasContract = true
			item.LeaveContractDay = int32(row.EndTime.Sub(now).Hours() / 24)
			item.Month = row.Month.Int32
			item.Money = row.ContractPrice.Int32 * item.Month
			item.ContractId = row.ContractId.Int32
		}

		items = append(items, item)
	}
	//获得最新交租日志
	logRowsId, err := house.Table("contract_log").Where("contract_id in (?)", contractId).
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
		Select("money,contract_at,contract_id").Scan(&logs).Error; err != nil {
		common.ErrorResposne(c, err.Error())
		return
	}
	logsMap := make(map[int32]*model.ContractLog)
	for _, log := range logs {
		logsMap[int32(log.ContractId)] = log
	}
	contracts := make([]*model.Contract, 0)

	house.Table("contract").Where("start_time<=? and end_time>=? and id in (?)", now.Format("2006-01-02"), now.Format("2006-01-02"), contractId).Find(&contracts)

	for _, item := range items {
		if log, ok := logsMap[item.ContractId]; ok {
			item.LastContractTime = log.ContractAt.Format("2006-01-02")
			item.PayDay = int(log.ContractAt.AddDate(0, int(item.Month), 0).Sub(now).Hours() / 24)
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
