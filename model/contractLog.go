package model

import "time"

type ContractLog struct {
	ID         int
	Money      int
	ContractAt time.Time
	Remark     string
	ContractId int
	ExtraMoney int
}

func (p *ContractLog) TableName() string {
	return "contract_log"
}

type ContractLogShow struct {
	ID         int
	Money      int
	ContractAt time.Time
	Remark     string
	ContractId int
	ExtraMoney int
	HouseName  string
	HouseId    int
}
