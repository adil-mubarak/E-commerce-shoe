##**This is my E-commerse project with Golang**

# You can start the project with below commands
go run main.go

***USER CAN USE WITHOUT LOGIN***

__**SIGNUP FUNCTION API CALL (POST REQUEST)**

**http://localhost:8888/signup
```json
{
  "name":"adil",
  "email":"adil@gamil.com"
  "password":"1234"
}
```

Response : "User registered successfully"

__**LOGIN FUNCTION API CALL (POST REQUEST)**

**http://localhost:888/login 

```json
{
  "email":"adil@gmail.com"
  "password":"1234"
}
```

Response : "Login successfully"
token : token,
refreshtoken : refreshtoken,
role : user

__**USER CAN SEE ALL PRODUCTS (GET REQUEST)**

**http://localhost:888/products

this time we can see all products
{
	ID   : 1     
	Name  : "nike shoe"      
	Description  : "good quality shoe"
	Price       : 2999
	Stock      :  35
	Category  :  "running shoe"  
	ImageURL :   image url / upload from system  
}

__**USER CAN SORT PRODUCTS (GET REQUEST)**

**http://localhost:888/sort/products

user can sort products on the base of price , stock , catogery , ...

__**USER CAN FILTER PRODUCTS (GET REQUEST)**

**http://localhost:888/filter/products

user can filter the products on the base of price , catogery , stock ,...

__**USER CAN GET THE REFRESH TOKEN (POST REQUEST)**

**http://localhost:888/refresh-token

user refresh the existed token that had expired

***USERS FUNCTIONALITY***

__**USER CAN ADD ITEMS INTO WISHLIST (POST REQUEST)**

**http://localhost:888/user/wishlist
this can use to users to add products into wishlist 

__**USER CAN DELETE WISHLIST ITEMS (DELETE REQUEST)**

**http://localhost:888/user/wishlist/:id
this can use to delete items from wishlist

__**USER CAN SEE ITEMS FROM WISHLIST (GET REQUEST)**

**http://localhost:888/user/wishlist
this can use to see all items from wishlist

__**USER CAN ADD ADDRESS (POST REQUEST)**

**http://localhost:888/user/addresses
this can use to add new Address

__**USER CAN DELETE ADDRESS (DELETE REQUEST)**

**http://localhost:888/user/addresses/:id
this can use to delete addresses

__**USER CAN EDIT ADDRESS (PUT REQUEST)**

**http://localhost:888/user/addresses
this can use edit addresses

__**USER CAN ADD ITEMS INTO CARTS (POST REQUEST)**

**http://localhost:888/user/cart
this can use to add items into carts 
```json
{
  "product_id" : 2,
  "quantity": 2
}
```

__**USER CAN UPDATE CART ITEMS (PUT REQUEST)**

**http://localhost:888/user/cart/:id
this can use to edit the quantity of product
```json
{
  "quantity" : 4
}
```

__**USER CAN DELETE CART ITEMS (DELETE REQUEST)**

**http://localhost:888/user/cart/:id
this can use to remove items from cart

__**USER CAN SEE ALL CART ITEMS (GET REQUEST)**

**http://localhost:888/user/carts
this can use to see all items from the cart

__**USER CAN PLACE ORDER (POST REQUEST)**

**http://localhost:888/user/order
this can use to order the items
```json
{
  "street":"thorakal",
  "city":"manjeri",
  "state":"kerala",
  "phone":4322453245,
  "postal_code":345423,
  "country":"india"
}
```

__**USER CAN SEE THE ORDERS (GET ORDERS)**

**http://localhost:888/user/orders
this can use to see all orders

__**USER NEED TO PAY (POST REQUEST)**

**http://localhost:888/user/payment
this can use to pay the amount of the product
```json
{
  "order_id":3,
  "total_price":3455,
}
```

***ADMIN FUNCATIONALITY***

__**ADMIN CAN ADD PRODUCTS (POST REQUEST)**

**http://localhost:8888/admin/product
this can use to add new products
```json
{
  "name":"puma"
  "description":"good"
  "price":2444,
  "stock":23,
  "category":"sneakers"
  "image_url": image url
}
```

__**ADMIN CAN DELETE PRODUCT (DELETE REQUEST)**

**http://localhost:888/admin/products/:id
this can use to remove products

__**ADMIN CAN UPDATE PRODUCT (PUT REQUEST)**

**http://localhost:888/admin/products/:id
this can use to edit the products

__**ADMIN CAN UPLOAD system FROM SYSTEM (PUT REQUEST)**

**http://localhost:888/admin/product/image/:id
this can use to upload image from system

__**ADMIN CAN SEE ORDERS (GET REQUEST)**

**http://localhost:888/admin/orders
this can use to see all the orders

__**ADMIN CAN SEE ALL USERS (GET REQUEST)**

**http://localhost:888/admin/users
this can ues to see all the users

__**ADMIN CAN BAN USERS (PUT REQUEST)**

**http://localhost:888/admin/banusers/:id
this can use to block the users

__**ADMIN CAN UNBAN USERS (PUT REQUEST)**

**http://localhost:888/admin/unbanuser/:id
this can use to un block the users

__**ADMIN CAN UPDATE THE ORDER STATUS (PUT REQUEST)**

**http://localhost:888/admin/updatestatus/:id
this can use to update the order status to (shipped,deliverd,prnding...)

    



