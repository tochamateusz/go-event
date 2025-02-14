package entities

type Shown struct {
	ShowId string `json:"show_id" db:"show_id"`
	Amount uint16 `json:"amount" db:"amount"`
}
