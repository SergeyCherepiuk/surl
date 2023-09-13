package handlers

import (
	"net/http"

	"github.com/SergeyCherepiuk/surl/domain"
	"github.com/labstack/echo/v4"
)

type todoHandler struct{}

func NewTodoHandler() *todoHandler {
	return &todoHandler{}
}

func (h todoHandler) Get(c echo.Context) error {
	todos := []domain.Todo{
		{Name: "Todo1", Description: "Desc1", IsComplete: true},
		{Name: "Todo2", Description: "Desc2", IsComplete: false},
		{Name: "Todo3", Description: "Desc3", IsComplete: true},
	}
	return c.Render(http.StatusOK, "components/todo-list", todos)
}
