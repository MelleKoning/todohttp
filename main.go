package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/rs/cors"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
)

// TodoItem reference for the item in mongo
type TodoItem struct {
	ID        bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	Text      string        `json:"text" bson:"text"`
	Status    string        `json:"status" bson:"status"`
	CreatedAt time.Time     `json:"createdAt" bson:"created_at"`
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
var todoitems *mgo.Collection

func main() {
	// Connect to mongo
	session, err := mgo.Dial("mongo1:27017")
	if err != nil {
		log.Fatalln(err)
		log.Fatalln("mongo err")
		os.Exit(1)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	// Get todo collection
	todoitems = session.DB("mongotododb").C("todoitems")

	// Set up routes
	r := mux.NewRouter()
	r.HandleFunc("/addtodoitem", createTodoItem).
		Methods("POST")
	r.HandleFunc("/todolist", readTodoItems).
		Methods("GET")
	r.HandleFunc("/updatetodoitem", updateTodoItem).Methods("POST")
	http.ListenAndServe(":8080", cors.AllowAll().Handler(r))
	log.Println("Listening on port 8080...")
}

func createTodoItem(w http.ResponseWriter, r *http.Request) {
	// Read body
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Read todoItem
	todoitem := &TodoItem{}
	err = json.Unmarshal(data, todoitem)
	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}
	todoitem.CreatedAt = time.Now().UTC()

	// Insert new item
	if err := todoitems.Insert(todoitem); err != nil {
		responseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseJSON(w, todoitem)
}

func readTodoItems(w http.ResponseWriter, r *http.Request) {
	result := []TodoItem{}
	if err := todoitems.Find(nil).Sort("-created_at").All(&result); err != nil {
		responseError(w, err.Error(), http.StatusInternalServerError)
	} else {
		responseJSON(w, result)
	}
}

func updateTodoItem(w http.ResponseWriter, r *http.Request) {
	todoitem, err := readTodoItemFromBody(w, r)
	if err != nil {
		log.Print("Could not read item from json")
		return
	}
	result := TodoItem{}
	if err := todoitems.FindId(todoitem.ID).One(&result); err != nil {
		log.Print("Could not find that item, it should already exist")
		responseError(w, err.Error(), http.StatusInternalServerError)
	}
	// yes, valid item, please update
	if info, err := todoitems.UpsertId(todoitem.ID, &todoitem); err != nil {
		log.Printf("Upsert failed %v %v", info, err)
		responseError(w, err.Error(), http.StatusInternalServerError)
	}
	result = *todoitem
	responseJSON(w, result)

}

func readTodoItemFromBody(w http.ResponseWriter, r *http.Request) (*TodoItem, error) {
	// Read body
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}

	// Read todoItem
	todoitem := &TodoItem{}
	err = json.Unmarshal(data, todoitem)
	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}
	return todoitem, nil
}
func responseError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func responseJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
