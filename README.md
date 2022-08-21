# todo-list-api

<b> To run the API: </b>

1. `docker-compose up -d` runs the postgresdb Docker container
2. `go run main.go` starts the server 

The API will be running on localhost:5000.

These are the available endpoints:

- `POST /api/todolists` create a todo list
- `GET /api/todolists` get all todo lists
- `GET /api/todolists/{id}` get a todo list by id
- `PUT /api/todolists/{id}` update a todo list by id
- `DELETE /api/todolists/{id}` delete a todo list by id
- `POST /api/todolists/{id}/todos` create a todo in the todo list
- `GET /api/todolists/{id}/todos` get all todos in a todo list
- `GET /api/todolists/{id}/todos/{todoId}` get a todo in a todo list
- `PUT /api/todolists/{id}/todos/{todoId}` update a todo in a todo list
- `DELETE /api/todolists/{id}/todos/{todoId}` delete a todo in a todo list

