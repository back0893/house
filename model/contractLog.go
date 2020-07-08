package model

import "time"

type ContractLog struct{
	Id int
	HouseId int
	HouseName string
	Monoey int
	ContractAt time.Time
	Remark string
}