package mongotododb

import (
	"log"

	"github.com/MelleKoning/todohttp/models"
	"gopkg.in/mgo.v2"
)

var databaseName = "mongotododb"
var collectionTodoitems = "todoitems"

// TodoRepository  is the interface for database operations related to todoitems
type TodoRepository interface {
	Insert(todoItem *models.TodoItem) (*models.TodoItem, error)
}

// TodoDatabase handles operations on the mongotododb database
type TodoDatabase struct {
	endpoint string
	session  *mgo.Session
	// unresponsiveThreshold int
}

// NewTodoDatabase creates a new instance of the mongotododb
// so that sessions can be opened/closed
func NewTodoDatabase(endpoint string) (TodoRepository, error) {
	// Create initial session
	session, err := mgo.Dial(endpoint)
	if err != nil {
		return nil, err
	}

	return &TodoDatabase{endpoint, session}, nil
}

// Insert inserts the item and returns the added item to include the created Bson.Id
func (t *TodoDatabase) Insert(todoItem *models.TodoItem) (*models.TodoItem, error) {
	session := t.session.Copy()
	defer session.Close()

	err := session.DB(databaseName).C(collectionTodoitems).Insert(todoItem)

	result := &models.TodoItem{}
	// we also want to return the created bson.Id so we Find back the last inserted item
	// the new mongo driver has InsertOne, but mgo.v2 does not have that yet
	if err := session.DB(databaseName).C(collectionTodoitems).Find(nil).Sort("-created_at").One(&result); err != nil {
		log.Print("could not find back the inserted item")
		return nil, err
	}
	return result, err
}

/*
func (t *TodoRepository) Delete(id string) error {
	session := t.session.Copy()
	defer session.Close()

	query := bson.M{taskIDKey: id}
	err := session.DB(databaseName).C(taskCollectionName).Remove(query)
	if err != nil {
		log.Errorf("Could not delete task with provided id: %v", id)
		return err
	}

	return nil
}
*/
