package main

import (
	"log"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {

	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	
	router := gin.Default() //initialising gin router

	//middleware to pass the database connection to handlers
	router.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	router.POST("/library", createLibrary)
	router.GET("/library", getAllLibrary)
	router.POST("/book/:email", createBook)
	router.GET("/book", getAllBook)
	router.PUT("/book/:isbn/:email", updateBookByISBN)
	router.DELETE("/book/:isbn/:email", deleteBookByISBN)
	router.POST("/user", createUser)
	router.PUT("/user/:id/:email", updateUserByID)
	router.GET("/book/:title", getBookByTitle)
	router.POST("/request/:email", createRequest)
	router.GET("/request/:email", getAllRequest)
	router.PUT("/request/:reqid/:email", updateRequestByReqID)
	router.GET("/user/:email", getAllUser)
	router.POST("/issue/:email", createIssue)

	router.Run(":8870")
}
