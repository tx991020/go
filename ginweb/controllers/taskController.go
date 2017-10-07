package controllers


import (
	"net/http"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/gin-gonic/gin"
	"taskmanager/common"
	"taskmanager/data"
)

// CreateTask insert a new Task document
// Handler for HTTP Post - "/tasks
func CreateTask(c *gin.Context) {
	var dataResource TaskResource
	// Decode the incoming Task json
	err := c.BindJSON(&dataResource)
	if err != nil {
		common.DisplayAppError(
			c,
			err,
			"Invalid Task data",
			500,
		)
		return
	}
	task := &dataResource.Data
	context := NewContext()
	defer context.Close()
	context.User = task.CreatedBy

	col := context.DbCollection("tasks")
	repo := &data.TaskRepository{C: col}

	// Insert a task document
	repo.Create(task)

	c.JSON(http.StatusCreated, TaskResource{Data: *task})
}

// GetTasks returns all Task document
// Handler for HTTP Get - "/tasks"
func GetTasks(c *gin.Context) {
	context := NewContext()
	defer context.Close()
	col := context.DbCollection("tasks")
	repo := &data.TaskRepository{C: col}
	tasks := repo.GetAll()

	c.JSON(http.StatusOK, TasksResource{Data: tasks})
}

// GetTaskByID returns a single Task document by id
// Handler for HTTP Get - "/tasks/t/:id"
func GetTaskByID(c *gin.Context) {
	// Get id from the incoming url
	id := c.Param("id")
	context := NewContext()
	defer context.Close()

	col := context.DbCollection("tasks")
	repo := &data.TaskRepository{C: col}
	task, err := repo.GetById(id)
	if err != nil {
		if err == mgo.ErrNotFound {
			c.JSON(http.StatusNoContent, gin.H{})

		} else {
			common.DisplayAppError(
				c,
				err,
				"An unexpected error has occurred",
				500,
			)

		}
		return
	}
	c.JSON(http.StatusOK, TaskResource{Data: task})
}

// GetTasksByUser returns all Tasks created by a User
// Handler for HTTP Get - "/tasks/users/:email"
func GetTasksByUser(c *gin.Context) {
	// Get id from the incoming url
	user := c.Param("email")
	context := NewContext()
	defer context.Close()
	col := context.DbCollection("tasks")
	repo := &data.TaskRepository{C: col}
	tasks := repo.GetByUser(user)

	c.JSON(http.StatusOK, TasksResource{Data: tasks})
}

// UpdateTask update an existing Task document
// Handler for HTTP Put - "/tasks/:id"
func UpdateTask(c *gin.Context) {
	// Get id from the incoming url
	id := bson.ObjectIdHex(c.Param("id"))
	var dataResource TaskResource
	// Decode the incoming Task json
	err := c.BindJSON(&dataResource)
	if err != nil {
		common.DisplayAppError(
			c,
			err,
			"Invalid Task data",
			500,
		)
		return
	}
	task := &dataResource.Data
	task.Id = id
	context := NewContext()
	defer context.Close()
	col := context.DbCollection("tasks")
	repo := &data.TaskRepository{C: col}
	// Update an existing Task document
	if err := repo.Update(task); err != nil {
		common.DisplayAppError(
			c,
			err,
			"An unexpected error has occurred",
			500,
		)
		return
	}
	c.Status(http.StatusNoContent)
}

// DeleteTask deelete an existing Task document
// Handler for HTTP Delete - "/tasks/:id"
func DeleteTask(c *gin.Context) {

	id := c.Param("id")
	context := NewContext()
	defer context.Close()
	col := context.DbCollection("tasks")
	repo := &data.TaskRepository{C: col}
	// Delete an existing Task document
	err := repo.Delete(id)
	if err != nil {
		common.DisplayAppError(
			c,
			err,
			"An unexpected error has occurred",
			500,
		)
		return
	}
	c.Status(http.StatusNoContent)
}

