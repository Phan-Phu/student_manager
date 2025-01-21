package models

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	// General
	ErrorCodeIDIsWrong = 0001

	// Student
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

	// Teacher
	ErrorCodeCanNotCreateTeacher       = 2001
	ErrorCodeFailCreateTeacher         = 2002
	ErrorCodeFailDeleteTeacher         = 2003
	ErrorCodeFailRetrieveTeacher       = 2004
	ErrorCodeTeacherIsNotFound         = 2005
	ErrorCodeTeacherIsEmpty            = 2006
	ErrorCodeCanNotUpdateTeacher       = 2007
	ErrorCodeCanNotFindTeacher         = 2008
	ErrorCodeMaxTeacherInClass         = 2009
	ErrorCodeNotFoundTeacher           = 2010
	ErrorCodeInputIsWrongTeacher       = 2011
	ErrorCodeGetTeacherDetails         = 2012
	ErrorCodeUpdateClassIsEmpty        = 2013
	ErrorCodeTeacherNotAssignedToClass = 2014

	ErrorCodeCanNotFindCourse        = 3006
	ErrorCodeCanNotUpdateCourse      = 3007
	ErrorCodeCanNotFindScoreInCourse = 3008
	ErrorCodeGetCourseDetails        = 3012
	ErrorCodeCreateManyScore         = 3013
	ErrorCodeGetScoreDetails         = 3015

	ErrorCodeCanNotCreateClass = 4015
	ErrorCodeCanNotUpdateClass = 4016
	ErrorCodeFailDeleteClass   = 4017
)

// Error messages mapped with error codes
var errorMessages = map[int]string{
	//General
	ErrorCodeIDIsWrong: "Id Is Wrong",
	// student
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

	//teacher
	ErrorCodeCanNotCreateTeacher:       "Can not create teacher",
	ErrorCodeFailCreateTeacher:         "Failed to create teacher",
	ErrorCodeFailDeleteTeacher:         "Failed to delete teacher",
	ErrorCodeFailRetrieveTeacher:       "Failed to retrieve teacher",
	ErrorCodeTeacherIsNotFound:         "Teacher is not found",
	ErrorCodeTeacherIsEmpty:            "Teacher is empty",
	ErrorCodeCanNotUpdateTeacher:       "Can not update teacher",
	ErrorCodeCanNotFindTeacher:         "Teacher not found",
	ErrorCodeMaxTeacherInClass:         "Max teacher in class",
	ErrorCodeNotFoundTeacher:           "Teacher not found",
	ErrorCodeInputIsWrongTeacher:       "Input is wrong",
	ErrorCodeGetTeacherDetails:         "Error getting teacher details",
	ErrorCodeUpdateClassIsEmpty:        "Error update class is empty",
	ErrorCodeTeacherNotAssignedToClass: "Error teacher is not assign in class",

	ErrorCodeCanNotCreateClass: "Error class can not create",
	ErrorCodeCanNotUpdateClass: "Error class can not update",
	ErrorCodeFailDeleteClass:   "Error class can not delete",

	ErrorCodeGetCourseDetails:        "Error  can not get detail course",
	ErrorCodeCanNotFindScoreInCourse: "Error  can not get find score id in course",
	ErrorCodeGetScoreDetails:         "Error  can not get score details",
	ErrorCodeCanNotUpdateCourse:      "Error course can not update",
	ErrorCodeCanNotFindCourse:        "Error course can not find",
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
