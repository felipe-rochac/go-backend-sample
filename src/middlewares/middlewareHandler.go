package middlewares

import (
	"backend-sample/common"
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

type KeyValue struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

var ErrorCodeKey KeyValue

type ErrorResponse struct {
	ErrorCode string `json,yaml:"error_code"`
	Message   string `json,yaml:"message"`
}

func MiddlewareHandler(c *gin.Context) {

	c.Next()

	handleError(c)

	formatResponse(c)
}

func formatResponse(c *gin.Context) {
	// Retrieve the response data
	response, exists := c.Get("response")
	if !exists {
		return
	}

	formatHttpResponse(http.StatusOK, response, c)
}

func handleError(c *gin.Context) {
	err := c.Errors.Last()

	if err != nil {
		errCode := 500
		var errResponse ErrorResponse
		if berr, ok := err.Err.(*common.BackendError); ok {
			if errorCode, err := common.EncryptAES([]byte(ErrorCodeKey.Key), berr.Identifier); err != nil {
				errResponse = ErrorResponse{
					Message: "Error generating error code",
				}
			} else {
				errCode = berr.Code
				errResponse = ErrorResponse{
					ErrorCode: errorCode,
					Message:   berr.Message,
				}
			}
		} else {
			errResponse = ErrorResponse{
				Message: err.Error(),
			}
		}

		formatHttpResponse(errCode, errResponse, c)

		c.Abort()
	}
}

func formatHttpResponse(statusCode int, response interface{}, c *gin.Context) {
	switch c.GetHeader("Accept") {
	case "application/x-yaml":
		c.Header("Content-Type", "application/x-yaml")
		yamlData, err := yaml.Marshal(response)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate YAML response"})
			return
		}
		c.String(statusCode, string(yamlData))
	default:
		c.Header("Content-Type", "application/json")
		c.JSON(statusCode, response)
	}
}
