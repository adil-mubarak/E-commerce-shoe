package controllers

import (
	"ecommerce/database"
	"ecommerce/models"
	"ecommerce/tokenjwt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddToCart(c *gin.Context) {
	claims, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not autherzide"})
		return
	}

	userClaims, ok := claims.(*tokenjwt.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token data"})
		return
	}

	userID := userClaims.UserID

	var cart models.Cart
	if err := c.ShouldBindJSON(&cart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cart.UserID = userID

	var product models.Product
	if err := database.DB.First(&product, cart.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if cart.Quantity > product.Stock {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Requasted quantity exceeds available stock"})
		return
	}

	var existingCart models.Cart
	if err := database.DB.Where("user_id = ? AND product_id = ?", cart.UserID, cart.ProductID).First(&existingCart).Error; err == nil {

		existingCart.Quantity += cart.Quantity

		if existingCart.Quantity > product.Stock {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Total quantity exceeds available stock"})
			return
		}

		if err := database.DB.Save(&existingCart).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update cart quantity"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Cart updated successfully"})
		return

	}

	if err := database.DB.Create(&cart).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not add to cart"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product added to cart successfully","cart":cart})
}

func RemoveFromCart(c *gin.Context) {
	cartID := c.Param("id")
	var cart models.Cart

	if err := database.DB.First(&cart, cartID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart item not found"})
		return
	}

	database.DB.Delete(&cart)
	c.JSON(http.StatusOK, gin.H{"message": "Item removed from cart"})
}

func ViewCart(c *gin.Context) {
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

	var cartItems []models.Cart
	if err := database.DB.Where("user_id = ?", userIDUint).Find(&cartItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cartItems)

}

func UpdateCartQuantity(c *gin.Context) {
	cartID := c.Param("id")
	var cart models.Cart

	if err := database.DB.First(&cart, cartID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart item not found"})
		return
	}

	quantitystr := c.Query("quantity")
	quantity, err := strconv.Atoi(quantitystr)
	if err != nil || quantity < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quantity"})
		return
	}

	cart.Quantity = uint(quantity)
	database.DB.Save(&cart)
	c.JSON(http.StatusOK, gin.H{"message": "Cart quantity updated"})
}
