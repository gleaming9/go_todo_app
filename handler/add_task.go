package handler

import (
	"github.com/go-playground/validator"
	"go_todo_app/store"
)

type AddTask struct {
	Store     *store.TaskStore
	Validator *validator.Validate
}
