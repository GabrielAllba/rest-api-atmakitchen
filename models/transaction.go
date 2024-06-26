package models

type Transaction struct {
	ID                int     `gorm:"primaryKey;" json:"id"`
	InvoiceNumber     string  `gorm:"index" json:"invoice_number"`
	UserId            int     `gorm:"index" json:"user_id"`
	User              User    `gorm:"foreignKey:UserId" json:"user"`
	LunasPada         string  `gorm:"type:varchar(255)" json:"lunas_pada"`
	TanggalPemesanan  string  `gorm:"type:varchar(255)" json:"tanggal_pemesanan"`
	NamaPenerima      string  `gorm:"type:varchar(255)" json:"nama_penerima"`
	AlamatPenerima    string  `gorm:"type:varchar(255)" json:"alamat_penerima"`
	NoTelpPenerima    string  `gorm:"type:varchar(255)" json:"no_telp_penerima"`
	Delivery          string  `gorm:"type:varchar(255)" json:"delivery"`
	TransactionStatus string  `gorm:"type:varchar(255)" json:"transaction_status"`
	Distance          float64 `gorm:"type:float;" json:"distance"`
	DeliveryFee       float64 `gorm:"type:float;" json:"delivery_fee"`
	TotalPrice        float64 `gorm:"type:float;" json:"total_price"`
	TransferNominal   float64 `gorm:"type:float;" json:"transfer_nominal"`
	PointUsed         float64 `gorm:"type:float;" json:"point_used"`
	PointIncome       float64 `gorm:"type:float;" json:"point_income"`
	PaymentDate       string  `gorm:"type:varchar(255);" json:"payment_date"`
	PaymentProof      string  `gorm:"type:varchar(255)" json:"payment_proof"`
	TotalPoinUser     float64 `gorm:"type:float;" json:"total_poin_user"`
	UserTransfer      float64 `gorm:"type:float;" json:"user_transfer"`
	Tips              float64 `gorm:"type:float;" json:"tips"`
	// TransactionDetails []TransactionDetail `gorm:"foreignKey:InvoiceNumber;references:InvoiceNumber"`
}
