package models

type PembelianBahanBaku struct {
	Id               int     `gorm:"primaryKey;" json:"id"`
	BahanId          int     `gorm:"index" json:"bahan_id"`
	Bahan            Bahan   `gorm:"foreignKey:BahanId" json:"bahan"`
	Jumlah           float64 `gorm:"type:float;" json:"jumlah"`
	Keterangan       string  `gorm:"type:string;" json:"keterangan"`
	TanggalPembelian string  `gorm:"type:varchar(255)" json:"tanggal_pembelian"`
}
