package model

type Contract struct {
	Base
	Version string `json:"version"`
	Content string `json:"content"`
}

func (Contract) TableName() string {
	return "baas_contract"
}
