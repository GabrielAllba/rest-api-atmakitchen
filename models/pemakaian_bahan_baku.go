package models

type PemakaianBahanBaku struct {
	Id                  int               `gorm:"primaryKey;" json:"id"`
	BahanId             int               `gorm:"index" json:"bahan_id"`
	Bahan               Bahan             `gorm:"foreignKey:BahanId" json:"bahan"`
	TransactionDetailId int               `gorm:"index" json:"transaction_detail_id"`
	TransactionDetail   TransactionDetail `gorm:"foreignKey:TransactionDetailId" json:"transaction_detail"`
	Jumlah              float64           `gorm:"type:float" json:"jumlah"`
	Tanggal             string            `gorm:"type:string" json:"tanggal"`
}
