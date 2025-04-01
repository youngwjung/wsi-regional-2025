package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

var users = []User{
	{ID: 1, Name: "Alice", Age: 25},
	{ID: 2, Name: "Bob", Age: 30},
	{ID: 3, Name: "Charlie", Age: 35},
}

var products = []Product{
	{ID: 1, Name: "Laptop", Price: 999.99},
	{ID: 2, Name: "Smartphone", Price: 499.99},
	{ID: 3, Name: "Tablet", Price: 299.99},
}

func authMiddleware(c *gin.Context) {
	resp, err := http.Get("http://auth-server.local")
	if err != nil || resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		c.Abort()
		return
	}
	c.Next()
}

func main() {
	r := gin.Default()

	r.Use(authMiddleware)

	r.GET("/api/v1/users", func(c *gin.Context) {
		c.JSON(200, users)
	})

	r.GET("/api/v1/user/:id", func(c *gin.Context) {
		id := c.Param("id")
		for _, user := range users {
			if id == string(rune(user.ID)) {
				c.JSON(200, user)
				return
			}
		}
		c.JSON(404, gin.H{"message": "User not found"})
	})

	r.GET("/api/v1/products", func(c *gin.Context) {
		c.JSON(200, products)
	})

	r.GET("/api/v1/product/:id", func(c *gin.Context) {
		id := c.Param("id")
		for _, product := range products {
			if id == string(rune(product.ID)) {
				c.JSON(200, product)
				return
			}
		}
		c.JSON(404, gin.H{"message": "Product not found"})
	})

	r.Run(":8888")
}
