package model

import "time"

type Contract struct {
	Id        int
	StartTime time.Time
	EndTime   time.Time
	Price     int
	Month     int
	HouseId   int
}
