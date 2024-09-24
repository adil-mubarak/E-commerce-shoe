package routes

import (
	"ecommerce/controllers"
	"ecommerce/middlewares"

	"github.com/gin-gonic/gin"
)

func SetUpRouter() *gin.Engine{
	router := gin.Default()

	router.POST("/signup",controllers.Register)
	router.POST("/login",controllers.Login)
	router.GET("/products",controllers.GetProducts)

	user := router.Group("/user")
	user.Use(middlewares.AuthMiddleWare("user"))
	{
		user.POST("/cart",controllers.AddToCart)
		user.GET("/cart",controllers.ViewCart)
		user.PUT("/cart/:id",controllers.UpdateCartQuantity)
		user.DELETE("/cart/:id",controllers.RemoveFromCart)

		user.POST("/wishlist",controllers.AddToWishlist)
		user.GET("/wishlist",controllers.ViewWishlist)
		user.DELETE("/wishlist/:id",controllers.RemoveFromWishlist)

		user.POST("/order",controllers.CheckOutOrder)
		user.GET("/orders",controllers.GetOrders)
	}


	admin := router.Group("/admin")
	admin.Use(middlewares.AdminAuthMiddleware())
	{
		admin.POST("/products",controllers.CreateProduct)
		admin.PUT("/products/:id",controllers.UpdateProduct)
		admin.DELETE("/products/:id",controllers.DeleteProduct)
		admin.GET("/orders",controllers.GetAllOrder)
	}

	router.POST("/logout",controllers.Logout)

	return router

}