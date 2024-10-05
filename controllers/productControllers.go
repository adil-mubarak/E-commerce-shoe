package controllers

import (
	"ecommerce/database"
	"ecommerce/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)



func SearchProduct(c *gin.Context){
	var input struct{
		Name string `json:"name"`
	}

	if err := c.BindJSON(&input); err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":"Invalid input"})
		return
	}

	var products []models.Product
	log.Printf("Searching for product with name: %s",input.Name)

	result := database.DB.Where("name LIKE ?","%"+input.Name+"%").Find(&products)

	if result.Error != nil{
		log.Printf("Database query error: %v",result.Error)
		c.JSON(http.StatusInternalServerError,gin.H{"error":"Database query error","details":result.Error.Error()})
		return
	}

	if len(products) == 0{
		c.JSON(http.StatusNotFound,gin.H{"error":"Could not find products"})
		return
	}
	c.JSON(http.StatusOK,products)
}



func GetProducts(c *gin.Context) {
	var products []models.Product
	database.DB.Find(&products)
	c.JSON(http.StatusOK, gin.H{"product": products})
}

func Sorting(c *gin.Context){
	var products []models.Product
	query := database.DB

	sortBy := c.Query("sort_by")
	sortOrder := c.DefaultQuery("sort_order","asc")

	allowedSortBy := map[string]bool{
		"price":true,
		"name":true,
		"stock":true,
		"id": true,
	}

	if sortBy == "" || !allowedSortBy[sortBy]{
		sortBy = "id"
	}

	if sortOrder != "asc" && sortOrder != "desc"{
		sortOrder = "asc"
	}

	query = query.Order(sortBy+" "+sortOrder)

	result := query.Find(&products)
	if result.Error != nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":"Database query error", "details":result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK,products)

}

func Filtering(c *gin.Context){
	var products []models.Product
	query := database.DB

	category := c.Query("category")
	minPrice := c.Query("min_price")
	maxPrice := c.Query("max_price")
	inStock := c.Query("in_stock")

	if category != ""{
		query = query.Where("category = ?",category)
	}

	if minPrice != "" && maxPrice != ""{
		min,errmin := strconv.ParseFloat(minPrice,64)
		max,errmax := strconv.ParseFloat(maxPrice,64)
		if errmin == nil && errmax == nil{
			query = query.Where("price BETWEEN ? AND ?", min,max)
		}
	}

	if inStock == "true"{
		query = query.Where("stock > 0")
	}

	result := query.Find(&products)
	if result.Error != nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":"Database query error","details":result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK,products)
}