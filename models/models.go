package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Email     string    `gorm:"unique;type:varchar(100)" json:"email"`
	Password  string    `gorm:"type:varchar(100)" json:"password"`
	Role      string    `gorm:"type:varchar(20)" json:"role"`
	Banned    bool      `json:"banned" gorm:"default:false"`
	Addresses []Address `json:"addresses,omitempty"`
}

type Product struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	Name        string  `gorm:"type:varchar(100)" json:"name"`
	Description string  `gorm:"type:varchar(255)" json:"description"`
	Price       float64 `json:"price"`
	Stock       uint    `json:"stock"`
	Category    string  `gorm:"type:varchar(100)" json:"category"`
}

type Cart struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	UserID    uint    `json:"user_id" gorm:"not null"`
	ProductID uint    `json:"product_id" gorm:"not null"`
	Quantity  uint    `json:"quantity"`
	User      User    `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`
	Product   Product `gorm:"foreignKey:ProductID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"product,omitempty"`
}

type Whishlist struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	UserID    uint    `json:"user_id" gorm:"not null"`
	ProductID uint    `json:"product_id" gorm:"not null"`
	User      User    `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`
	Product   Product `gorm:"foreignKey:ProductID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"product,omitempty"`
}

type Address struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	UserID     uint   `json:"user_id" gorm:"not null"`
	Street     string `gorm:"type:varchar(100)" json:"street"`
	City       string `gorm:"type:varchar(50)" json:"city"`
	State      string `gorm:"type:varchar(50)" json:"state"`
	PostalCode string `gorm:"type:varchar(20)" json:"postal_code"`
	Country    string `gorm:"type:varchar(50)" json:"country"`
	User       User   `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`
}

type Order struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `json:"user_id" gorm:"not null"`
	Total         float64   `json:"total"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	AddressID     uint      `json:"address_id" gorm:"not null"`
	Address       Address   `gorm:"foreignKey:AddressID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"address,omitempty"`
	User          User      `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`
	PaymentStatus string    `json:"payment_status" gorm:"type:varchar(20);default:'Pending'"`
	Status        string    `json:"status" gorm:"type:varchar(20);default:'Pending'"`
}

type OrderedItems struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	OrderID   uint      `json:"order_id" gorm:"not null"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	Total     float64   `json:"total"`
	AddressID uint      `json:"address_id"`
	Address   Address   `gorm:"foreignKey:AddressID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"address,omitempty"`
	Status    string    `json:"status" gorm:"type:varchar(20);default:'Pending'"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

type Payment struct {
	ID             uint      `gorm:"primarykey" json:"id"`
	OrderID        uint      `json:"order_id" gorm:"not null"`
	UserID         uint      `json:"user_id" gorm:"not null"`
	Amount         float64   `json:"amount"`
	Status         string    `json:"status" gorm:"varchar(20);default:'Pending'"`
	PaymentMethod  string    `json:"payment_method"`
	PaymentID      string    `json:"payment_id"`
	StripeChargeID string    `json:"stripe_charge_id"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
}
