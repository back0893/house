package model

type House struct {
	ID    int
	Name  string
	Price float64
}

func (h *House) TableName() string {
	return "house"
}
