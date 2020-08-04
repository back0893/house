package model

import "time"

type Contract struct {
	ID        int
	StartTime time.Time
	EndTime   time.Time
	Price     int
	Month     int
	HouseId   int
	CardName  string
	CardNum   string
	House     *House `gorm:"ForeignKey:HouseId"`
}

func (c *Contract) TableName() string {
	return "contract"
}
