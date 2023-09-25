package models

type Postback struct {
	ID             int    `db:"id" json:"-"`
	MerchantID     int    `db:"merchant_id" json:"-"`
	PostbackUrl    string `db:"postback_url" json:"postback_url" binding:"required"`
	PostbackMethod string `db:"postback_method" json:"postback_method" binding:"required"`
	IsEnabled      bool   `db:"is_enabled" json:"is_enabled"`
}
