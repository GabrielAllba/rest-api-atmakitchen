package models

type Role struct {
	Id         int     `gorm:"primaryKey;" json:"id"`
	Name       string  `gorm:"type:varchar(255)" json:"name"`
	GajiHarian float32 `gorm:"type:float;" json:"gaji_harian"`
}
