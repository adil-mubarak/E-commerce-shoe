package controllers

import (
	"ecommerce/database"
	"ecommerce/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddToCart(c *gin.Context) {
	var cart models.Cart
	if err := c.ShouldBindJSON(&cart); err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		return
	}

	var product models.Product
	if err := database.DB.First(&product,cart.ProductID).Error;err != nil{
		c.JSON(http.StatusNotFound,gin.H{"error":"Product not found"})
		return
	}

	if err := database.DB.Create(&cart).Error; err != nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":"Could not add to cart"})
		return
	}
	c.JSON(http.StatusOK,gin.H{"message":"Product added to cart successfully"})
}

func RemoveFromCart(c *gin.Context){
	cartID := c.Param("id")
	var cart models.Cart

	if err := database.DB.First(&cart,cartID).Error; err != nil{
		c.JSON(http.StatusNotFound,gin.H{"error":"Cart item not found"})
		return
	}

	database.DB.Delete(&cart)
	c.JSON(http.StatusOK,gin.H{"message":"Item removed from cart"})
}

func ViewCart(c *gin.Context){
	userID,exists := c.Get("user_id")
	if !exists{
		c.JSON(http.StatusBadRequest,gin.H{"error":"User not authenticated"})
		return
	}

	var cartItems []models.Cart
	database.DB.Where("user_id = ?",userID).Find(&cartItems)
	c.JSON(http.StatusOK,cartItems)
}

func UpdateCartQuantity(c *gin.Context){
	cartID := c.Param("id")
	var cart models.Cart

	if err := database.DB.First(&cart,cartID).Error; err != nil{
		c.JSON(http.StatusNotFound,gin.H{"error":"Cart item not found"})
		return
	}

	quantitystr := c.Query("quantity")
	quantity,err := strconv.Atoi(quantitystr)
	if err != nil || quantity < 1{
		c.JSON(http.StatusBadRequest,gin.H{"error":"Invalid quantity"})
		return
	}

	cart.Quantity = uint(quantity)
	database.DB.Save(&cart)
	c.JSON(http.StatusOK,gin.H{"message":"Cart quantity updated"})
}