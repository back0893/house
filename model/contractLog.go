package model

import "time"

type ContractLog struct {
	ID         int
	HouseId    int
	HouseName  string
	Money      int
	ContractAt time.Time
	Remark     string
	ContractId int
}

func (cl *ContractLog) TableName() string {
	return "contract_log"
}
