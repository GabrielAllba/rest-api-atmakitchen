package models

type Quota struct {
	Id        int      `gorm:"primaryKey;" json:"id"`
	ProductId *int     `gorm:"index" json:"product_id"`
	Product   *Product `gorm:"foreignKey:ProductId" json:"product"`
	HampersId *int     `gorm:"index" json:"hampers_id"`
	Hampers   *Hampers `gorm:"foreignKey:HampersId" json:"hampers"`
	Tanggal   string   `gorm:"type:varchar(255)" json:"tanggal"`
	Quota     float64  `gorm:"type:float" json:"quota"`
}
