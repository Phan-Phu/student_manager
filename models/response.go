package models

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	ErrorCodeCanNotCreateStudent = 1001
	ErrorCodeFailCreateStudent   = 1002
	ErrorCodeFailDeleteStudent   = 1003
	ErrorCodeFailRetrieveStudent = 1004
	ErrorCodeStudentIsNotFound   = 1005
	ErrorCodeStudentIsEmpty      = 1006
	ErrorCodeCanNotUpdateStudent = 1007
	ErrorCodeCanNotFindClass     = 1008
	ErrorCodeMaxStudentInClass   = 1009
	ErrorCodeNotFoundClass       = 1010
	ErrorCodeInputIsWrong        = 1011
	ErrorCodeGetStudentDetails   = 1012
)

// Error messages mapped with error codes
var errorMessages = map[int]string{
	ErrorCodeCanNotCreateStudent: "Can not create student",
	ErrorCodeFailCreateStudent:   "Failed to create student",
	ErrorCodeFailDeleteStudent:   "Failed to delete student",
	ErrorCodeFailRetrieveStudent: "Failed to retrieve student",
	ErrorCodeStudentIsNotFound:   "Student is not found",
	ErrorCodeStudentIsEmpty:      "Student is empty",
	ErrorCodeCanNotUpdateStudent: "Can not update student",
	ErrorCodeCanNotFindClass:     "Not Found Class",
	ErrorCodeMaxStudentInClass:   "Max student in class",
	ErrorCodeNotFoundClass:       "Not found class",
	ErrorCodeInputIsWrong:        "Input is wrong",
	ErrorCodeGetStudentDetails:   "Error Get Student details",
}

func GetErrorMessage(code int) string {
	if message, exists := errorMessages[code]; exists {
		return message
	}
	return "Unknown error"
}

// Response Base response
type Response struct {
	StatusCode int            `json:"-"`
	Success    bool           `json:"success"`
	Message    string         `json:"message,omitempty"`
	Data       map[string]any `json:"data,omitempty"`
}

func (response *Response) SendResponse(c *gin.Context) {
	c.AbortWithStatusJSON(response.StatusCode, response)
}

func SendResponseData(c *gin.Context, data gin.H) {
	response := &Response{
		StatusCode: http.StatusOK,
		Success:    true,
		Data:       data,
	}
	response.SendResponse(c)
}

func SendErrorResponse(c *gin.Context, status int, message string) {
	response := &Response{
		StatusCode: status,
		Success:    false,
		Message:    message,
	}
	response.SendResponse(c)
}

func GetErrorResponse(status int, data string) *Response {
	response := &Response{
		StatusCode: status,
		Success:    false,
		Message:    data,
	}
	return response
}
