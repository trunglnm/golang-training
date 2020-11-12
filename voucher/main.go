package main

import (
	"database/sql"
	"log"
	"os"

	"./model"
	"./storage"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db, err := sql.Open("mysql", "default:secret@/voucher")
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(30)
	defer db.Close()
	r := gin.Default()
	r.POST("/register", func(c *gin.Context) {
		voucher := model.Voucher{}
		voucherStorage := storage.Voucher{
			DB: db,
		}
		// S1: Lay data tu client
		if err := c.ShouldBindJSON(&voucher); err != nil {
			c.JSON(400, err)
			return
		}
		// S2: Kiem tra co ton tai trong db khong
		// isExist, err := voucherStorage.IsExit(voucher)
		// if err != nil {
		// 	c.JSON(400, err)
		// 	return
		// }

		// if isExist {
		// 	c.JSON(400, gin.H{
		// 		"error": "Exist",
		// 	})
		// 	return
		// }
		// S3: Insert vo db neu no chua ton tai
		voucherStorage.RegisterIsolation(&voucher)
		c.JSON(200, voucher)
	})
	r.GET("/verify", func(c *gin.Context) {
		voucher := model.Voucher{}
		c.JSON(200, voucher)
	})
	port := os.Getenv("PORT")
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong from " + port,
		})
	})
	r.Run(":" + port) // listen and serve on 0.0.0.0:8080
}
