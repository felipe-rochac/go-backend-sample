package middlewares

import (
	"backend-sample/common"

	"github.com/gin-gonic/gin"
)

type KeyValue struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

var ErrorCodeKey KeyValue

func ErrorHandler(c *gin.Context) {

	c.Next()

	err := c.Errors.Last()

	if err != nil {
		if berr, ok := err.Err.(*common.BackendError); ok {
			if errorCode, err := common.EncryptAES([]byte(ErrorCodeKey.Key), berr.Identifier); err != nil {
				c.JSON(500, gin.H{
					"message": "Error generating error code",
				})
			} else {
				c.JSON(berr.Code, gin.H{
					"error_code": errorCode,
					"message":    "An error occurred",
				})
			}
		} else {
			c.JSON(500, gin.H{
				"message": err.Error(),
			})
		}

		c.Abort()
	}
}
