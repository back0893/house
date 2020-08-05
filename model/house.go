package model

type House struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func (h *House) TableName() string {
	return "house"
}
