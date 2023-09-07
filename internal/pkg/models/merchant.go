package models

type Merchant struct {
	ID     int    `db:"id" json:"id,omitempty"`
	Title  string `db:"title" json:"title,omitempty" binding:"required"`
	ApiKey string `db:"api_key" json:"api_key,omitempty"`
}
