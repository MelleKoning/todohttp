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
(Run once) Create a external named volume for MongoDB
* docker volume create --name=mongodata
* docker run -d -p 27017:27017 -v mongodata:/data/db --name=tododb-mongo-container mongo

code todoitems
==============
This is a bare implementation using HTTP.

THe following commands were used to setup the project dependencies
* dep init
* dep ensure

initial todo item struct
------------------------
// TodoItem reference for the item in mongo
type TodoItem struct {
	Text      string    `json:"text" bson:"text"`
	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
}

running program
---------------
in visual studio code just press F5 to run the main.go program
and set some break points

adding a todoitem
-----------------
in-browser POST message:
http://localhost:8080/addtodoitem
payload:

```
{
"text": "sometext2"
}
```

returns:
```
{
  "text": "sometext2",
  "createdAt": "2020-01-18T09:49:15.8543081Z"
}
```

read items back with GET:
http://localhost:8080/todolist

for example returns:
```
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
```

# updating items

To be able to update items, for example change the state from TODO to BUSY or DONE, we need a /updatetodoitem POST and also exchange
the actual underlying key so that frontend knows what item to actually update.

updating struct to include the id and a status:

// TodoItem reference for the item in mongo
```
type TodoItem struct {
	ID        bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	Text      string        `json:"text" bson:"text"`
	Status    string        `json:"status" bson:"status"`
	CreatedAt time.Time     `json:"createdAt" bson:"created_at"`
}
```

now todolist returns the mongo-db-id and the status that is still empty (as not yet updated):
```
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
```

now adding update call to handle updates for existing item
POST
http://localhost:8080/updatetodoitem
payload: todoitem
```
{
    "_id": "5e22d49bdb7638c86a902828",
    "text": "sometext2",
    "status": "done",
    "createdAt": "2020-01-01T00:00:00Z",
    "duedate" : "2020-02-02T00:00:00Z"
}
```


adding delete call
http://localhost:8080/deletetodoitem DELETE
payload: todoitem
```
{
    "_id": "5e22fa09db7638c86a902f1c",
    "text": "comment of the todoitem",
    "status": "todo",
    "duedate": "0001-01-01T00:00:00Z",
    "createdAt": "2020-01-18T13:28:57.671+01:00"
  }
```

  # Testing 

  Testing is done via browser, all integration testing.

UML (Mermaid)
=============

as used on https://mermaid-js.github.io/mermaid-live-editor/

[![](https://mermaid.ink/img/eyJjb2RlIjoic2VxdWVuY2VEaWFncmFtXG5Gcm9udGVuZCAtPiB0b2RvbGlzdHN2YzogUE9TVCBhZGR0b2RvaXRlbVxudG9kb2xpc3RzdmMgLT4gbW9uZ28gOiBzdG9yZSB0b2RvaXRlbVxudG9kb2xpc3RzdmMgLS0-IEZyb250ZW5kOiB0b2RvaXRlbSByZXNwb25zZVxuRnJvbnRlbmQgLT4gdG9kb2xpc3RzdmM6IEdFVCB0b2RvbGlzdFxudG9kb2xpc3RzdmMgLT4gbW9uZ28gOiByZXRyaWV2ZSBhbGwgaXRlbXNcbnRvZG9saXN0c3ZjIC0tPiBGcm9udGVuZDogdG9kb2xpc3QgcmVzcG9uc2VcbkZyb250ZW5kIC0-IHRvZG9saXN0c3ZjOiBQT1NUIC91cGRhdGV0b2RvaXRlbVxudG9kb2xpc3RzdmMgLT4gbW9uZ28gOiB1cGRhdGUgb25lIGl0ZW1cbnRvZG9saXN0c3ZjIC0tPiBGcm9udGVuZDogdG9kb2xpc3QgcmVzcG9uc2VcbiIsIm1lcm1haWQiOnsidGhlbWUiOiJkZWZhdWx0In19)](https://mermaid-js.github.io/mermaid-live-editor/#/edit/eyJjb2RlIjoic2VxdWVuY2VEaWFncmFtXG5Gcm9udGVuZCAtPiB0b2RvbGlzdHN2YzogUE9TVCBhZGR0b2RvaXRlbVxudG9kb2xpc3RzdmMgLT4gbW9uZ28gOiBzdG9yZSB0b2RvaXRlbVxudG9kb2xpc3RzdmMgLS0-IEZyb250ZW5kOiB0b2RvaXRlbSByZXNwb25zZVxuRnJvbnRlbmQgLT4gdG9kb2xpc3RzdmM6IEdFVCB0b2RvbGlzdFxudG9kb2xpc3RzdmMgLT4gbW9uZ28gOiByZXRyaWV2ZSBhbGwgaXRlbXNcbnRvZG9saXN0c3ZjIC0tPiBGcm9udGVuZDogdG9kb2xpc3QgcmVzcG9uc2VcbkZyb250ZW5kIC0-IHRvZG9saXN0c3ZjOiBQT1NUIC91cGRhdGV0b2RvaXRlbVxudG9kb2xpc3RzdmMgLT4gbW9uZ28gOiB1cGRhdGUgb25lIGl0ZW1cbnRvZG9saXN0c3ZjIC0tPiBGcm9udGVuZDogdG9kb2xpc3QgcmVzcG9uc2VcbiIsIm1lcm1haWQiOnsidGhlbWUiOiJkZWZhdWx0In19)

Deployment
==========

Deployment strategy.
- Deploy frontend to connect to the todolist service or LB.
- Deploy one mongo database docker container
    * Ideally a sharded version of mongo (Primary for writing and slaves for reading)
- Deploy multiple todolist service docker containers connecting to the same database
    * Why deploy multiple services behind a LB? Multiple as we can do cascade updates:
        update one service that is backwards compatible, check logs, maybe update mongodb in case todoitems have changes with dbscript 
        when it works: update other services as well
        if it does not work; rollback the new service and fix issues

Other ideas
===========
This is a very bare version http project whereby code is all in main. Simple and effective but hard to add (unit)tests. Please see other project for when using proto files for messaging between client/server and a better differentiation of models, databaserepository, services.

State: Some masterdata should still be persisted and exposed:
- todo-state: each item can have only one state 'TODO', 'BUSY', 'DONE'
- labels: each item can have one or  more labels to assign to items
  - the user might want to create or remove labels
