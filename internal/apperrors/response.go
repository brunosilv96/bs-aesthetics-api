package apperrors

import "github.com/gin-gonic/gin"

func WriteError(c *gin.Context, err *ApiError) {
	response := gin.H{
		"code":    err.Code,
		"message": err.Message,
	}

	if err.Details != nil {
		response["details"] = err.Details
	}

	c.AbortWithStatusJSON(err.Status, response)
}
