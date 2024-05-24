package models

type Transaction struct {
	ID                int     `gorm:"primaryKey;" json:"id"`
	InvoiceNumber     string  `gorm:"type:varchar(255)" json:"invoice_number"`
	UserId            int     `gorm:"index" json:"user_id"`
	User              User    `gorm:"foreignKey:UserId" json:"user"`
	LunasPada         string  `gorm:"type:varchar(255)" json:"lunas_pada"`
	TanggalAmbil      string  `gorm:"type:varchar(255)" json:"tanggal_ambil"`
	NamaPenerima      string  `gorm:"type:varchar(255)" json:"nama_penerima"`
	AlamatPenerima    string  `gorm:"type:varchar(255)" json:"alamat_penerima"`
	Delivery          string  `gorm:"type:varchar(255)" json:"delivery"`
	TransactionStatus string  `gorm:"type:varchar(255)" json:"transaction_status"`
	Distance          float64 `gorm:"type:float;" json:"distance"`
	DeliveryFee       float64 `gorm:"type:float;" json:"delivery_fee"`
	TotalPrice        float64 `gorm:"type:float;" json:"total_price"`
	TransferNominal   float64 `gorm:"type:float;" json:"transfer_nominal"`
	PointUser         float64 `gorm:"type:float;" json:"point_user"`
	PointIncome       float64 `gorm:"type:float;" json:"point_income"`
	PaymentDate       string  `gorm:"type:varchar(255);" json:"payment_date"`
	PaymentProof      string  `gorm:"type:varchar(255)" json:"payment_proof"`
	PoinUser          float64 `gorm:"type:float;" json:"poin_user"`
	TotalPoinUser     float64 `gorm:"type:float;" json:"total_poin_user"`
}
