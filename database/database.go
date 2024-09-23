package database

import (
	"ecommerce/models"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := "root:kl18jda183079@tcp(127.0.0.1:3306)/ecommerse?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil{
		log.Fatal("Failed to connect to database",err)
	}

	db.AutoMigrate(&models.User{},&models.Product{},&models.Cart{},&models.Order{},&models.Whishlist{})
	DB = db
}
