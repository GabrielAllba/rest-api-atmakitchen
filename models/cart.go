package models

type Cart struct {
	ID                 int     `gorm:"primaryKey" json:"id"`
	UserId             int     `gorm:"index" json:"user_id"`
	User               User    `gorm:"foreignKey:UserId" json:"user"`
	ProductId          *int    `gorm:"index" json:"product_id"`
	Product            Product `gorm:"foreignKey:ProductId" json:"product"`
	JenisItem          string  `gorm:"type:varchar(255)" json:"jenis_item"`
	HampersId          *int    `gorm:"index" json:"hampers_id"`
	Hampers            Hampers `gorm:"foreignKey:HampersId" json:"hampers"`
	Quantity           float64 `gorm:"type:float;" json:"quantity"`
	TotalPrice         float64 `gorm:"type:float;" json:"total_price"`
	Status             string  `gorm:"type:varchar(255);" json:"status"`
	Jenis              string  `gorm:"type:varchar(255);" json:"jenis"`
	OpsiPengambilan    string  `gorm:"type:varchar(255);" json:"opsi_pengambilan"`
	TanggalPengiriman  *string `json:"tanggal_pengiriman"`
	TanggalPengambilan *string `json:"tanggal_pengambilan"`
}