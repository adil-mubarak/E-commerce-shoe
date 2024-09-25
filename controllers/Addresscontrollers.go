package controllers

import (
	"ecommerce/database"
	"ecommerce/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUserAddresses(c *gin.Context) {
	userID := c.Param("id")
	var user models.User

	if err := database.DB.Preload("Addresses").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user.Addresses)
}

func CreateAddress(c *gin.Context){
	var address models.Address
	if err := c.ShouldBindJSON(&address); err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		return
	}

	if err := database.DB.Create(&address).Error; err != nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":"Unable to create address"})
		return
	}

	c.JSON(http.StatusCreated,address)
}

func UpdateAddress(c *gin.Context){
	var address models.Address
	id := c.Param("id")

	if err := database.DB.First(&address,id).Error; err != nil{
		c.JSON(http.StatusNotFound,gin.H{"error":"Address not found"})
		return
	}

	if err := c.ShouldBindJSON(&address); err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		return
	}

	database.DB.Save(&address)
	c.JSON(http.StatusOK, gin.H{"address":address})
}

func DeleteAddress(c *gin.Context){
	id := c.Param("id")
	var address models.Address
	if err := database.DB.Delete(&address,id).Error; err != nil{
		c.JSON(http.StatusNotFound,gin.H{"error":"Address not found"})
		return
	}

	c.JSON(http.StatusNoContent,nil)
}