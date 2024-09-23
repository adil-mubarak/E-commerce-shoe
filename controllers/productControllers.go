package controllers

import (
	"ecommerce/database"
	"ecommerce/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if result := database.DB.Create(&product); result.Error != nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":result.Error.Error()})
		return
	}
	c.JSON(http.StatusCreated,gin.H{"product":product})
}

func UpdateProduct(c *gin.Context){
	id := c.Param("id")
	var product models.Product

	if result := database.DB.First(&product,id); result.Error != nil{
		c.JSON(http.StatusNotFound,gin.H{"error":"Product not found"})
		return
	}

	if err := c.ShouldBindJSON(&product); err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"eror":err.Error()})
		return
	}

	database.DB.Save(&product)
	c.JSON(http.StatusOK,gin.H{"product":product})
}

func DeleteProduct(c *gin.Context){
	id := c.Param("id")
	var product models.Product

	if result := database.DB.First(&product,id); result.Error != nil{
		c.JSON(http.StatusNotFound,gin.H{"error":"Product not found"})
		return
	}

	database.DB.Delete(&product)
	c.JSON(http.StatusOK,gin.H{"message":"Product deleted successfully"})

}

func GetProducts(c *gin.Context){
	var products []models.Product
	database.DB.Find(&products)
	c.JSON(http.StatusOK,gin.H{"product":products})
}

