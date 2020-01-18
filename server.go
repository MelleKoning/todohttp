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

// ServerSvc exposes all necessary methods of the server
type ServerSvc interface {
	InsertTodoItem(todoitem *models.TodoItem) (*models.TodoItem, error)
}

// NewServer to instantiate a server ref
func NewServer(todoRepository mongotododb.TodoRepository) (ServerSvc, error) {
	if nil == todoRepository {
		return nil, errors.New("You have to provide a mock or instance for accessing MongoDB")
	}
	return &ServerPackage{
		todoRepository: todoRepository,
	}, nil

}

// InsertTodoItem inserts a new item to the list
func (s *ServerPackage) InsertTodoItem(todoitem *models.TodoItem) (*models.TodoItem, error) {
	// we could do some validation logic here that we can test..
	// as an example, suppose the item to add needs to have a valid status set.

	if err := todoitem.HasValidStatus(); err != nil {
		return nil, err
	}

	return s.todoRepository.Insert(todoitem) // just pass through result of db-insert, can be mocked in test
}
