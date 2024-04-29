package models

type DetailHampers struct {
	Id        int     `gorm:"primaryKey;" json:"id"`
	HampersId int     `gorm:"index" json:"hampers_id"`
	Hampers   Hampers `gorm:"foreignKey:HampersId" json:"hampers"`
	ProductId int     `gorm:"index" json:"product_id"`
	Product   Product `gorm:"foreignKey:ProductId" json:"product"`
	jumlah    int     `gorm:"type:int" json:"jumlah"`
}
