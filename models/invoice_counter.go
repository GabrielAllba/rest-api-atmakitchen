package models

type InvoiceCounter struct {
	ID          int `gorm:"primaryKey;" json:"id"`
	LastInvoice int
}
