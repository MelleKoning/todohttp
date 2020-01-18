package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/MelleKoning/todohttp/models"
	"github.com/rs/cors"

	mongotododb "github.com/MelleKoning/todohttp/mongo"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
)

var todoitems *mgo.Collection

var todoRepository mongotododb.TodoRepository
var serverService ServerSvc

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

	initializeServerPackage()

	// Set up routes
	r := mux.NewRouter()
	r.HandleFunc("/addtodoitem", createTodoItem).Methods("POST")
	r.HandleFunc("/todolist", readTodoItems).Methods("GET")
	r.HandleFunc("/updatetodoitem", updateTodoItem).Methods("POST") // or PUT
	r.HandleFunc("/deletetodoitem", deleteTodoItem).Methods("DELETE")

	r.HandleFunc("/todoitemstatus", readTodoItemStatus).Methods("GET") // masterdata... labels..

	http.ListenAndServe(":8080", cors.AllowAll().Handler(r))
	log.Println("Listening on port 8080...")
}

func initializeServerPackage() {

	// TODO setup repository, now connected to mongo1 docker instance
	// which is 127.0.0.1 on dev machine, should come from ENV var
	todoRepository, err := mongotododb.NewTodoDatabase("mongo1:27017")
	if err != nil {
		log.Fatalln(err)
		log.Fatalln("mongo err")
		os.Exit(1)
	}
	// inject the repository into the ServerPackage
	serverService, err = NewServer(todoRepository)
	if err != nil || serverService == nil {
		log.Fatalln(err)
		log.Fatalln("ServerPackage initialization err")
		os.Exit(1)
	}
	return

}
func createTodoItem(w http.ResponseWriter, r *http.Request) {
	todoitem, err := readTodoItemFromBody(w, r)
	if err != nil {
		log.Print("Could not read item from json")
		return
	}
	todoitem.CreatedAt = time.Now().UTC()

	// wire now to the serverService instance/ServerPackage struct
	// which has the status-validation 'business logic'
	result, err := serverService.InsertTodoItem(todoitem)
	if err != nil {
		responseError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	responseJSON(w, result)
}

func readTodoItems(w http.ResponseWriter, r *http.Request) {
	result := []models.TodoItem{}
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
	result := models.TodoItem{}
	if err := todoitems.FindId(todoitem.ID).One(&result); err != nil {
		log.Print("Could not find that item, it should exist")
		responseError(w, err.Error(), http.StatusGone)
		return
	}
	// yes, valid item, please update
	if info, err := todoitems.UpsertId(todoitem.ID, &todoitem); err != nil {
		log.Printf("Upsert failed %v %v", info, err)
		responseError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	result = *todoitem
	responseJSON(w, result)

}

func deleteTodoItem(w http.ResponseWriter, r *http.Request) {
	todoitem, err := readTodoItemFromBody(w, r)
	if err != nil {
		log.Print("Could not read item from json")
		return
	}

	// find item that is to be deleted
	result := models.TodoItem{}
	if err := todoitems.FindId(todoitem.ID).One(&result); err != nil {
		log.Print("Could not find that item, it should exist to be deletable")
		responseError(w, err.Error(), http.StatusGone)
		return
	}

	// delete the found item
	filter := bson.M{"_id": result.ID}
	if err := todoitems.Remove(filter); err != nil {
		log.Printf("Remove failed %v", err)
		responseError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	result = *todoitem
	responseJSON(w, result)

}
func readTodoItemFromBody(w http.ResponseWriter, r *http.Request) (*models.TodoItem, error) {
	// Read body
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}

	// Read todoItem
	todoitem := &models.TodoItem{}
	err = json.Unmarshal(data, todoitem)
	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}
	return todoitem, nil
}

func readTodoItemStatus(w http.ResponseWriter, r *http.Request) { //} []*TodoItemStatusItem {
	todoitemstatuslist := []*models.TodoItemStatusItem{

		&models.TodoItemStatusItem{
			Status: "TODO",
		},
		&models.TodoItemStatusItem{
			Status: "BUSY",
		},
		&models.TodoItemStatusItem{
			Status: "DONE",
		},
	}
	responseJSON(w, todoitemstatuslist)

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
