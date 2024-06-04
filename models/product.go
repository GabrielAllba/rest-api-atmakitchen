package models

import "mime/multipart"

type Product struct {
	Id             int          `gorm:"primaryKey;" json:"id"`
	Name           string       `gorm:"type:varchar(255)" json:"name"`
	Price          float64      `gorm:"type:float;" json:"price"`
	Description    string       `gorm:"type:varchar(255)" json:"description"`
	Photo          string      `gorm:"type:varchar(255)" json:"photo"`
	Stock          float64      `gorm:"type:float;" json:"stock"`
	DailyQuota     float64      `gorm:"type:float;" json:"daily_quota"`
	Status         string       `gorm:"type:varchar(255)" json:"status"`
	Tag            string      `gorm:"type:varchar(255)" json:"tag"`
	ProductTypeId  int          `gorm:"index" json:"product_type_id"`
	ProductType    ProductType  `gorm:"foreignKey:ProductTypeId" json:"product_type"`
	ConsignationId *int         `gorm:"index" json:"consignation_id"`
	Consignation   Consignation `gorm:"foreignKey:ConsignationId" json:"consignation"`
	BahanReseps []BahanResep `gorm:"foreignKey:ProductId" json:"bahan_reseps"`

}

type FileUpload struct {
	FileHeader *multipart.FileHeader `form:"file" binding:"required"`
}