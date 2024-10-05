package controllers

import (
	"ecommerce/database"
	"ecommerce/models"
	"ecommerce/tokenjwt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/paymentintent"
)

func ProcessPayment(c *gin.Context) {
	claims, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userClaims, ok := claims.(*tokenjwt.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	userID := userClaims.UserID 
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	var payment models.Payment
	if err := c.ShouldBindJSON(&payment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if payment.OrderID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order ID is required"})
		return
	}

	var order models.Order
	if err := database.DB.Where("id = ?", payment.OrderID).First(&order).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order not found"})
		return
	}

	if order.Status != "Pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order is not pending"})
		return
	}

	if payment.Amount < order.Total {
		remaining := order.Total - payment.Amount
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Insufficient payment",
			"remaining": remaining,
			"total_price":order.Total,
		})
		return
	}

	if payment.Amount > order.Total{
		extra := payment.Amount - order.Total
		c.JSON(http.StatusBadRequest,gin.H{
			"error":"Insufficient payment",
			"extra": extra,
			"total_price":order.Total,
		})
		return
	}

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(payment.Amount * 100)), 
		Currency: stripe.String("inr"),
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment intent"})
		return
	}

	_, err = paymentintent.Confirm(
		pi.ID,
		&stripe.PaymentIntentConfirmParams{
			PaymentMethod: stripe.String("pm_card_visa"), 
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm payment intent"})
		return
	}

	payment.UserID = userID 
	payment.Status = "succeeded"
	payment.PaymentID = pi.ID
	payment.OrderID = order.ID

	if result := database.DB.Create(&payment); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not process payment"})
		return
	}

	order.PaymentStatus = "Paid"
	if result := database.DB.Save(&order); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update order status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Payment successful",
		"order_details": order,
	})
}
