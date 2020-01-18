# todohttp
a test project todolist
bare http trial.

# mongodb

if no docker image is available first:
* docker pull mongo:latest

folder
.\docker 
contains the script to start a mongodb instance in docker with
* docker-compose -f ./docker/docker-compose.yml up -d 

on windows a volume is needed to store data
* docker volume create --name=mongodata
* docker run -d -p 27017:27017 -v mongodata:/data/db --name=tododb-mongo-container mongo

# code todoitems

This is a bare implementation using HTTP.

THe following commands were used to setup the project dependencies
* dep init
* dep ensure

# initial todo item struct

// TodoItem reference for the item in mongo
type TodoItem struct {
	Text      string    `json:"text" bson:"text"`
	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
}

# running program

in visual studio code just press F5 to run the main.go program

in-browser POST message:
http://localhost:8080/addtodoitem
payload:
{
"text": "sometext2"
}

returns:
{
  "text": "sometext2",
  "createdAt": "2020-01-18T09:49:15.8543081Z"
}

read items back with GET:
http://localhost:8080/todolist

for example returns:
[
  {
    "text": "sometext2",
    "createdAt": "2020-01-18T10:49:15.854+01:00"
  },
  {
    "text": "sometext",
    "createdAt": "2020-01-18T10:48:43.142+01:00"
  }
]

# updating items

To be able to update items, for example change the state from TODO to BUSY or DONE, we need a /updatetodoitem POST and also exchange
the actual underlying key so that frontend knows what item to actually update.

updating struct to include the id and a status:

// TodoItem reference for the item in mongo
type TodoItem struct {
	ID        bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	Text      string        `json:"text" bson:"text"`
	Status    string        `json:"status" bson:"status"`
	CreatedAt time.Time     `json:"createdAt" bson:"created_at"`
}

now todolist returns the mongo-db-id and the status that is still empty (as not yet updated):
[
  {
    "_id": "5e22d49bdb7638c86a902828",
    "text": "sometext2",
    "status": "",
    "createdAt": "2020-01-18T10:49:15.854+01:00"
  },
  {
    "_id": "5e22d47cdb7638c86a902814",
    "text": "sometext",
    "status": "",
    "createdAt": "2020-01-18T10:48:43.142+01:00"
  }
]

now adding update call to handle updates for existing item