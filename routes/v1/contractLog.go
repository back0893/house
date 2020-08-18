package v1

import (
	"main/common"
	"main/model"
	"time"

	"github.com/gin-gonic/gin"
)

type ContractLogIn struct {
	ContractId int    `json:"contract_id"`
	ContractAt string `json:"contract_at"`
	Money      int    `json:"money"`
	Remark     string `json:"content"`
	ExtraMoney int    `json:"extra_money"`
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
		ExtraMoney: in.ExtraMoney,
	}
	var err error
	log.ContractAt, err = common.ParseTime("2006-01-02", in.ContractAt)
	if err != nil {
		common.ErrorResposne(c, "时间传递错误")
		return
	}
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
	ContractId int       `form:"contract_id"`
	HouseId    int       `form:"houseid"`
	StartTime  time.Time `form:"start_time"`
	EndTime    time.Time `form:"end_time"`
}

type ContractLogIndexOut struct {
	Items []*ContractLogIndexOutItem `json:"item"`
}
type ContractLogIndexOutItem struct {
	ID         int    `json:"id"`
	Money      int    `json:"money"`
	ContractAt string `json:"created_at"`
	Remark     string `json:"remark"`
	ExtraMoney int    `json:"extra_money"`
	HouseName  string `json:"house_name"`
	HouseId    int    `json:"house_id"`
	ContractId int    `json:"contract_id"`
}

func getContractId(houseId, contractId int) []int {
	house := common.DbConnections.Get("house")
	query := house.Model(&model.Contract{}).Where("house_id=?", houseId)
	if contractId != 0 {
		query.Where("id=?", contractId)
	}
	var IDs []int
	query.Pluck("id", &IDs)
	return IDs

}
func ContractLogIndex(c *gin.Context) {
	in := ContractLogIndexIn{}
	if err := c.BindQuery(&in); err != nil {
		common.ErrorResposne(c, err.Error())
		return
	}
	query := common.DbConnections.Get("house").Table("contract_log log").
		Joins("inner join contract c on c.id=log.contract_id").Joins("inner join house h on h.id=c.house_id")
	logs := make([]*model.ContractLogShow, 0)
	if !in.StartTime.IsZero() {
		query.Where("log.contract_at>=?", in.StartTime.Format("2006-01-02"))
	}
	if !in.EndTime.IsZero() {
		query.Where("log.contract_at<=?", in.EndTime.Format("2006-01-02"))
	}
	if in.HouseId > 0 {
		query.Where("log.contract_id in (?)", getContractId(in.HouseId, in.ContractId))
	}

	if err := query.Select("log.*,h.name as house_name,h.id as house_id").
		Find(&logs).Error; err != nil {
		common.ErrorResposne(c, err.Error())
		return
	}
	out := ContractLogIndexOut{}
	for _, log := range logs {
		item := &ContractLogIndexOutItem{
			ID:         log.ID,
			Money:      log.Money,
			ContractAt: log.ContractAt.Format("2006-01-02"),
			Remark:     log.Remark,
			ExtraMoney: log.ExtraMoney,
			HouseId:    log.HouseId,
			HouseName:  log.HouseName,
			ContractId: log.ContractId,
		}
		out.Items = append(out.Items, item)
	}
	common.SuccessResposne(c, out)
}

type StatiscHouseIn struct {
	HouseId   int       `json:"hosue_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

type StatiscHouseItem struct {
	ContractAt time.Time `json:"contract_at"`
	Money      int       `json:"money"`
	ExtraMoney int       `json:"extra_money"`
	HouseId    int       `json:"house_id"`
	HouseName  string    `json:"house_name" gorm:"column:name"`
}

type StatiscHouseOutItem struct {
	HouseId   int                 `json:"house_id"`
	HouseName string              `json:"house_name"`
	Items     []*StatiscHouseItem `josn:"items"`
}

//统计一段时间内房间的收租情况
func StatiscHouse(c *gin.Context) {
	in := StatiscHouseIn{}
	if err := c.BindJSON(&in); err != nil {
		common.ErrorResposne(c, err.Error())
		return
	}
	var items []*StatiscHouseItem
	house := common.DbConnections.Get("house")
	if err := house.Table("contract_log log").Joins("contract as c on c.id=log.contract_id").
		Joins("house h on h.id=c.house_id").
		Where("c.house_id=? and log.contract_at>=? and log.contract_at<=?", in.HouseId, in.StartTime.Format("2006-01-02"), in.EndTime.Format("2006-01-02")).Select("log.contract_at,log.money,log.extra_money,h.id,h.name").Find(&items).Error; err != nil {
		common.ErrorResposne(c, err.Error())
		return
	}
	tmp := make(map[int][]*StatiscHouseItem)
	for _, row := range items {
		if _, ok := tmp[row.HouseId]; !ok {
			tmp[row.HouseId] = make([]*StatiscHouseItem, 0)
		}
		tmp[row.HouseId] = append(tmp[row.HouseId], row)
	}

	out := make([]*StatiscHouseOutItem, 0)
	for key, value := range tmp {
		out = append(out, &StatiscHouseOutItem{
			HouseId:   key,
			HouseName: value[0].HouseName,
			Items:     value,
		})
	}
	common.SuccessResposne(c, out)
}
