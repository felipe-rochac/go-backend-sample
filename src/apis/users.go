package apis

import (
	"backend-sample/database"
	"backend-sample/workflows"
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

	response, err := userWorkflow.GetUsers(userId, name)

	if err != nil {
		c.Errors = append(c.Errors, c.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

func AddUser(c *gin.Context) {
	var body workflows.UserRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "invalid payload"})
		return
	}

	response, err := userWorkflow.Create(body)

	if err != nil {
		c.Errors = append(c.Errors, c.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}
