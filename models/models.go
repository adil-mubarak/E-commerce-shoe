package models

import "time"

type User struct {
	ID       uint   `gorm:"primarykey"`
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type Product struct {
	ID          uint    `gorm:"primarykey"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       uint    `json:"stock"`
	Category    string  `json:"catogery"`
}

type Cart struct {
	ID        uint `gorm:"primarykey"`
	UserID    uint `json:"user_id"`
	ProductID uint `json:"product_id"`
	Quantity  uint `json:"quantity"`
}

type Whishlist struct {
	ID        uint `gorm:"primarykey"`
	UserID    uint `json:"user_id"`
	ProductID uint `json:"product_id"`
}

type Order struct {
	ID        uint    `gorm:"primarykey"`
	UserID    uint    `json:"user_id"`
	Total     float64 `json:"total"`
	CreatedAt time.Time `json:"created_at"`
}

