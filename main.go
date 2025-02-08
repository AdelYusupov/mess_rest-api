package main

import (
	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

type Message struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}
type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

var db *gorm.DB

func initDB() {
	dsn := "host=localhost user=postgres password=yourpassword dbname=postgres port=5432 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	db.AutoMigrate(&Message{})
}

var nextID = 1

func patchHadler(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: "id not correct",
		})
	}
	var updatedMessage Message
	if err := c.Bind(&updatedMessage); err != nil {
		return c.JSON(http.StatusBadRequest, Response{Status: "error", Message: "Invaild input"})
	}

	if err := db.Model(&Message{}).Where("id = ?", id).Update("text", updatedMessage.Text).Error; err != nil {
		return c.JSON(http.StatusBadRequest, Response{Status: "error", Message: "Could not update the message"})
	}

	return c.JSON(http.StatusOK, Response{Status: "success", Message: "Message updated successfully"})
}

func deleteHandler(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: "id not correct",
		})
	}
	if err := db.Where("id = ?", id).Delete(&Message{}).Error; err != nil {
		return c.JSON(http.StatusBadRequest, Response{Status: "error", Message: "Could not delete the message"})
	}
	return c.JSON(http.StatusOK, Response{Status: "success", Message: "Message deleted successfully"})
}

func getHandler(c echo.Context) error {
	var messages []Message

	if err := db.Find(&messages).Error; err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Could not find the messages"})
	}

	return c.JSON(http.StatusOK, &messages)
}
func postHandler(c echo.Context) error {
	var message Message
	if err := c.Bind(&message); err != nil {
		return c.JSON(http.StatusBadRequest, Response{Status: "error", Message: "Could not add the message"})
	}
	if err := db.Create(&message).Error; err != nil {
		return c.JSON(http.StatusBadRequest, Response{Status: "error", Message: "Could not created the message"})
	}
	return c.JSON(http.StatusOK, Response{Status: "ok", Message: "Message created successfully"})
}
func main() {
	initDB()
	e := echo.New()
	e.GET("/messages", getHandler)
	e.POST("/messages", postHandler)
	e.PATCH("/messages/:id", patchHadler)
	e.DELETE("/messages/:id", deleteHandler)
	err := e.Start(":9090")
	if err != nil {
		return
	}
}
