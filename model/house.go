package model

type House struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

func (h *House) TableName() string {
	return "house"
}
