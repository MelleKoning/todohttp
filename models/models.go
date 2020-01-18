package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// TodoItem reference for the item in mongo
type TodoItem struct {
	ID        bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	Text      string        `json:"text" bson:"text"`
	Status    string        `json:"status" bson:"status"`
	DueDate   time.Time     `json:"duedate" bson:"duedate"`
	CreatedAt time.Time     `json:"createdAt" bson:"created_at"`
}

// TodoItemStatusItem is one of the available status,
// should be available to frontend as a list of known status, for example (TODO, BUSY, DONE)
type TodoItemStatusItem struct {
	// ID        bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	Status string `json:"status" bson:"status"`
}

/*
example json
{
  "todo_item": {
    "todoitem_id": "34",
    "description": "string",
    "time_due": "2020-01-17T13:57:57.974Z",
    "status": "BUSY",
    "todo_labels": [
      "RED"
    ]
  }
}
*/
