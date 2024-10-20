package controllers

import (
	"ecommerce/database"
	"ecommerce/models"
	"ecommerce/tokenjwt"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CheckOutOrder(c *gin.Context) {
	claims, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authorized"})
		return
	}

	userClaims, ok := claims.(*tokenjwt.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token data"})
		return
	}

	userID := userClaims.UserID

	var cartItems []models.Cart
	if err := database.DB.Where("user_id = ?", userID).Find(&cartItems).Error; err != nil || len(cartItems) == 0 {
		log.Printf("Cart retrieval error or cart empty: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart is empty"})
		return
	}

	log.Printf("Processing checkout for user ID: %d with %d items in cart\n", userID, len(cartItems))

	var total float64
	for _, item := range cartItems {
		var product models.Product
		if err := database.DB.First(&product, item.ProductID).Error; err != nil {
			log.Printf("Product not found in cart: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Product not found in cart"})
			return
		}

		if item.Quantity > product.Stock {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Not enough stock for product %s. Available: %d, Requested: %d", product.Name, product.Stock, item.Quantity),
			})
			return
		}

		total += float64(item.Quantity) * product.Price
	}

	log.Printf("Total price calculated: %f\n", total)

	var newAddress models.Address
	if err := c.ShouldBindJSON(&newAddress); err != nil {
		log.Printf("Error binding address data: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address data", "details": err.Error()})
		return
	}

	PhoneNumber := len(fmt.Sprint(newAddress.Phone))
	if PhoneNumber > 10 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number is more than 10 digits"})
		return
	} else if PhoneNumber < 10 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number is less than 10 digits"})
		return
	}

	PostalCode := len(fmt.Sprint(newAddress.PostalCode))
	if PostalCode > 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Postal code is more than 6 digits"})
		return
	} else if PostalCode < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Postal code is less than 6 digits"})
		return
	}

	newAddress.UserID = userID
	if err := database.DB.Create(&newAddress).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save address", "details": err.Error()})
		return
	}

	order := models.Order{
		UserID:    userID,
		Total:     total,
		AddressID: newAddress.ID,
		CreatedAt: time.Now(),
	}

	if err := database.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not place order", "details": err.Error()})
		return
	}

	for _, item := range cartItems {
		var product models.Product
		if err := database.DB.First(&product, item.ProductID).Error; err != nil {
			log.Printf("Product not found: %v", err)
			continue
		}

		product.Stock -= item.Quantity
		if err := database.DB.Save(&product).Error; err != nil {
			log.Printf("failed to update stock for product %d: %v", product.ID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update stock"})
			return
		}
	}

	if err := database.DB.Where("user_id = ?", userID).Delete(&models.Cart{}).Error; err != nil {
		log.Printf("Failed to clear cart after checkout: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear cart after checkout"})
		return
	}

	log.Printf("Cart cleared successfully for user ID: %d\n", userID)

	c.JSON(http.StatusOK, gin.H{
		"message":     "Checkout successful",
		"total_price": total,
		"Order_id":    order.ID,
	})
}

func GetOrders(c *gin.Context) {
	userIDif, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not authenticated"})
		return
	}

	var order models.Order
	database.DB.Preload("Address").Preload("User").First(&order, order.ID)

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
