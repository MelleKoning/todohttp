package server

import (
	"errors"
	"testing"

	mocks "github.com/MelleKoning/todohttp/mocks"
	"github.com/stretchr/testify/mock"

	"github.com/MelleKoning/todohttp/models"

	"github.com/stretchr/testify/assert"
)

// Resembles ServerPackage
type ServerFactory struct {
	todoRepositoryMock *mocks.TodoRepository
}

// CreateServerFactory creates server with mocked database
// so that we can fake database interaction and test on results
func CreateServerFactory() *ServerFactory {
	return &ServerFactory{
		todoRepositoryMock: &mocks.TodoRepository{},
	}
}

// Init initializes a Server instance based on the ServerFactory
func (s *ServerFactory) Init() (ServerSvc, error) {
	if s.todoRepositoryMock == nil {
		return nil, errors.New("assignment can not be nil")
	}
	return NewServer(s.todoRepositoryMock)

}

func TestAddTodoItem_ThrowsError_BecauseDatabaseError(t *testing.T) {

	// Arrange a server but now with the mocked repository interface
	f := CreateServerFactory()
	s, _ := f.Init()

	todoItemToAdd := &models.TodoItem{
		Text:   "add a new item",
		Status: "TODO",
	}

	// arrange some error on the databaserepo
	f.todoRepositoryMock.On("Insert", mock.Anything, mock.Anything).Return(nil, errors.New("Database error"))

	// Act
	newTodoItem, err := s.InsertTodoItem(todoItemToAdd)

	// Assert
	assert.Nil(t, newTodoItem) // no item created
	assert.Error(t, err)
	assert.Equal(t, "Database error", err.Error())
}

func TestAddTodoItem_ThrowsError_BecauseHasInvalidStatus(t *testing.T) {

	// Arrange a server but with the mocked repository interface

	f := CreateServerFactory()
	s, _ := f.Init()

	todoItemToAdd := &models.TodoItem{
		Text:   "add a new item",
		Status: "FRONTEND_THOUGHTUP_STAZTUS", // this is invalid..
	}

	// arrange as if database insert would go all right, if executed
	f.todoRepositoryMock.On("Insert", mock.Anything, mock.Anything).Return(todoItemToAdd, nil)

	// Act... this is to test some 'business logic' of the model
	newTodoItem, err := s.InsertTodoItem(todoItemToAdd)

	// Assert
	assert.Nil(t, newTodoItem) // no item created
	assert.Error(t, err)
	assert.Equal(t, "status not valid, must be in TODO, BUSY, DONE", err.Error())
}
