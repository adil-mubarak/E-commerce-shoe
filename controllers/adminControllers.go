package controllers

import (
	"ecommerce/database"
	"ecommerce/models"
	"ecommerce/tokenjwt"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func GetAllUsers(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	claims, err := tokenjwt.ValidateToken(tokenString)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	if claims.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden,Admins only"})
		return
	}

	limit, offset := Paginate(c)
	var total int64

	var users []models.User
	if err := database.DB.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	database.DB.Model(&models.User{}).Count(&total)

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"pagination": gin.H{
			"total": total,
			"page":  c.Query("page"),
			"limit": limit,
			"pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

func UpdateOrderStatus(c *gin.Context) {
	orderID := c.Param("id")
	var Order models.Order

	if err := database.DB.First(&Order, orderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	var input struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	Order.Status = input.Status
	if err := database.DB.Save(&Order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update order status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order status updated successfully", "ordereditems": Order})

}

func BanUser(c *gin.Context) {
	userID := c.Param("id")
	var user models.User

	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.Banned = true
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not ban user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User banned successfully", "user": user})
}

func UnBanUser(c *gin.Context) {
	userID := c.Param("id")
	var user models.User

	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.Banned = false
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not unban user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User unbanned successfullt", "user": user})

}

func CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if result := database.DB.Create(&product); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"product": product})
}

func UploadFile(c *gin.Context) {
	productID := c.Param("id")

	var product models.Product
	if err := database.DB.First(&product, productID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	uploadDir := "./uploads/products"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err := os.MkdirAll(uploadDir, os.ModePerm)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create upload directory"})
			return
		}
	}

	fileName := fmt.Sprintf("%d_%s", product.ID, file.Filename)
	filePath := filepath.Join(uploadDir, fileName)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	product.ImageURL = filePath
	if err := database.DB.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product with image URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Product image updated successfully",
		"Product":    product,
		"image_path": product.ImageURL,
	})
}

func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	if result := database.DB.First(&product, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"eror": err.Error()})
		return
	}

	database.DB.Save(&product)
	c.JSON(http.StatusOK, gin.H{"product": product})
}

func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	if result := database.DB.First(&product, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	database.DB.Delete(&product)
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})

}

func GetAllOrder(c *gin.Context) {
	var orders []models.Order
	database.DB.Find(&orders)
	c.JSON(http.StatusOK, orders)
}
