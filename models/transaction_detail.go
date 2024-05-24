package models

type TransactionDetail struct {
	ID                int      `gorm:"primaryKey;" json:"id"`
	InvoiceNumber     *string  `gorm:"type:varchar(255)" json:"invoice_number"`
	Jenis             *string  `gorm:"type:varchar(255)" json:"jenis"`
	ProductId         *int     `gorm:"index" json:"product_id"`
	Product           *Product `gorm:"foreignKey:ProductId" json:"product"`
	ProductQuantity   *float64 `gorm:"type:float" json:"product_quantity"`
	ProductPrice      *float64 `gorm:"type:float;" json:"product_price"`
	HampersId         *int     `gorm:"index" json:"hampers_id"`
	Hampers           *Hampers `gorm:"foreignKey:HampersId" json:"hampers"`
	HampersQuantity   *float64 `gorm:"type:float" json:"hampers_quantity"`
	HampersPrice      *float64 `gorm:"type:float;" json:"hampers_price"`
	TanggalPengiriman *string  `gorm:"type:varchar(255)" json:"tanggal_pengiriman"`
	TransactionStatus *string  `gorm:"type:varchar(255)" json:"transaction_status"`
	JenisItem         *string  `gorm:"type:varchar(255)" json:"jenis_item"`
}
