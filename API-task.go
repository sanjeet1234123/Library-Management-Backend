package main

import (
	"log"
    _ "github.com/lib/pq"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"

)

func main() {
	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := gin.Default()

	
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"} 
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	router.Use(cors.New(config))

	router.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	router.POST("/library", createLibrary)
	router.GET("/library", getAllLibrary)
	router.POST("/book/create", createBook)
	router.GET("/book", getAllBook)
	router.PUT("/update/book/:isbn", updateBookByISBN)
	router.DELETE("/book/:isbn", deleteBookByISBN)
	router.POST("/user", createUser)
	router.POST("/user/login",loginUser)
	router.PUT("/user/:id", updateUserByID)
	router.GET("/book/:isbn", getBookByISBN)
	router.POST("/request", createRequest)
	router.GET("/getrequest", getAllRequest)
	router.PUT("/request/:reqid", updateRequestByReqID)
	router.GET("/user/:email", getAllUser)
	router.POST("/issue", createIssue)

	router.Run(":8870")
}