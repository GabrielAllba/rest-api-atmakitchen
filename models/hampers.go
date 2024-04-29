package models

type Hampers struct {
	Id          int     `gorm:"primaryKey;" json:"id"`
	HampersName string  `gorm:"type:varchar(255)" json:"hampers_name"`
	Price       float64 `gorm:"type:float;" json:"price"`
	Deskripsi   string  `gorm:"type:varchar(255)" json:"deskripsi"`
	Stock       float64 `gorm:"type:float;" json:"stock"`
	DailyQuota  float64 `gorm:"type:float;" json:"daily_quota"`
	Photo       string  `gorm:"type:varchar(255)" json:"photo"`
}
