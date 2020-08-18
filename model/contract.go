package model

import "time"

type Contract struct {
	ID        int       `json:"id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Price     int       `json:"price"`
	Month     int       `json:"month"`
	HouseId   int       `json:"house_id"`
	CardName  string    `json:"card_name"`
	CardNum   string    `json:"card_num"`
	Cancel    int       `json:"cancel"`
	Phone     string    `json:"phone"`
}

func (c *Contract) TableName() string {
	return "contract"
}
