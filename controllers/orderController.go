package controllers

import (
	"ecommerce/database"
	"ecommerce/models"
	"ecommerce/tokenjwt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/paymentintent"
)

const stripeSecretKey = "sk_test_51Q3VDSA4q5XQgDF4Exe4KJGgrjRZl2NQjCliTKMWDFLk4licXMyUibc3OGJ6IGCpNrR6hPkGw47D3xbu5utTNyDG00mgJvskaw"

func init() {
	stripe.Key = stripeSecretKey
}

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
		total += float64(item.Quantity) * product.Price
	}

	log.Printf("Total price calculated: %f\n", total)

	// var address models.Address
	// if err := database.DB.Where("user_id = ?", userID).First(&address).Error; err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "User address not found"})
	// 	return
	// }

	var newAddress models.Address
	if err := c.ShouldBindJSON(&newAddress); err != nil{
		log.Printf("Error binding address data: %v",err)
		c.JSON(http.StatusBadRequest,gin.H{"error":"Invalid address data","details":err.Error()})
		return
	}

	newAddress.UserID = userID
	if err := database.DB.Create(&newAddress).Error; err != nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":"Failed to save address","details":err.Error()})
		return
	}

	order := models.Order{
		UserID:    userID,
		Total:     total,
		// AddressID: address.ID,
		AddressID: newAddress.ID,
		CreatedAt: time.Now(),
	}

	if err := database.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not place order", "details": err.Error()})
		return
	}

	amountInPaisa := int64(total * 100)

	params := &stripe.PaymentIntentParams{
		Amount:             stripe.Int64(amountInPaisa),
		Currency:           stripe.String("inr"),
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create payment intent",
			"details": err.Error(),
		})
		return
	}

	payment := models.Payment{
		UserID:        userID,
		Amount:        total,
		Status:        "Pending",
		PaymentID:     pi.ID,
		OrderID:       order.ID,
		PaymentMethod: "Stripe",
	}

	if err := database.DB.Create(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to save payment",
			"details": err.Error(),
		})
		return
	}

	log.Printf("Checkout saved successfully: %v", payment)

	c.JSON(http.StatusOK, gin.H{
		"message":       "Checkout successful",
		"total_price":   total,
		"client_secret": pi.ClientSecret,
		"payment_id":    pi.ID,
		"Order_id":      order.ID,
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


