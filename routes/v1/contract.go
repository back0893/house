package v1

import (
	"errors"
	"main/common"
	"main/model"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type ContracIn struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Price     int    `json:"price"`
	Month     int    `json:"month"`
	HouseId   int    `json:"house_id"`
	Id        int    `json:"id"`
	CardName  string `json:"card_name"`
	CardNum   string `json:"card_num"`
}

//检查合同时间是否交叉
func CheckContractCocos(startTime, endTime time.Time, houseId int) error {
	//新增时间判断,合同时间不能交叉
	//3个 被包含, 包含,交叉
	start, end := startTime.Format("2006-01-02 15:04:05"), endTime.Format("2006-01-02 15:04:05")
	house := common.DbConnections.Get("house")
	cocos := 0
	if err := house.Model(&model.Contract{HouseId: houseId}).Where("cancel=1").
		Where("(start_time >=? and end_time<=?) or  (start_time<=? and end_time>=?) or (start_time<=? and end_time>=?)", start, end, start, start, end, end).Count(&cocos).Error; err != nil {
		return err
	}
	if cocos > 0 {
		return errors.New("合同时间交叉")
	}
	return nil
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
	if err := CheckContractCocos(contractModel.StartTime, contractModel.EndTime, in.HouseId); err != nil {
		common.ErrorResposne(c, err.Error())
		return
	}
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
	if err := CheckContractCocos(contractModel.StartTime, contractModel.EndTime, in.HouseId); err != nil {
		common.ErrorResposne(c, err.Error())
		return
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

func ContractCancle(c *gin.Context) {
	in := ContractDeleteIn{}
	if err := c.BindJSON(&in); err != nil {
		common.ErrorResposne(c, err.Error())
		return
	}
	house := common.DbConnections.Get("house")
	if house.Model(&model.Contract{ID: in.Id}).Update("cancel", 1).Error != nil {
		common.ErrorResposne(c, "取消失败")
		return
	}
	common.SuccessResposne(c, gin.H{"message": "取消成功"})
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
	HouseId   int    `json:"house_id"`
	Cancel    bool   `json:"cancel"`
}

func ContractIndex(c *gin.Context) {
	houseId := c.Param("houseId")
	house := common.DbConnections.Get("house")
	query := house.Table("contract c").
		Joins("inner join house h on h.id=c.house_id").
		Select("c.id,c.start_time,c.end_time,c.price,c.month,h.name,c.card_num,c.card_name,c.house_id,c.cancel")
	if houseId != "" {
		query.Where("h.id=?", houseId)
	}
	rows, err := query.Rows()
	if err != nil {
		common.ErrorResposne(c, err.Error())
		return
	}
	var startTime, endTime time.Time
	var items []*ContractIndexOut
	var cancel int
	for rows.Next() {
		item := ContractIndexOut{}
		rows.Scan(&item.Id, &startTime, &endTime, &item.Price, &item.Month, &item.HouseName, &item.CardNum, &item.CardName, &item.HouseId, &cancel)
		item.StartTime = startTime.Format("2006-01-02")
		item.EndTime = endTime.Format("2006-01-02")
		item.Cancel = cancel == 1
		items = append(items, &item)

	}
	common.SuccessResposne(c, gin.H{"data": items})
}

type FormItem struct {
	Field string      //显示
	Value interface{} //值
	Type  string      //类型
}

//但前有效的合同
func ContractValid(c *gin.Context) {
	hosueId := c.Query("houseId")
	now := common.Now().Format("2006-01-02")
	house := common.DbConnections.Get("house")
	contract := model.Contract{}
	if err := house.Where("start_time<=? and end_time>=? and cancel=0 and house_id=?", now, now, hosueId).First(&contract).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			common.ErrorResposne(c, "无法查询")
			return
		}
		common.ErrorResposne(c, err.Error())
		return
	}

	var items []*FormItem
	items = append(items, &FormItem{
		Field: "开始时间",
		Value: contract.StartTime,
		Type:  "date",
	})
	items = append(items, &FormItem{
		Field: "结束时间",
		Value: contract.EndTime,
		Type:  "date",
	})
	items = append(items, &FormItem{
		Field: "合同房租",
		Value: contract.Price,
		Type:  "number",
	})
	items = append(items, &FormItem{
		Field: "交租间隔",
		Value: contract.Month,
		Type:  "number",
	})
	items = append(items, &FormItem{
		Field: "租房人",
		Value: contract.CardName,
		Type:  "string",
	})
	items = append(items, &FormItem{
		Field: "合同人身份证",
		Value: contract.CardNum,
		Type:  "string",
	})

	common.SuccessResposne(c, contract)
}
