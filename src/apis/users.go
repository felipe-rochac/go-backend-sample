package apis

import (
	"backend-sample/workflows"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

var userWorkflow workflows.UserWorkflowService

func GetUser(c *gin.Context) {
	userId, _ := c.GetQuery("user_id")
	name, _ := c.GetQuery("name")

	response, err := userWorkflow.GetUsers(userId, name)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	b, backErr := json.Marshal(response)

	if backErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": backErr.Error(),
		})
		return
	}

	content := string(b)

	c.JSON(http.StatusOK, content)
}
