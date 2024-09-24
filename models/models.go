package models

import "time"

type User struct {
	ID       uint   `gorm:"primarykey"`
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique;type:varchar(100)"`
	Password string `json:"password" gorm:"type:varchar(100)"`
	Role     string `json:"role" gorm:"type:varchar(20)"`
}

type Product struct {
	ID          uint    `gorm:"primarykey"`
	Name        string  `json:"name" gorm:"type:varchar(100)"`
	Description string  `json:"description" gorm:"type:varchar(255)"`
	Price       float64 `json:"price"`
	Stock       uint    `json:"stock"`
	Category    string  `json:"catogery" gorm:"type:varchar(100)"`
}

type Cart struct {
	ID        uint `gorm:"primarykey"`
	UserID    uint `json:"user_id" gorm:"not null"`
	ProductID uint `json:"product_id" gorm:"not null"`
	Quantity  uint `json:"quantity"`
	User		User	`gorm:"foreignkey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Product		Product		`gorm:"foreignkey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Whishlist struct {
	ID        uint `gorm:"primarykey"`
	UserID    uint `json:"user_id" gorm:"not null"`
	ProductID uint `json:"product_id" gorm:"not null"`
	User	User	`gorm:"foreignkey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Product Product	`gorm:"foreignkey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Address struct {
	ID         uint   `gorm:"primarykey"`
	Street     string `json:"street" gorm:"type:varchar(100)"`
	City       string `json:"city"  gorm:"type:varchar(50)"`
	State      string `json:"state" gorm:"type:varchar(50)"`
	PostalCode string `json:"postal_code" gorm:"type:varchar(20)"`
	Country    string `json:"country"  gorm:"type:varchar(50)"`
}

type Order struct {
	ID        uint      `gorm:"primarykey"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	Total     float64   `json:"total"`
	CreatedAt time.Time `json:"created_at"`
	AddressID uint      `json:"address_id" gorm:"not null"`
	Address   Address   `gorm:"foreignkey:AddressID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	User      User      `gorm:"foreignkey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
