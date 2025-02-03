package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type Message struct {
	Text string `json:"text"`
}
type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

var messages []Message

func getHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, &messages)
}
func postHandler(c echo.Context) error {
	var message Message
	if err := c.Bind(&message); err != nil {
		return c.JSON(http.StatusBadRequest, Response{Status: "error", Message: "Could not add the message"})
	}
	messages = append(messages, message)
	return c.JSON(http.StatusCreated, Response{Status: "ok", Message: "Message added successfully"})
}
func main() {
	e := echo.New()
	e.GET("/messages", getHandler)
	e.POST("/messages", postHandler)
	err := e.Start(":9090")
	if err != nil {
		return
	}
}
