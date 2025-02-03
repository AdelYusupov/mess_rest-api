package main

import (
	"github.com/labstack/echo/v4"
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

var messages = make(map[int]Message)
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
		return c.JSON(http.StatusBadRequest, Response{Status: "error", Message: "Could not update the message"})
	}
	if _, exists := messages[id]; !exists {
		return c.JSON(http.StatusBadRequest, Response{Status: "error", Message: "Message was not found"})
	}
	updatedMessage.ID = id
	messages[id] = updatedMessage
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
	if _, exists := messages[id]; !exists {
		return c.JSON(http.StatusBadRequest, Response{Status: "error", Message: "Message was not found"})
	}
	delete(messages, id)
	return c.JSON(http.StatusOK, Response{Status: "success", Message: "Message deleted successfully"})
}

func getHandler(c echo.Context) error {
	var msgSlice []Message

	for _, msg := range messages {
		msgSlice = append(msgSlice, msg)
	}
	return c.JSON(http.StatusOK, msgSlice)
}
func postHandler(c echo.Context) error {
	var message Message
	if err := c.Bind(&message); err != nil {
		return c.JSON(http.StatusBadRequest, Response{Status: "error", Message: "Could not add the message"})
	}
	message.ID = nextID
	nextID++
	messages[message.ID] = message
	return c.JSON(http.StatusCreated, Response{Status: "ok", Message: "Message added successfully"})
}
func main() {
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
