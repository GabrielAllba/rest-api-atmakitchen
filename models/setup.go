package models

import (
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := os.Getenv("DB")
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	database.AutoMigrate(&User{}, &Role{}, &ProductType{}, &Consignation{}, &Product{}, &Bank{}, &Hampers{}, &DetailHampers{}, &Resep{}, &Bahan{}, &Token{}, &BahanProduct{}, &BahanResep{}, &PembelianBahanBaku{}, &Cart{}, &Transaction{}, &TransactionDetail{}, &InvoiceCounter{}, &Quota{}, &PengeluaranLain{}, &Presensi{})
	DB = database
}
