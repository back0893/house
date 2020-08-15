package v1

import (
	"database/sql"
	"main/common"
	"main/model"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
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

type HouseListOut struct {
	Id               int    `json:"id"`
	Name             string `json:"name"`
	Price            int    `json:"price"`
	HasContract      bool   `json:"has_contract"`
	ContractId       int    `json:"contract_id"`
	LeaveContractDay int    `json:"leave_contract_day"`
	PayDay           int    `json:"pay_day"`
	Month            int    `json:"month"`
	LastContractTime string `json:"last_contract_time"`
	Money            int    `json:"money"`
}

//对于进入来说,应该显示名称,出租价格,是否出租,剩余租期,是否应该交租,每隔几月交租,上次交租时间,租金(交)
func HousenIndex(c *gin.Context) {
	house := common.DbConnections.Get("house")
	now := time.Now().In(common.GetLoc())

	var houses []*model.House
	if err := house.Find(&houses).Error; err != nil {
		common.ErrorResposne(c, err.Error())
		return
	}

	var items []*HouseListOut

	var houseID []int

	for _, row := range houses {
		houseID = append(houseID, row.ID)
		items = append(items, &HouseListOut{
			Id:    row.ID,
			Name:  row.Name,
			Price: row.Price,
		})
	}
	//获得合同
	var contracts []*model.Contract
	if err := house.Where("start_time<=? and end_time>=? and cancel=0", now.Format("2006-01-02"), now.Format("2006-01-02")).Find(&contracts).Error; err != nil {
		common.ErrorResposne(c, err.Error())
		return
	}

	var contractId []int
	for _, contract := range contracts {
		contractId = append(contractId, contract.ID)
		for _, item := range items {
			if item.Id == contract.HouseId {
				item.HasContract = true
				item.ContractId = contract.ID
				item.Money = contract.Month
				item.Money = contract.Month * contract.Price
				item.LeaveContractDay = int(contract.EndTime.Sub(now).Hours() / 24)
				break
			}
		}
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
	logsMap := make(map[int]*model.ContractLog)
	for _, log := range logs {
		logsMap[log.ContractId] = log
	}

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
		if err == gorm.ErrRecordNotFound {
			common.ErrorResposne(c, "查询失败")
			return
		}
		common.ErrorResposne(c, err.Error())
		return
	}
	common.SuccessResposne(c, house)
}

type HouseSimpleOut struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func HouseSimple(c *gin.Context) {
	house := common.DbConnections.Get("house")
	houses := make([]*model.House, 0)
	house.Select("id,name").Find(&houses)
	var items []*HouseSimpleOut
	for _, house := range houses {
		items = append(items, &HouseSimpleOut{
			ID:   house.ID,
			Name: house.Name,
		})
	}
	common.SuccessResposne(c, items)
}
