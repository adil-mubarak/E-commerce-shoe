package controllers

import (
	"ecommerce/database"
	"ecommerce/models"
	"ecommerce/tokenjwt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CheckOutOrder(c *gin.Context) {
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var cartItems []models.Cart
	 if err := database.DB.Where("user_id = ?", order.UserID).Find(&cartItems).Error; err != nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":"Could not retrieve cart items"})
		return
	 }

	var total float64
	for _, item := range cartItems {
		var product models.Product
		if err := database.DB.First(&product, item.ProductID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found in cart"})
			return
		}
		total += float64(item.Quantity) * product.Price
	}

	order.Total = total
	order.CreatedAt = time.Now()

	if err := database.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not place order"})
		return
	}

	if err := database.DB.Where("user_id", order.UserID).Delete(&models.Cart{}).Error; err != nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":"Could not clear cart"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Order placed successfully", "order": order})
}

func GetOrders(c *gin.Context) {
	userIDif, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not authenticated"})
		return
	}

	claims, ok := userIDif.(*tokenjwt.Claims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid token claims"})
		return
	}

	userIDUint := claims.UserID

	var orders []models.Order
	if err := database.DB.Where("user_id = ?", userIDUint).Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func GetAllOrder(c *gin.Context) {
	var orders []models.Order
	database.DB.Find(&orders)
	c.JSON(http.StatusOK, orders)
}
