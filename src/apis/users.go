package apis

import (
	"backend-sample/database"
	"backend-sample/workflows"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

var userWorkflow workflows.UserWorkflowService

// Initialize sets up the necessary services and repositories for APIs
func Initialize(db database.MySqlDatabaseService) {
	// Initialize the repository
	repository := database.NewRepository(db)

	// Initialize the UserWorkflowService with the repository
	userWorkflow = *workflows.NewUserWorkflow(repository)
}

func GetUser(c *gin.Context) {
	userId, _ := c.GetQuery("user_id")
	name, _ := c.GetQuery("name")

	result, err := userWorkflow.GetUsers(userId, name)

	if err != nil {
		c.Errors = append(c.Errors, c.Error(err))
		return
	}

	b, backErr := json.Marshal(result)

	if backErr != nil {
		c.Errors = append(c.Errors, c.Error(err))
		return
	}

	response := string(b)

	c.JSON(http.StatusOK, response)
}

func AddUser(c *gin.Context) {
	var body workflows.UserRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "invalid payload"})
		return
	}

	result, err := userWorkflow.Create(body)

	if err != nil {
		c.Errors = append(c.Errors, c.Error(err))
		return
	}

	b, backErr := json.Marshal(result)

	if backErr != nil {
		c.Errors = append(c.Errors, c.Error(err))
		return
	}

	response := string(b)

	c.JSON(http.StatusOK, response)
}
