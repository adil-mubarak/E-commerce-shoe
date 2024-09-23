package controllers

import (
	"ecommerce/database"
	"ecommerce/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddToWishlist(c *gin.Context){
	var wishlist models.Whishlist
	if err := c.ShouldBindJSON(&wishlist); err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		return
	}

	var product models.Product
	if err := database.DB.First(&product,wishlist.ProductID).Error; err != nil{
		c.JSON(http.StatusNotFound,gin.H{"error":"Product not found"})
		return
	}

	if err := database.DB.Create(&wishlist).Error; err != nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":"Could not add to wishlist"})
		return
	}
	c.JSON(http.StatusOK,gin.H{"message":"Product added to wishlist"})
}

func RemoveFromWishlist(c *gin.Context){
	wishlistID := c.Param("id")
	var wishlist models.Whishlist

	if err := database.DB.First(&wishlist,wishlistID).Error; err != nil{
		c.JSON(http.StatusNotFound,gin.H{"error":"Wishlist item not found"})
		return
	}

	database.DB.Delete(&wishlist)
	c.JSON(http.StatusOK,gin.H{"message":"Item removed from wishlist"})
}

func ViewWishlist(c *gin.Context){
	userID,exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusBadRequest,gin.H{"error":"User not authenticated"})
		return
	}

	var wishlistItems []models.Whishlist
	database.DB.Where("user_id = ?",userID).Find(&wishlistItems)
	c.JSON(http.StatusOK,wishlistItems)
}