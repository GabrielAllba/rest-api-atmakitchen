package models

type Token struct {
	LogoutToken string `gorm:"type:varchar(255)" json:"logout_token"`
	Expiration  string `gorm:"type:varchar(255)" json:"expiration"`
	UserId      int    `gorm:"index" json:"user_id"`
	User        User   `gorm:"foreignKey:UserId" json:"user"`
}
