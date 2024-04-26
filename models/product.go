package models

type Product struct {
	Id             int          `gorm:"primaryKey;" json:"id"`
	Name           string       `gorm:"type:varchar(255)" json:"name"`
	Price          float32      `gorm:"type:float;" json:"price"`
	Description    string       `gorm:"type:varchar(255)" json:"description"`
	Photo          string       `gorm:"type:varchar(255)" json:"photo"`
	Stock          float32      `gorm:"type:float;" json:"stock"`
	DailyQuota     float32      `gorm:"type:float;" json:"daily_quota"`
	RewardPoin     int          `gorm:"type:int" json:"reward_poin"`
	Status         string       `gorm:"type:varchar(255)" json:"status"`
	ProductTypeId  int          `gorm:"index" json:"product_type_id"`
	ProductType    ProductType  `gorm:"foreignKey:ProductTypeId" json:"product_type"`
	ConsignationId *int         `gorm:"index" json:"consignation_id"`
	Consignation   Consignation `gorm:"foreignKey:ConsignationId" json:"consignation"`
}
