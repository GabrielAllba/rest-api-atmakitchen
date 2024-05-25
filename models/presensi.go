package models

type Presensi struct {
	Id      int    `gorm:"primaryKey;" json:"id"`
	Tanggal string `gorm:"type:varchar(255)" json:"tanggal"`
	Status  string `gorm:"type:varchar(255)" json:"status"`
	UserId  int    `gorm:"index" json:"user_id"`
	User    User   `gorm:"foreignKey:UserId" json:"user"`
}
