package models

type User struct {
	Id          int    `gorm:"primaryKey;" json:"id"`
	Name        string `gorm:"type:varchar(255)" json:"name"`
	Email       string `gorm:"type:varchar(255)" json:"email"`
	Username    string `gorm:"type:varchar(255)" json:"username"`
	Password    string `gorm:"type:varchar(255)" json:"password"`
	BornDate    string `gorm:"type:varchar(255)" json:"born_date"`
	PhoneNumber string `gorm:"type:varchar(255)" json:"phone_number"`
	TotalPoint  int    `json:"total_point"`
	Balance     float64 `json:"balance"` // Tambahkan atribut saldo
	RoleId      int    `gorm:"index" json:"role_id"`
	Role        Role   `gorm:"foreignKey:RoleId" json:"role"`
}
