package controllers

import (
	"ecommerce/database"
	"ecommerce/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CheckOutOrder(c *gin.Context) {
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		return
	}

	var cartItems []models.Cart
	database.DB.Where("user_id = ?",order.UserID).Find(&cartItems)
	

	var total float64
	for _,item := range cartItems{
		var product models.Product
		if err := database.DB.First(&product,item.ProductID).Error; err != nil{
			c.JSON(http.StatusNotFound,gin.H{"error":"Product not found in cart"})
			return
		}
		total += float64(item.Quantity) * product.Price
	}

	order.Total = total
	order.CreatedAt = time.Now()

	if err := database.DB.Create(&order).Error; err != nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":"Could not place order"})
		return
	}

	database.DB.Where("user_id",order.UserID).Delete(&models.Cart{})
	c.JSON(http.StatusOK,gin.H{"message":"Order placed successfully","order":order})
}

func GetOrders(c *gin.Context){
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusBadRequest,gin.H{"error":"User not authenticated"})
		return
	}

	var orders []models.Order
	database.DB.Where("user_id = ?",userID).Find(&orders)
	c.JSON(http.StatusOK,orders)
}

func GetAllOrder(c *gin.Context){
	var orders []models.Order
	database.DB.Find(&orders)
	c.JSON(http.StatusOK,orders)
}
