package main

import (
	"errors"

	models "github.com/MelleKoning/todohttp/models"
	mongotododb "github.com/MelleKoning/todohttp/mongo"
)

// ServerPackage serves the concrete server..
type ServerPackage struct {
	todoRepository mongotododb.TodoRepository
}

// NewServer to instantiate a server ref
func NewServer(todoRepository mongotododb.TodoRepository) (*ServerPackage, error) {
	if nil == todoRepository {
		return nil, errors.New("You have to provide a mock or instance for accessing MongoDB")
	}
	return &ServerPackage{
		todoRepository: todoRepository,
	}, nil

}

// InsertTodoItem inserts a new item to the list
func (s *ServerPackage) InsertTodoItem(todoitem *models.TodoItem) (*models.TodoItem, error) {
	return s.todoRepository.Insert(todoitem) // just pass through for now
}
