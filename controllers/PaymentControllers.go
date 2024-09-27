package controllers

import (
	"ecommerce/database"
	"ecommerce/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/paymentintent"
)

func init() {
	stripe.Key = "sk_test_51Q3VDSA4q5XQgDF4Exe4KJGgrjRZl2NQjCliTKMWDFLk4licXMyUibc3OGJ6IGCpNrR6hPkGw47D3xbu5utTNyDG00mgJvskaw"
}

func InitiatePayment(c *gin.Context) {
	var reqBody struct{
		OrderID uint `json:"order_id"`
		UserID uint `json:"user_id"`
		Total float64 `json:"total"`
	}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	amountInPaisa := int64(reqBody.Total * 100)

	params := &stripe.PaymentIntentParams{
		Amount:             stripe.Int64(amountInPaisa),
		Currency:           stripe.String("INR"),
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
		OrderID:       reqBody.OrderID,
		UserID:        reqBody.UserID,
		Amount:        reqBody.Total,
		Status:        "Pending",
		PaymentID:     pi.ID,
		PaymentMethod: "Stripe",
	}

	if err := database.DB.Create(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to save payment",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"client_secret": pi.ClientSecret,
		"amount":       float64(amountInPaisa)/100,
		"currency":      "INR",
	})
}

func HandlePaymentSuccess(c *gin.Context) {
	var paymentDetails struct {
		PaymentIntentID string `json:"payment_intent_id"`
	}

	if err := c.ShouldBindJSON(&paymentDetails); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var payment models.Payment
	if err := database.DB.Where("payment_id = ?", paymentDetails.PaymentIntentID).First(&payment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	var order models.Order
	if err := database.DB.First(&order, payment.OrderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	order.PaymentStatus = "Paid"
	payment.Status = "Paid"

	if err := database.DB.Save(&payment).Error; err != nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":"Failed to update payment status"})
		return
	}

	if err := database.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order status", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment processed successfully"})
}
