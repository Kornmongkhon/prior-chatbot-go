package handler

import (
	"fmt"
	"prior-chat-bot/internal/core/port"
)

func ExecuteHandlerHealthCheck(c port.MyServer) {
	fmt.Println("ExecuteHandlerHealthCheck")
	c.ToResponse(200, "S0000", "Success", nil)
}
